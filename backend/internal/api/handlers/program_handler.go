package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/devdavidalonso/cecor/backend/internal/models"
)

// ProgramHandler handles program endpoints.
type ProgramHandler struct {
	db *gorm.DB
}

// NewProgramHandler creates ProgramHandler.
func NewProgramHandler(db *gorm.DB) *ProgramHandler {
	return &ProgramHandler{db: db}
}

// ListPrograms returns programs available in the educational center.
// Query params:
// - activeOnly (bool, default true): when true returns only active programs.
func (h *ProgramHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	activeOnly := true
	if v := r.URL.Query().Get("activeOnly"); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			activeOnly = parsed
		}
	}

	query := h.db.WithContext(r.Context()).
		Model(&models.Program{}).
		Preload("Center").
		Order("name ASC")

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	var programs []models.Program
	if err := query.Find(&programs).Error; err != nil {
		http.Error(w, "failed to list programs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(programs)
}
