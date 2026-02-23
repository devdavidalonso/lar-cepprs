package postgres

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository"
)

// userRepository implements the repository.UserRepository interface for PostgreSQL
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of repository.UserRepository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Returns nil without error when not found
		}
		return nil, fmt.Errorf("error finding user by ID: %w", result.Error)
	}

	return &user, nil
}

// Preload Address and UserContacts for GetProfessorByID inside the service
func (r *userRepository) FindByIDWithAssociations(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).
		Preload("Address").
		Preload("UserContacts").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding user by ID: %w", result.Error)
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	// Criar um novo contexto com timeout maior (30 segundos)
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result := r.db.WithContext(timeoutCtx).
		Where("LOWER(email) = LOWER(?) AND deleted_at IS NULL", email).
		First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Returns nil without error when not found
		}
		return nil, fmt.Errorf("error finding user by email: %w", result.Error)
	}

	return &user, nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "unique constraint") {
			if strings.Contains(result.Error.Error(), "email") {
				return fmt.Errorf("a user with this email already exists")
			}
			if strings.Contains(result.Error.Error(), "cpf") {
				return fmt.Errorf("a user with this CPF already exists")
			}
			return fmt.Errorf("uniqueness violation: %w", result.Error)
		}
		return fmt.Errorf("error creating user: %w", result.Error)
	}

	// Compatibilidade com schema legado: coluna users.profile (texto) obrigatória.
	// Mantemos profile_id como fonte principal e sincronizamos profile textual.
	profileText := profileTextFromID(user.ProfileID)
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("profile", profileText).Error; err != nil {
		return fmt.Errorf("error syncing legacy profile text: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Model(user).Updates(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "unique constraint") {
			if strings.Contains(result.Error.Error(), "email") {
				return fmt.Errorf("another user with this email already exists")
			}
			if strings.Contains(result.Error.Error(), "cpf") {
				return fmt.Errorf("another user with this CPF already exists")
			}
			return fmt.Errorf("uniqueness violation: %w", result.Error)
		}
		return fmt.Errorf("error updating user: %w", result.Error)
	}

	// Sync coluna legado users.profile para evitar inconsistencias em bancos antigos.
	profileText := profileTextFromID(user.ProfileID)
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Update("profile", profileText).Error; err != nil {
		return fmt.Errorf("error syncing legacy profile text: %w", err)
	}

	return nil
}

// UpdateWithAssociations updates an existing user and replaces Address and UserContacts
func (r *userRepository) UpdateWithAssociations(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Atualiza campos normais do usuario
		if err := tx.Model(user).Updates(user).Error; err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return fmt.Errorf("uniqueness violation: %w", err)
			}
			return fmt.Errorf("error updating user: %w", err)
		}

		// 2. Substitui ou atualiza o Endereço
		if user.Address != nil {
			if err := tx.Model(user).Association("Address").Replace(user.Address); err != nil {
				return fmt.Errorf("error replacing address: %w", err)
			}
		} else {
			// Opcional: remover endereco se nil
			tx.Model(user).Association("Address").Clear()
		}

		// 3. Substitui os Contatos de Emergência
		if user.UserContacts != nil {
			if err := tx.Model(user).Association("UserContacts").Replace(user.UserContacts); err != nil {
				return fmt.Errorf("error replacing contacts: %w", err)
			}
		}

		return nil
	})
}

// Delete removes a user (soft delete)
func (r *userRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("error deleting user: %w", result.Error)
	}

	return nil
}

// GetUserProfiles gets the profiles of a user (legacy - now returns single profile based on ProfileID)
func (r *userRepository) GetUserProfiles(ctx context.Context, userID uint) ([]models.UserProfile, error) {
	var profiles []models.UserProfile

	// First get the user to find their ProfileID
	var user models.User
	result := r.db.WithContext(ctx).Select("profile_id").Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return profiles, nil // Return empty if user not found
		}
		return nil, fmt.Errorf("error finding user: %w", result.Error)
	}

	// Get the profile
	var profile models.UserProfile
	result = r.db.WithContext(ctx).Where("id = ?", user.ProfileID).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return profiles, nil // Return empty if profile not found
		}
		return nil, fmt.Errorf("error finding profile: %w", result.Error)
	}

	profiles = append(profiles, profile)
	return profiles, nil
}

// FindProfileByID finds a profile by ID
func (r *userRepository) FindProfileByID(ctx context.Context, id uint) (*models.UserProfile, error) {
	var profile models.UserProfile

	result := r.db.WithContext(ctx).Where("id = ?", id).First(&profile)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding profile by ID: %w", result.Error)
	}

	return &profile, nil
}

// UpdateLastLogin updates the last login date
func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint, timestamp time.Time) error {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login", timestamp)

	if result.Error != nil {
		return fmt.Errorf("error updating last login: %w", result.Error)
	}

	return nil
}

// FindByProfile finds users by profile name (legacy - maps name to ID)
func (r *userRepository) FindByProfile(ctx context.Context, profile string) ([]models.User, error) {
	// Map profile name to ID
	var profileID uint
	switch profile {
	case "admin", "administrator":
		profileID = 1
	case "teacher", "professor":
		profileID = 2
	case "student", "aluno":
		profileID = 3
	default:
		profileID = 3 // Default to student
	}

	return r.FindByProfileID(ctx, profileID)
}

func profileTextFromID(profileID uint) string {
	switch profileID {
	case 1:
		return "administrator"
	case 2:
		return "teacher"
	case 3:
		return "student"
	default:
		return "student"
	}
}

// FindByProfileID finds users by profile ID
func (r *userRepository) FindByProfileID(ctx context.Context, profileID uint) ([]models.User, error) {
	var users []models.User

	result := r.db.WithContext(ctx).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("error finding users by profile ID: %w", result.Error)
	}

	return users, nil
}

// UpsertTeacherProfile creates/updates teacher metadata linked to a user.
func (r *userRepository) UpsertTeacherProfile(ctx context.Context, userID uint, specialization, bio, phone string, active bool) error {
	teacher, err := r.GetTeacherProfileByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if teacher == nil {
		newTeacher := models.Teacher{
			UserID:         userID,
			Specialization: specialization,
			Bio:            bio,
			Phone:          phone,
			Active:         active,
		}
		if err := r.db.WithContext(ctx).Create(&newTeacher).Error; err != nil {
			return fmt.Errorf("error creating teacher profile: %w", err)
		}
		return nil
	}

	updates := map[string]interface{}{
		"specialization": specialization,
		"bio":            bio,
		"phone":          phone,
		"active":         active,
	}
	if err := r.db.WithContext(ctx).Model(&models.Teacher{}).
		Where("id = ?", teacher.ID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("error updating teacher profile: %w", err)
	}

	return nil
}

// GetTeacherProfileByUserID returns teacher row for a user.
func (r *userRepository) GetTeacherProfileByUserID(ctx context.Context, userID uint) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&teacher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("error finding teacher profile: %w", err)
	}
	return &teacher, nil
}

// ReplaceTeacherProgramsByUserID replaces all teacher-program links for a given user.
func (r *userRepository) ReplaceTeacherProgramsByUserID(ctx context.Context, userID uint, programIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var teacher models.Teacher
		err := tx.Where("user_id = ?", userID).First(&teacher).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error finding teacher profile for programs: %w", err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newTeacher := models.Teacher{
				UserID: userID,
				Active: true,
			}
			if err := tx.Create(&newTeacher).Error; err != nil {
				return fmt.Errorf("error creating teacher profile for programs: %w", err)
			}
			teacher = newTeacher
		}

		if err := tx.Where("teacher_id = ?", teacher.ID).Delete(&models.TeacherProgram{}).Error; err != nil {
			return fmt.Errorf("error clearing teacher programs: %w", err)
		}

		if len(programIDs) == 0 {
			return nil
		}

		unique := make(map[uint]struct{}, len(programIDs))
		for _, id := range programIDs {
			if id > 0 {
				unique[id] = struct{}{}
			}
		}
		ids := make([]uint, 0, len(unique))
		for id := range unique {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

		links := make([]models.TeacherProgram, 0, len(ids))
		for _, programID := range ids {
			links = append(links, models.TeacherProgram{
				TeacherID: teacher.ID,
				ProgramID: programID,
				Role:      "teacher",
				IsActive:  true,
			})
		}

		if err := tx.Create(&links).Error; err != nil {
			return fmt.Errorf("error creating teacher programs: %w", err)
		}

		return nil
	})
}

// GetTeacherProgramIDsByUserID lists all linked programs for the teacher represented by userID.
func (r *userRepository) GetTeacherProgramIDsByUserID(ctx context.Context, userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.WithContext(ctx).
		Table("teacher_programs tp").
		Select("tp.program_id").
		Joins("JOIN teachers t ON t.id = tp.teacher_id").
		Where("t.user_id = ?", userID).
		Order("tp.program_id ASC").
		Scan(&ids).Error
	if err != nil {
		return nil, fmt.Errorf("error listing teacher programs: %w", err)
	}
	return ids, nil
}
