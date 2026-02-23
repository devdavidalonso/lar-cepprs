package routes

import (
	"net/http"

	"github.com/devdavidalonso/cecor/backend/internal/api/handlers"
	"github.com/devdavidalonso/cecor/backend/internal/api/middleware"
	"github.com/devdavidalonso/cecor/backend/internal/config"
	"github.com/go-chi/chi/v5"
)

// Register configura as rotas base da API v1 em um router já prefixado
func Register(r chi.Router, cfg *config.Config, authHandler *handlers.AuthHandler, courseHandler *handlers.CourseHandler, enrollmentHandler *handlers.EnrollmentHandler, attendanceHandler *handlers.AttendanceHandler, reportHandler *handlers.ReportHandler, teacherHandler *handlers.TeacherHandler) {
	// Rotas públicas da v1 (sem autenticação)
	r.Group(func(r chi.Router) {
		r.Get("/auth/sso/login", authHandler.SSOLogin)
		r.Get("/auth/sso/callback", authHandler.SSOCallback)
	})

	// Rotas protegidas da v1
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(cfg))

		// Endpoint de verificação de token
		r.Get("/auth/verify", authHandler.Verify)

		// Courses
		r.Route("/courses", func(r chi.Router) {
			r.Get("/", courseHandler.ListCourses)
			r.Post("/", courseHandler.CreateCourse)
			r.Get("/{id}", courseHandler.GetCourse)
			r.Put("/{id}", courseHandler.UpdateCourse)
			r.Delete("/{id}", courseHandler.DeleteCourse)
		})

		// Enrollments
		r.Route("/enrollments", func(r chi.Router) {
			r.Get("/", enrollmentHandler.ListEnrollments)
			r.Post("/", enrollmentHandler.EnrollStudent)
			r.Get("/{id}", enrollmentHandler.GetEnrollment)
			r.Put("/{id}", enrollmentHandler.UpdateEnrollment)
			r.Delete("/{id}", enrollmentHandler.DeleteEnrollment)
		})

		// Attendance
		r.Route("/attendance", func(r chi.Router) {
			r.Post("/record", attendanceHandler.RecordBatch)
			r.Get("/course/{id}/date/{date}", attendanceHandler.GetClassAttendance)
			r.Get("/student/{id}", attendanceHandler.GetStudentHistory)
			r.Get("/student/{id}/percentage", attendanceHandler.GetStudentPercentage)
		})

		// Notifications
		r.Route("/notifications", func(r chi.Router) {
			r.Get("/", http.NotFound)
			r.Post("/", http.NotFound)
			r.Get("/{id}", http.NotFound)
			r.Put("/{id}/read", http.NotFound)
		})

		// Reports
		r.Route("/reports", func(r chi.Router) {
			r.Get("/attendance/course/{id}", reportHandler.GetCourseAttendanceReport)
			r.Get("/attendance/student/{id}", reportHandler.GetStudentAttendanceReport)
			r.Get("/students", http.NotFound)
			r.Get("/courses", http.NotFound)
		})

		// Interviews placeholder
		r.Route("/interviews", func(r chi.Router) {
			r.Get("/", http.NotFound)
			r.Post("/", http.NotFound)
			r.Get("/{id}", http.NotFound)
			r.Put("/{id}", http.NotFound)
		})

		// Volunteering placeholder
		r.Route("/volunteering", func(r chi.Router) {
			r.Get("/terms", http.NotFound)
			r.Post("/terms", http.NotFound)
			r.Get("/terms/{id}", http.NotFound)
			r.Post("/terms/{id}/sign", http.NotFound)
		})

		// Users & Permissions (admin)
		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.RequireAdmin)
			r.Get("/", http.NotFound)
			r.Post("/", http.NotFound)
			r.Get("/teachers", courseHandler.ListProfessors)
			r.Get("/{id}", http.NotFound)
			r.Put("/{id}", http.NotFound)
			r.Delete("/{id}", http.NotFound)
		})

		// Teachers
		r.Route("/teachers", func(r chi.Router) {
			r.Post("/", teacherHandler.CreateTeacher)
			r.Get("/", teacherHandler.GetTeachers)
			r.Get("/{id}", teacherHandler.GetTeacherByID)
			r.Put("/{id}", teacherHandler.UpdateTeacher)
			r.Delete("/{id}", teacherHandler.DeleteTeacher)
		})
	})
}
