// backend/internal/service/users/user_service.go
package users

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// userService implements the Service interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new instance of the user service
func NewUserService(userRepo repository.UserRepository) Service {
	return &userService{
		userRepo: userRepo,
	}
}

// Authenticate authenticates a user with email and password
func (s *userService) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	log.Printf("Debugging email: %s", email) // Linha adicionada para debugging
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.Active {
		return nil, fmt.Errorf("user is inactive")
	}

	// Compare password
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	// 	return nil, fmt.Errorf("incorrect password")
	// }

	// Load user profile relationship
	profile, err := s.userRepo.FindProfileByID(ctx, user.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("error loading profile: %w", err)
	}
	user.Profile = *profile

	return user, nil
}

// GetUserByID gets a user by ID
func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Load user profile relationship
	profile, err := s.userRepo.FindProfileByID(ctx, user.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("error loading profile: %w", err)
	}
	user.Profile = *profile

	return user, nil
}

// GetUserByEmail gets a user by email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Load user profile relationship
	profile, err := s.userRepo.FindProfileByID(ctx, user.ProfileID)
	if err != nil {
		return nil, fmt.Errorf("error loading profile: %w", err)
	}
	user.Profile = *profile

	return user, nil
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	// Check required fields
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("name, email, and password are required")
	}

	// Check if a user with the same email already exists
	existing, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("error checking existing email: %w", err)
	}

	if existing != nil {
		return fmt.Errorf("a user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error generating password hash: %w", err)
	}

	user.Password = string(hashedPassword)

	// Set default status
	if !user.Active {
		user.Active = true
	}

	// Create user
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	// Check if user exists
	existing, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("user not found")
	}

	// Check email change
	if user.Email != existing.Email {
		emailExists, err := s.userRepo.FindByEmail(ctx, user.Email)
		if err != nil {
			return fmt.Errorf("error checking existing email: %w", err)
		}

		if emailExists != nil && emailExists.ID != user.ID {
			return fmt.Errorf("another user with this email already exists")
		}
	}

	// Check password change
	if user.Password != "" && user.Password != existing.Password {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error generating password hash: %w", err)
		}

		user.Password = string(hashedPassword)
	} else {
		user.Password = existing.Password
	}

	// Update user
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	// Check if user exists
	existing, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	if existing == nil {
		return fmt.Errorf("user not found")
	}

	// Delete user (soft delete)
	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp
func (s *userService) UpdateLastLogin(ctx context.Context, id uint) error {
	now := time.Now()
	err := s.userRepo.UpdateLastLogin(ctx, id, now)
	if err != nil {
		return fmt.Errorf("error updating last login: %w", err)
	}

	return nil
}

// Profile name to ID mapping
const (
	ProfileAdmin     = 1
	ProfileProfessor = 2
	ProfileStudent   = 3
)

// FindOrCreateByEmail finds a user by email or creates a new one
func (s *userService) FindOrCreateByEmail(ctx context.Context, email, name, profile string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	// Map Keycloak profiles to local profile IDs
	var localProfileID uint
	switch strings.ToLower(strings.TrimSpace(profile)) {
	case "administrator", "admin":
		localProfileID = ProfileAdmin
	case "gestor":
		localProfileID = ProfileAdmin // Manager treated as admin for now
	case "teacher", "professor":
		localProfileID = ProfileProfessor
	case "student", "aluno":
		localProfileID = ProfileStudent
	case "responsável", "responsavel":
		localProfileID = ProfileStudent // Guardian treated as student for now
	default:
		localProfileID = ProfileStudent
	}

	if user != nil {
		// Update user info if changed
		if user.Name != name || user.ProfileID != localProfileID {
			user.Name = name
			user.ProfileID = localProfileID
			if err := s.userRepo.Update(ctx, user); err != nil {
				return nil, fmt.Errorf("error updating user info from SSO: %w", err)
			}
		}
		return user, nil
	}

	// Create new user
	newUser := &models.User{
		Email:     email,
		Name:      name,
		Password:  "", // No password for SSO users
		ProfileID: localProfileID,
		Active:    true,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return newUser, nil
}
