package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/service/teachers"
)

// TeacherHandler handles HTTP requests for teachers
type TeacherHandler struct {
	service teachers.Service
}

// NewTeacherHandler creates a new instance of TeacherHandler
func NewTeacherHandler(service teachers.Service) *TeacherHandler {
	return &TeacherHandler{
		service: service,
	}
}

// CreateTeacher handles the creation of a new teacher.
func (h *TeacherHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	var teacher models.User
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTeacher(r.Context(), &teacher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(teacher)
}

// GetTeachers handles the retrieval of all teachers.
func (h *TeacherHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.service.GetTeachers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

// GetTeacherByID handles the retrieval of a teacher by ID.
func (h *TeacherHandler) GetTeacherByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	teacher, err := h.service.GetTeacherByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// UpdateTeacher handles the update of a teacher.
func (h *TeacherHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var teacher models.User
	if err := json.NewDecoder(r.Body).Decode(&teacher); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	teacher.ID = uint(id)

	if err := h.service.UpdateTeacher(r.Context(), &teacher); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// DeleteTeacher handles the deletion of a teacher.
func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTeacher(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Deprecated: use CreateTeacher.
func (h *TeacherHandler) CreateProfessor(w http.ResponseWriter, r *http.Request) {
	h.CreateTeacher(w, r)
}

// Deprecated: use GetTeachers.
func (h *TeacherHandler) GetProfessors(w http.ResponseWriter, r *http.Request) {
	h.GetTeachers(w, r)
}

// Deprecated: use GetTeacherByID.
func (h *TeacherHandler) GetProfessorByID(w http.ResponseWriter, r *http.Request) {
	h.GetTeacherByID(w, r)
}

// Deprecated: use UpdateTeacher.
func (h *TeacherHandler) UpdateProfessor(w http.ResponseWriter, r *http.Request) {
	h.UpdateTeacher(w, r)
}

// Deprecated: use DeleteTeacher.
func (h *TeacherHandler) DeleteProfessor(w http.ResponseWriter, r *http.Request) {
	h.DeleteTeacher(w, r)
}
