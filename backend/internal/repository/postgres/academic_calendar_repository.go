package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository"
)

type academicCalendarRepository struct {
	db *gorm.DB
}

// NewAcademicCalendarRepository creates a new instance of repository.AcademicCalendarRepository
func NewAcademicCalendarRepository(db *gorm.DB) repository.AcademicCalendarRepository {
	return &academicCalendarRepository{db: db}
}

func (r *academicCalendarRepository) Create(ctx context.Context, event *models.AcademicCalendar) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *academicCalendarRepository) FindRecessesBetween(ctx context.Context, programID *uint, start, end time.Time) ([]models.AcademicCalendar, error) {
	var events []models.AcademicCalendar
	query := r.db.WithContext(ctx).
		Where("((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND end_date >= ?))", start, end, start, end, start, end).
		Where("is_active = ?", true).
		Where("type = ?", "recess")

	if programID != nil {
		query = query.Where("(program_id = ? OR program_id IS NULL)", *programID)
	} else {
		query = query.Where("program_id IS NULL")
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *academicCalendarRepository) ListAll(ctx context.Context) ([]models.AcademicCalendar, error) {
	var events []models.AcademicCalendar
	if err := r.db.WithContext(ctx).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *academicCalendarRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.AcademicCalendar{}, id).Error
}
