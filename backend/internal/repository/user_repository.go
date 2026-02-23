package repository

import (
	"context"
	"time"

	"github.com/devdavidalonso/cecor/backend/internal/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uint) (*models.User, error)

	// FindByIDWithAssociations finds a user by ID and eager loads its Address and UserContacts
	FindByIDWithAssociations(ctx context.Context, id uint) (*models.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *models.User) error

	// UpdateWithAssociations updates an existing user and its Address/Contacts
	UpdateWithAssociations(ctx context.Context, user *models.User) error

	// Delete performs a logical deletion of a user
	Delete(ctx context.Context, id uint) error

	// GetUserProfiles gets the profiles of a user (legacy - returns profile by user ID)
	GetUserProfiles(ctx context.Context, userID uint) ([]models.UserProfile, error)

	// FindProfileByID finds a profile by ID
	FindProfileByID(ctx context.Context, id uint) (*models.UserProfile, error)

	// UpdateLastLogin updates the last login date
	UpdateLastLogin(ctx context.Context, id uint, timestamp time.Time) error

	// FindByProfile finds users by profile name (legacy)
	FindByProfile(ctx context.Context, profile string) ([]models.User, error)

	// FindByProfileID finds users by profile ID
	FindByProfileID(ctx context.Context, profileID uint) ([]models.User, error)

	// UpsertTeacherProfile creates/updates teacher record linked to a user
	UpsertTeacherProfile(ctx context.Context, userID uint, specialization, bio, phone string, active bool) error

	// GetTeacherProfileByUserID returns teacher profile row linked to a user
	GetTeacherProfileByUserID(ctx context.Context, userID uint) (*models.Teacher, error)

	// ReplaceTeacherProgramsByUserID replaces all program links for a teacher (via user_id)
	ReplaceTeacherProgramsByUserID(ctx context.Context, userID uint, programIDs []uint) error

	// GetTeacherProgramIDsByUserID returns program IDs linked to a teacher (via user_id)
	GetTeacherProgramIDsByUserID(ctx context.Context, userID uint) ([]uint, error)
}
