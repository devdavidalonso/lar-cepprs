package routes

import (
	"github.com/devdavidalonso/cecor/backend/internal/api/handlers"
	"github.com/devdavidalonso/cecor/backend/internal/api/middleware"
	"github.com/devdavidalonso/cecor/backend/internal/config"
	"github.com/go-chi/chi/v5"
)

// SetupStudentRoutes configures routes for student resources
func SetupStudentRoutes(r chi.Router, cfg *config.Config, handler *handlers.StudentHandler) {
	// Public routes if needed
	r.Group(func(r chi.Router) {
		// There are no public student-related routes
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(cfg))

		// Student routes
		r.Route("/students", func(r chi.Router) {
			r.Get("/", handler.GetStudents)          // List students
			r.Post("/", handler.CreateStudent)       // Create student
			r.Get("/{id}", handler.GetStudent)       // Get student by ID
			r.Put("/{id}", handler.UpdateStudent)    // Update student
			r.Delete("/{id}", handler.DeleteStudent) // Delete student

			// Sub-routes for guardians
			r.Route("/{id}/guardians", func(r chi.Router) {
				r.Get("/", handler.GetGuardians) // List guardians
				r.Post("/", handler.AddGuardian) // Add guardian
			})

			// Sub-routes for documents
			r.Route("/{id}/documents", func(r chi.Router) {
				r.Get("/", handler.GetDocuments) // List documents
				r.Post("/", handler.AddDocument) // Add document
			})

			// Sub-routes for notes
			r.Route("/{id}/notes", func(r chi.Router) {
				r.Get("/", handler.GetNotes) // List notes
				r.Post("/", handler.AddNote) // Add note
			})
		})

		// Guardian routes
		r.Route("/guardians", func(r chi.Router) {
			r.Put("/{id}", handler.UpdateGuardian)    // Update guardian
			r.Delete("/{id}", handler.DeleteGuardian) // Delete guardian
		})

		// Document routes
		r.Route("/documents", func(r chi.Router) {
			r.Delete("/{id}", handler.DeleteDocument) // Delete document
		})
	})
}
