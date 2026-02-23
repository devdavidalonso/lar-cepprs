package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository"
)

type programRepository struct {
	db *gorm.DB
}

// NewProgramRepository creates a new instance of repository.ProgramRepository
func NewProgramRepository(db *gorm.DB) repository.ProgramRepository {
	return &programRepository{db: db}
}

func (r *programRepository) FindByID(ctx context.Context, id uint) (*models.Program, error) {
	var program models.Program
	if err := r.db.WithContext(ctx).First(&program, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find program by ID: %w", err)
	}
	return &program, nil
}

func (r *programRepository) ListAll(ctx context.Context, activeOnly bool) ([]models.Program, error) {
	var programs []models.Program
	query := r.db.WithContext(ctx).Model(&models.Program{})
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	if err := query.Find(&programs).Error; err != nil {
		return nil, fmt.Errorf("failed to list programs: %w", err)
	}
	return programs, nil
}

func (r *programRepository) FindByIDs(ctx context.Context, ids []uint) ([]models.Program, error) {
	var programs []models.Program
	if len(ids) == 0 {
		return programs, nil
	}
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&programs).Error; err != nil {
		return nil, fmt.Errorf("failed to find programs by IDs: %w", err)
	}
	return programs, nil
}
