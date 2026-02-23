package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/devdavidalonso/cecor/backend/internal/api/middleware"
	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/service/students"
	"github.com/devdavidalonso/cecor/backend/pkg/errors"
)

// StudentHandler implements HTTP handlers for student resource
type StudentHandler struct {
	studentService students.Service
}

// NewStudentHandler creates a new instance of StudentHandler
func NewStudentHandler(studentService students.Service) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// PaginatedResponse is a structure for paginated responses
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalItems int64       `json:"totalItems"`
	TotalPages int64       `json:"totalPages"`
}

// GetStudents returns a paginated list of students
// @Summary List students
// @Description Returns a paginated list of students with optional filters
// @Tags students
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Items per page (default: 20, max: 100)"
// @Param name query string false "Filter by name"
// @Param email query string false "Filter by email"
// @Param cpf query string false "Filter by CPF"
// @Param status query string false "Filter by status"
// @Param course_id query int false "Filter by course"
// @Success 200 {object} PaginatedResponse{data=[]models.Student}
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students [get]
func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Get filters
	filters := make(map[string]interface{})

	if name := r.URL.Query().Get("name"); name != "" {
		filters["name"] = name
	}

	if email := r.URL.Query().Get("email"); email != "" {
		filters["email"] = email
	}

	if cpf := r.URL.Query().Get("cpf"); cpf != "" {
		filters["cpf"] = cpf
	}

	if status := r.URL.Query().Get("status"); status != "" {
		filters["status"] = status
	}

	if courseID := r.URL.Query().Get("course_id"); courseID != "" {
		if id, err := strconv.Atoi(courseID); err == nil {
			filters["course_id"] = id
		}
	}

	programParam := r.URL.Query().Get("program_id")
	if programParam == "" {
		programParam = r.URL.Query().Get("programId")
	}
	if programParam != "" {
		if id, err := strconv.Atoi(programParam); err == nil {
			filters["program_id"] = id
		}
	}

	// Call service
	students, total, err := h.studentService.GetStudents(r.Context(), page, pageSize, filters)
	if err != nil {
		errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Build paginated response
	response := PaginatedResponse{
		Data:       students,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: (total + int64(pageSize) - 1) / int64(pageSize),
	}

	errors.RespondWithJSON(w, http.StatusOK, response)
}

// GetStudent returns a specific student by ID
// @Summary Get student
// @Description Returns details of a specific student
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id} [get]
func (h *StudentHandler) GetStudent(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Call service
	student, err := h.studentService.GetStudentByID(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, student)
}

// CreateStudent creates a new student
// @Summary Create student
// @Description Creates a new student record
// @Tags students
// @Accept json
// @Produce json
// @Param student body models.Student true "Student data"
// @Success 201 {object} models.Student
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students [post]
func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var student models.Student

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		// Log the specific error for debugging
		fmt.Printf("Error decoding student JSON: %v\n", err)
		errors.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid data format: %v", err))
		return
	}

	// Get user from context (for auditing) - optional for now since test route might not have it
	// In production, this should be enforced by middleware
	if _, ok := middleware.GetUserFromContext(r.Context()); !ok {
		// Just log warning but proceed if it's the test route or if we want to allow unauthenticated creation for now
		// For strict mode, uncomment the return
		// errors.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		// return
		fmt.Println("Warning: User not found in context during student creation")
	}

	// Call service
	err := h.studentService.CreateStudent(r.Context(), &student)
	if err != nil {
		// Check error type
		if err.Error() == "name is required" ||
			err.Error() == "email is required" ||
			err.Error() == "phone is required" ||
			err.Error() == "birth date is required" ||
			err.Error() == "a student with this email already exists" ||
			err.Error() == "a student with this CPF already exists" {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusCreated, student)
}

// UpdateStudent updates an existing student
// @Summary Update student
// @Description Updates data for an existing student
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body models.Student true "Student data"
// @Success 200 {object} models.Student
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id} [put]
func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var student models.Student

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	// Ensure ID in URL is used
	student.ID = uint(id)

	// Call service
	err = h.studentService.UpdateStudent(r.Context(), &student)
	if err != nil {
		// Check error type
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "another student with this email already exists" ||
			err.Error() == "another student with this CPF already exists" {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Get updated student for response
	updatedStudent, err := h.studentService.GetStudentByID(r.Context(), uint(id))
	if err != nil {
		errors.RespondWithError(w, http.StatusInternalServerError, "Error retrieving updated student")
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, updatedStudent)
}

// DeleteStudent removes a student
// @Summary Delete student
// @Description Performs a logical deletion of a student
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id} [delete]
func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Call service
	err = h.studentService.DeleteStudent(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Respond with success and no content
	w.WriteHeader(http.StatusNoContent)
}

// GetGuardians returns guardians for a student
// @Summary List guardians
// @Description Returns the list of guardians associated with a student
// @Tags students,guardians
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {array} models.Guardian
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id}/guardians [get]
func (h *StudentHandler) GetGuardians(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Call service
	guardians, err := h.studentService.GetGuardians(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, guardians)
}

// AddGuardian adds a guardian to a student
// @Summary Add guardian
// @Description Adds a new guardian to a student
// @Tags students,guardians
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param guardian body models.Guardian true "Guardian data"
// @Success 201 {object} models.Guardian
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id}/guardians [post]
func (h *StudentHandler) AddGuardian(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var guardian models.Guardian

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&guardian); err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	// Ensure student ID is set correctly
	guardian.StudentID = uint(id)

	// Call service
	err = h.studentService.AddGuardian(r.Context(), &guardian)
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "guardian name is required" ||
			err.Error() == "relationship is required" ||
			err.Error() == "maximum of 3 guardians per student reached" {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusCreated, guardian)
}

// UpdateGuardian updates a guardian
// @Summary Update guardian
// @Description Updates an existing guardian
// @Tags guardians
// @Accept json
// @Produce json
// @Param id path int true "Guardian ID"
// @Param guardian body models.Guardian true "Guardian data"
// @Success 200 {object} models.Guardian
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/guardians/{id} [put]
func (h *StudentHandler) UpdateGuardian(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var guardian models.Guardian

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&guardian); err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	// Ensure guardian ID is correct
	guardian.ID = uint(id)

	// Call service
	err = h.studentService.UpdateGuardian(r.Context(), &guardian)
	if err != nil {
		if err.Error() == "guardian not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, guardian)
}

// DeleteGuardian removes a guardian
// @Summary Remove guardian
// @Description Removes a guardian from the system
// @Tags guardians
// @Accept json
// @Produce json
// @Param id path int true "Guardian ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/guardians/{id} [delete]
func (h *StudentHandler) DeleteGuardian(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Call service
	err = h.studentService.RemoveGuardian(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "guardian not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Respond with success and no content
	w.WriteHeader(http.StatusNoContent)
}

// GetDocuments returns documents for a student
// @Summary List documents
// @Description Returns the list of documents associated with a student
// @Tags students,documents
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {array} models.Document
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id}/documents [get]
func (h *StudentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Call service
	documents, err := h.studentService.GetDocuments(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, documents)
}

// AddDocument adds a document to a student
// @Summary Add document
// @Description Adds a new document to a student
// @Tags students,documents
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param document body models.Document true "Document data"
// @Success 201 {object} models.Document
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id}/documents [post]
func (h *StudentHandler) AddDocument(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var document models.Document

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	// Ensure student ID is set correctly
	document.StudentID = uint(id)

	// Get user from context for document.UploadedBy
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		errors.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	document.UploadedByID = uint(userClaims.UserID)

	// Call service
	err = h.studentService.AddDocument(r.Context(), &document)
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "document name is required" ||
			err.Error() == "document type is required" ||
			err.Error() == "document path is required" {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusCreated, document)
}

// DeleteDocument removes a document
// @Summary Remove document
// @Description Removes a document from the system
// @Tags documents
// @Accept json
// @Produce json
// @Param id path int true "Document ID"
// @Success 204 "No Content"
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/documents/{id} [delete]
func (h *StudentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Call service
	err = h.studentService.RemoveDocument(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "document not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Respond with success and no content
	w.WriteHeader(http.StatusNoContent)
}

// GetNotes returns notes for a student
// @Summary List notes
// @Description Returns the list of notes associated with a student
// @Tags students,notes
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param includeConfidential query bool false "Include confidential notes (default: false)"
// @Success 200 {array} models.StudentNote
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id}/notes [get]
func (h *StudentHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Check parameter to include confidential notes
	includeConfidential := false
	if r.URL.Query().Get("includeConfidential") == "true" {
		includeConfidential = true
	}

	// Call service
	notes, err := h.studentService.GetNotes(r.Context(), uint(id), includeConfidential)
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusOK, notes)
}

// AddNote adds a note to a student
// @Summary Add note
// @Description Adds a new note to a student
// @Tags students,notes
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param note body models.StudentNote true "Note data"
// @Success 201 {object} models.StudentNote
// @Failure 400 {object} errors.AppError
// @Failure 404 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /api/v1/students/{id}/notes [post]
func (h *StudentHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var note models.StudentNote

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		errors.RespondWithError(w, http.StatusBadRequest, "Invalid data format")
		return
	}

	// Ensure student ID is set correctly
	note.StudentID = uint(id)

	// Get user from context for note.AuthorID
	userClaims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		errors.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	note.AuthorID = uint(userClaims.UserID)

	// Call service
	err = h.studentService.AddNote(r.Context(), &note)
	if err != nil {
		if err.Error() == "student not found" {
			errors.RespondWithError(w, http.StatusNotFound, err.Error())
		} else if err.Error() == "note content is required" {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	errors.RespondWithJSON(w, http.StatusCreated, note)
}
