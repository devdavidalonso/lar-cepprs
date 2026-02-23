package teachers

import (
	"context"
	"fmt"
	"strings"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository"
	"github.com/devdavidalonso/cecor/backend/internal/service/email"
	"github.com/devdavidalonso/cecor/backend/internal/service/keycloak"
)

// Service defines the teacher service interface.
type Service interface {
	CreateTeacher(ctx context.Context, teacher *models.User) error
	GetTeachers(ctx context.Context) ([]models.User, error)
	GetTeacherByID(ctx context.Context, id uint) (*models.User, error)
	UpdateTeacher(ctx context.Context, teacher *models.User) error
	DeleteTeacher(ctx context.Context, id uint) error

	// Deprecated: compatibility aliases kept for gradual migration.
	CreateProfessor(ctx context.Context, teacher *models.User) error
	GetProfessors(ctx context.Context) ([]models.User, error)
	GetProfessorByID(ctx context.Context, id uint) (*models.User, error)
	UpdateProfessor(ctx context.Context, teacher *models.User) error
	DeleteProfessor(ctx context.Context, id uint) error
}

// teacherService implements the Service interface.
type teacherService struct {
	userRepo    repository.UserRepository
	programRepo repository.ProgramRepository
	keycloak    *keycloak.KeycloakService
	email       *email.EmailService
}

// NewService creates a new instance of teacherService.
func NewService(userRepo repository.UserRepository, programRepo repository.ProgramRepository, keycloak *keycloak.KeycloakService, email *email.EmailService) Service {
	return &teacherService{
		userRepo:    userRepo,
		programRepo: programRepo,
		keycloak:    keycloak,
		email:       email,
	}
}

// CreateTeacher creates a new teacher.
func (s *teacherService) CreateTeacher(ctx context.Context, teacher *models.User) error {
	// Validate required fields
	if teacher.Name == "" || teacher.Email == "" {
		return fmt.Errorf("name and email are required")
	}

	// Set profile as 'professor' (ProfileID = 2)
	teacher.ProfileID = 2
	teacher.Active = true

	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(ctx, teacher.Email)
	if err != nil {
		return fmt.Errorf("error checking existing email: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("a user with this email already exists")
	}

	// Create user in database
	// Set placeholder password
	teacher.Password = "temp123456"

	if err := s.userRepo.Create(ctx, teacher); err != nil {
		return fmt.Errorf("error creating professor in database: %w", err)
	}

	if err := s.userRepo.UpsertTeacherProfile(
		ctx,
		teacher.ID,
		teacher.Specialization,
		teacher.Bio,
		teacher.Phone,
		teacher.Active,
	); err != nil {
		return fmt.Errorf("error creating teacher profile metadata: %w", err)
	}

	if teacher.ProgramIDs != nil {
		if err := s.userRepo.ReplaceTeacherProgramsByUserID(ctx, teacher.ID, teacher.ProgramIDs); err != nil {
			return fmt.Errorf("error linking teacher to programs: %w", err)
		}
	}

	// Create in Keycloak
	if s.keycloak != nil {
		// Generate temporary password
		tempPassword := "prof123" // Or generate random

		// Create user
		// Split name
		nameParts := strings.Fields(teacher.Name)
		firstName := nameParts[0]
		lastName := ""
		if len(nameParts) > 1 {
			lastName = strings.Join(nameParts[1:], " ")
		}

		req := keycloak.CreateUserRequest{
			Username:      teacher.Email,
			Email:         teacher.Email,
			FirstName:     firstName,
			LastName:      lastName,
			Enabled:       true,
			EmailVerified: true, // Auto-verify for simplicity
		}

		keycloakID, err := s.keycloak.CreateUser(ctx, req)
		if err != nil {
			fmt.Printf("Warning: failed to create Keycloak user: %v\n", err)
		} else {
			// Assign role
			if err := s.keycloak.AssignRole(ctx, keycloakID, "teacher"); err != nil {
				fmt.Printf("Warning: failed to assign role: %v\n", err)
			}

			// Dynamic Group Assignment
			if len(teacher.ProgramIDs) > 0 {
				programs, err := s.programRepo.FindByIDs(ctx, teacher.ProgramIDs)
				if err != nil {
					fmt.Printf("Warning: failed to fetch programs for group assignment: %v\n", err)
				} else {
					for _, p := range programs {
						if err := s.keycloak.AddUserToGroup(ctx, keycloakID, p.Code); err != nil {
							fmt.Printf("Warning: failed to add user to group '%s': %v\n", err)
						}
					}
				}
			}

			// Set password
			if err := s.keycloak.SetTemporaryPassword(ctx, keycloakID, tempPassword); err != nil {
				fmt.Printf("Warning: failed to set password: %v\n", err)
			}

			// Update user with Keycloak ID
			teacher.KeycloakUserID = &keycloakID
			s.userRepo.Update(ctx, teacher)

			// Send email
			if s.email != nil {
				s.email.SendWelcomeEmail(teacher.Email, teacher.Name, tempPassword)
			}
		}
	}

	return nil
}

// GetTeachers returns all teachers.
func (s *teacherService) GetTeachers(ctx context.Context) ([]models.User, error) {
	professors, err := s.userRepo.FindByProfileID(ctx, 2)
	if err != nil {
		return nil, err
	}

	for i := range professors {
		if err := s.hydrateTeacherMetadata(ctx, &professors[i]); err != nil {
			return nil, err
		}
	}

	return professors, nil
}

// GetTeacherByID returns a teacher by ID.
func (s *teacherService) GetTeacherByID(ctx context.Context, id uint) (*models.User, error) {
	// Call FindByIDWithAssociations instead of FindByID to load contacts and address
	user, err := s.userRepo.FindByIDWithAssociations(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("professor not found")
	}
	if user.ProfileID != 2 {
		return nil, fmt.Errorf("user is not a professor")
	}
	if err := s.hydrateTeacherMetadata(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateTeacher updates a teacher.
func (s *teacherService) UpdateTeacher(ctx context.Context, teacher *models.User) error {
	existing, err := s.GetTeacherByID(ctx, teacher.ID)
	if err != nil {
		return err
	}

	// Update allowed fields
	existing.Name = teacher.Name
	existing.Phone = teacher.Phone
	existing.CPF = teacher.CPF
	existing.Specialization = teacher.Specialization
	existing.Bio = teacher.Bio
	existing.LinkedinURL = teacher.LinkedinURL
	existing.ProgramIDs = teacher.ProgramIDs

	// Include associations payload
	existing.Address = teacher.Address
	existing.UserContacts = teacher.UserContacts

	if err := s.userRepo.UpdateWithAssociations(ctx, existing); err != nil {
		return err
	}

	if err := s.userRepo.UpsertTeacherProfile(
		ctx,
		existing.ID,
		existing.Specialization,
		existing.Bio,
		existing.Phone,
		existing.Active,
	); err != nil {
		return fmt.Errorf("error updating teacher profile metadata: %w", err)
	}

	if teacher.ProgramIDs != nil {
		if err := s.userRepo.ReplaceTeacherProgramsByUserID(ctx, existing.ID, teacher.ProgramIDs); err != nil {
			return fmt.Errorf("error updating teacher programs: %w", err)
		}
	}

	return nil
}

// DeleteTeacher deletes a teacher.
func (s *teacherService) DeleteTeacher(ctx context.Context, id uint) error {
	// Check if exists
	_, err := s.GetTeacherByID(ctx, id)
	if err != nil {
		return err
	}

	// Disable in Keycloak if needed (omitted for brevity)

	return s.userRepo.Delete(ctx, id)
}

func (s *teacherService) hydrateTeacherMetadata(ctx context.Context, user *models.User) error {
	teacherProfile, err := s.userRepo.GetTeacherProfileByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error loading teacher profile metadata: %w", err)
	}
	if teacherProfile != nil {
		user.Specialization = teacherProfile.Specialization
		user.Bio = teacherProfile.Bio
		if user.Phone == "" {
			user.Phone = teacherProfile.Phone
		}
		user.Active = teacherProfile.Active
	}

	programIDs, err := s.userRepo.GetTeacherProgramIDsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error loading teacher program links: %w", err)
	}
	user.ProgramIDs = programIDs

	return nil
}

// Deprecated: use CreateTeacher.
func (s *teacherService) CreateProfessor(ctx context.Context, teacher *models.User) error {
	return s.CreateTeacher(ctx, teacher)
}

// Deprecated: use GetTeachers.
func (s *teacherService) GetProfessors(ctx context.Context) ([]models.User, error) {
	return s.GetTeachers(ctx)
}

// Deprecated: use GetTeacherByID.
func (s *teacherService) GetProfessorByID(ctx context.Context, id uint) (*models.User, error) {
	return s.GetTeacherByID(ctx, id)
}

// Deprecated: use UpdateTeacher.
func (s *teacherService) UpdateProfessor(ctx context.Context, teacher *models.User) error {
	return s.UpdateTeacher(ctx, teacher)
}

// Deprecated: use DeleteTeacher.
func (s *teacherService) DeleteProfessor(ctx context.Context, id uint) error {
	return s.DeleteTeacher(ctx, id)
}
