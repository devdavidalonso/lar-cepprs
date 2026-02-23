package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt" // Adicione esta importação para o bcrypt
	"gorm.io/gorm"

	"github.com/devdavidalonso/cecor/backend/internal/api/handlers"
	apiMiddleware "github.com/devdavidalonso/cecor/backend/internal/api/middleware"
	"github.com/devdavidalonso/cecor/backend/internal/api/routes"
	"github.com/devdavidalonso/cecor/backend/internal/auth"
	"github.com/devdavidalonso/cecor/backend/internal/config"
	"github.com/devdavidalonso/cecor/backend/internal/database"
	"github.com/devdavidalonso/cecor/backend/internal/infrastructure/googleapis"
	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository/mongodb"
	"github.com/devdavidalonso/cecor/backend/internal/repository/postgres"
	"github.com/devdavidalonso/cecor/backend/internal/service/attendance" // Adicionar importação de attendance
	"github.com/devdavidalonso/cecor/backend/internal/service/courses"    // Adicionar importação de courses
	"github.com/devdavidalonso/cecor/backend/internal/service/email"
	"github.com/devdavidalonso/cecor/backend/internal/service/enrollments" // Adicionar importação de enrollments
	"github.com/devdavidalonso/cecor/backend/internal/service/incidents"
	"github.com/devdavidalonso/cecor/backend/internal/service/interviews" // Import interviews service
	"github.com/devdavidalonso/cecor/backend/internal/service/keycloak"
	"github.com/devdavidalonso/cecor/backend/internal/service/reports"  // Adicionar importação de reports
	"github.com/devdavidalonso/cecor/backend/internal/service/students" // Adicionar importação de students
	"github.com/devdavidalonso/cecor/backend/internal/service/teacherportal"
	"github.com/devdavidalonso/cecor/backend/internal/service/teachers" // Adicionar importação de professors
	"github.com/devdavidalonso/cecor/backend/internal/service/users"    // Adicionar esta importação
	"github.com/devdavidalonso/cecor/backend/pkg/logger"
)

// Adicione esta função para atualizar a senha do usuário
func updateUserPassword(db *gorm.DB, email, password string) {
	// Gerar novo hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erro ao gerar hash: %v", err)
		return
	}

	// Atualizar no banco
	result := db.Model(&models.User{}).
		Where("email = ?", email).
		Update("password", string(hashedPassword))

	if result.Error != nil {
		log.Printf("Erro ao atualizar senha: %v", result.Error)
	} else {
		log.Printf("Senha atualizada com sucesso para %s", email)
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	appLogger := logger.NewLogger()
	appLogger.Info("Starting CECOR Educational Management System")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		appLogger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize SSO Provider
	// Extract issuer URL from AuthURL (remove the /protocol/openid-connect/auth part)
	issuerURL := strings.TrimSuffix(cfg.SSO.AuthURL, "/protocol/openid-connect/auth")
	if err := auth.InitProvider(context.Background(), issuerURL); err != nil {
		appLogger.Fatal("Failed to initialize SSO provider", "error", err)
	}

	// Initialize database
	db, err := postgres.InitDB(cfg)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", "error", err)
	}

	// Perform database migrations
	err = postgres.MigrateDB(db)
	if err != nil {
		appLogger.Fatal("Failed to migrate database", "error", err)
	}
	appLogger.Info("Database migration completed successfully")

	// Atualizar senha do usuário de teste
	updateUserPassword(db, "maria.silva@cecor.org", "cecor2024!")

	// Initialize MongoDB
	mongoClient, err := database.InitMongoDB(cfg)
	if err != nil {
		appLogger.Warn("Failed to connect to MongoDB, proceeding without it", "error", err)
	} else {
		defer func() {
			if err := mongoClient.Disconnect(context.Background()); err != nil {
				appLogger.Error("Failed to disconnect MongoDB", "error", err)
			}
		}()
		appLogger.Info("Connected to MongoDB successfully")
	}

	// Initialize repositories
	studentRepo := postgres.NewStudentRepository(db)
	userRepo := postgres.NewUserRepository(db)             // Adicionar o repositório de usuários
	courseRepo := postgres.NewCourseRepository(db)         // Adicionar repositório de cursos
	enrollmentRepo := postgres.NewEnrollmentRepository(db) // Adicionar repositório de matrículas
	attendanceRepo := postgres.NewAttendanceRepository(db) // Adicionar repositório de presenças
	reportRepo := postgres.NewReportRepository(db)         // Adicionar repositório de relatórios
	programRepo := postgres.NewProgramRepository(db)       // Adicionar repositório de programas

	// Initialize Google Classroom Client
	classroomClient, err := googleapis.NewGoogleClassroomClient("credentials.json")
	if err != nil {
		appLogger.Error("Failed to initialize Google Classroom client (integration will be disabled)", "error", err)
		classroomClient = nil // explicitly nil if failed
	}

	// Initialize services
	keycloakService := keycloak.NewKeycloakService()                                                      // Inicializar keycloak service
	emailService := email.NewEmailService()                                                               // Inicializar email service
	studentService := students.NewStudentService(studentRepo, programRepo, keycloakService, emailService) // Inicializar student service com ProgramRepo
	userService := users.NewUserService(userRepo)                                                         // Adicionar o serviço de usuários
	teacherService := teachers.NewService(userRepo, programRepo, keycloakService, emailService)           // Adicionar serviço de teachers com ProgramRepo
	courseService := courses.NewService(courseRepo, classroomClient)                                      // Adicionar serviço de cursos
	enrollmentService := enrollments.NewService(enrollmentRepo, studentRepo, courseRepo, classroomClient) // Updated with dependencies
	attendanceService := attendance.NewService(attendanceRepo)                                            // Adicionar serviço de presenças
	reportService := reports.NewService(reportRepo)                                                       // Adicionar serviço de relatórios

	// Initialize MongoDB repositories
	formRepo := mongodb.NewFormRepository()
	interviewService := interviews.NewService(formRepo)           // Inicializar serviço de entrevistas
	interviewAdminService := interviews.NewAdminService(formRepo) // Inicializar serviço admin de entrevistas

	// Initialize SSO Config
	ssoConfig := auth.NewSSOConfig(cfg)

	// Initialize handlers
	studentHandler := handlers.NewStudentHandler(studentService)
	authHandler := handlers.NewAuthHandler(userService, cfg, ssoConfig)                     // Adicionar o handler de autenticação
	teacherHandler := handlers.NewTeacherHandler(teacherService)                            // Adicionar handler de professores
	courseHandler := handlers.NewCourseHandler(courseService, keycloakService)              // Adicionar handler de cursos
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentService, interviewService) // Adicionar handler de matrículas com entrevista
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)                   // Adicionar handler de presenças
	reportHandler := handlers.NewReportHandler(reportService)                               // Adicionar handler de relatórios
	interviewHandler := handlers.NewInterviewHandler(interviewService)                      // Adicionar handler de entrevistas
	interviewAdminHandler := handlers.NewInterviewAdminHandler(interviewAdminService)       // Adicionar handler admin de entrevistas

	// Initialize teacher portal service and handler
	teacherPortalService := teacherportal.NewService(db)
	teacherPortalHandler := handlers.NewTeacherPortalHandler(teacherPortalService)

	// Initialize incident service and handler
	incidentService := incidents.NewService(db)
	incidentHandler := handlers.NewIncidentHandler(incidentService)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// CORS middleware
	r.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"))
	r.Use(middleware.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization"))

	// Options handling for CORS preflight requests
	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Initialize new handlers for Phase 2
	studentPortalHandler := handlers.NewStudentPortalHandler(db)
	courseClassHandler := handlers.NewCourseClassHandler(db)
	skillHandler := handlers.NewSkillHandler(db)
	migrationHandler := handlers.NewMigrationHandler(db)
	programHandler := handlers.NewProgramHandler(db)

	// Registrar todas as rotas v1 sob um único prefixo /api/v1
	appLogger.Info("Registering v1 routes...")
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// Rotas base (públicas + protegidas)
			routes.Register(r, cfg, authHandler, courseHandler, enrollmentHandler, attendanceHandler, reportHandler, teacherHandler)

			// Módulos adicionais protegidos
			r.Group(func(r chi.Router) {
				r.Use(apiMiddleware.Authenticate(cfg))

				// Programs
				r.Get("/programs", programHandler.ListPrograms)

				// Interview endpoints
				interviewHandler.RegisterRoutes(r)

				// Teacher portal / incidents / student portal
				teacherPortalHandler.RegisterRoutes(r)
				incidentHandler.RegisterRoutes(r)
				studentPortalHandler.RegisterRoutes(r)

				// Course class / skills
				courseClassHandler.RegisterRoutes(r)
				skillHandler.RegisterRoutes(r)

				// Students CRUD
				r.Route("/students", func(r chi.Router) {
					r.Get("/", studentHandler.GetStudents)
					r.Post("/", studentHandler.CreateStudent)
					r.Get("/{id}", studentHandler.GetStudent)
					r.Put("/{id}", studentHandler.UpdateStudent)
					r.Delete("/{id}", studentHandler.DeleteStudent)
					r.Get("/{id}/guardians", studentHandler.GetGuardians)
					r.Post("/{id}/guardians", studentHandler.AddGuardian)
					r.Get("/{id}/documents", studentHandler.GetDocuments)
					r.Post("/{id}/documents", studentHandler.AddDocument)
					r.Get("/{id}/notes", studentHandler.GetNotes)
					r.Post("/{id}/notes", studentHandler.AddNote)
				})
			})

			// Rotas admin sob /api/v1/admin
			r.Route("/admin", func(r chi.Router) {
				r.Use(apiMiddleware.Authenticate(cfg))
				r.Use(apiMiddleware.RequireAdmin)

				interviewAdminHandler.RegisterAdminRoutes(r)
				r.Post("/migrations/run", migrationHandler.RunMigrations)
				r.Post("/migrations/data", migrationHandler.RunDataMigration)
				r.Post("/migrations/rollback", migrationHandler.RollbackMigrations)
				r.Get("/migrations/status", migrationHandler.GetMigrationStatus)
			})
		})
	})
	appLogger.Info("V1 routes registered")

	// Get server port
	port := cfg.Server.Port
	addr := fmt.Sprintf(":%d", port)

	// Create HTTP server
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Imprimir todas as rotas registradas
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	// Start server in a goroutine
	go func() {
		appLogger.Info("Server starting", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server gracefully stopped")
}
