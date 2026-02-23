package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/service/courses"
	"github.com/devdavidalonso/cecor/backend/internal/service/keycloak"
	"github.com/go-chi/chi/v5"
)

type CourseHandler struct {
	service         courses.Service
	keycloakService *keycloak.KeycloakService
}

func NewCourseHandler(service courses.Service, keycloakService *keycloak.KeycloakService) *CourseHandler {
	return &CourseHandler{
		service:         service,
		keycloakService: keycloakService,
	}
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCourse(r.Context(), &course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	course, err := h.service.GetCourseByID(r.Context(), uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if course == nil {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.service.ListCourses(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(courses)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var course models.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	course.ID = uint(id)
	if err := h.service.UpdateCourse(r.Context(), &course); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteCourse(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CourseHandler) ListProfessors(w http.ResponseWriter, r *http.Request) {
	professors, err := h.keycloakService.GetUsersByRole(r.Context(), "teacher")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(professors)
}
