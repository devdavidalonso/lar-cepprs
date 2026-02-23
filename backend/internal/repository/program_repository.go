package repository

import (
	"context"

	"github.com/devdavidalonso/cecor/backend/internal/models"
)

// ProgramRepository defines the interface for program data access operations
type ProgramRepository interface {
	// FindByID finds a program by ID
	FindByID(ctx context.Context, id uint) (*models.Program, error)

	// ListAll returns all programs
	ListAll(ctx context.Context, activeOnly bool) ([]models.Program, error)

	// FindByIDs returns multiple programs by their IDs
	FindByIDs(ctx context.Context, ids []uint) ([]models.Program, error)
}
