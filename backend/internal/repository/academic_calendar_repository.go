package repository

import (
	"context"
	"time"

	"github.com/devdavidalonso/cecor/backend/internal/models"
)

// AcademicCalendarRepository defines the interface for academic calendar data access
type AcademicCalendarRepository interface {
	Create(ctx context.Context, event *models.AcademicCalendar) error
	FindRecessesBetween(ctx context.Context, programID *uint, start, end time.Time) ([]models.AcademicCalendar, error)
	ListAll(ctx context.Context) ([]models.AcademicCalendar, error)
	Delete(ctx context.Context, id uint) error
}
