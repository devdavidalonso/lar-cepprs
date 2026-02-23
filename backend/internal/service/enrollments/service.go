package enrollments

import (
	"context"
	"fmt"
	"time"

	"github.com/devdavidalonso/cecor/backend/internal/infrastructure/googleapis"
	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository"
)

type Service interface {
	EnrollStudent(ctx context.Context, enrollment *models.Enrollment) error
	GetEnrollment(ctx context.Context, id uint) (*models.Enrollment, error)
	ListEnrollments(ctx context.Context) ([]models.Enrollment, error)
	ListByCourse(ctx context.Context, courseID uint) ([]models.Enrollment, error)
	UpdateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	DeleteEnrollment(ctx context.Context, id uint) error
}

type service struct {
	repo            repository.EnrollmentRepository
	studentRepo     repository.StudentRepository
	courseRepo      repository.CourseRepository
	classroomClient googleapis.GoogleClassroomClient
}

func NewService(repo repository.EnrollmentRepository, studentRepo repository.StudentRepository, courseRepo repository.CourseRepository, classroomClient googleapis.GoogleClassroomClient) Service {
	return &service{
		repo:            repo,
		studentRepo:     studentRepo,
		courseRepo:      courseRepo,
		classroomClient: classroomClient,
	}
}

func (s *service) EnrollStudent(ctx context.Context, enrollment *models.Enrollment) error {
	// 1. Check if student is already enrolled in this course
	existing, err := s.repo.FindByStudentAndCourse(ctx, enrollment.StudentID, enrollment.CourseID)
	if err != nil {
		return fmt.Errorf("failed to check existing enrollment: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("student is already enrolled in this course")
	}

	// 1.1 Validate Student Age (Must be >= 12)
	student, err := s.studentRepo.FindByID(ctx, enrollment.StudentID)
	if err != nil {
		return fmt.Errorf("failed to fetch student: %w", err)
	}
	// Calculate age
	now := time.Now()
	age := now.Year() - student.User.BirthDate.Year()
	if now.YearDay() < student.User.BirthDate.YearDay() {
		age--
	}
	if age < 12 {
		return fmt.Errorf("student must be at least 12 years old (current age: %d)", age)
	}

	// 1.2 Validate Schedule Conflict
	course, err := s.courseRepo.FindByID(ctx, enrollment.CourseID)
	if err != nil {
		return fmt.Errorf("failed to fetch course: %w", err)
	}

	// Get all active enrollments for the student
	activeEnrollments, err := s.repo.FindAll(ctx) // TODO: Create FindByStudent(ctx, studentID) for better performance
	if err != nil {
		return fmt.Errorf("failed to fetch student enrollments: %w", err)
	}

	for _, active := range activeEnrollments {
		if active.StudentID != enrollment.StudentID || active.Status != "active" {
			continue
		}

		activeCourse, err := s.courseRepo.FindByID(ctx, active.CourseID)
		if err != nil {
			continue
		}

		// Check overlap: Same WeekDay and Overlapping Time
		if activeCourse.WeekDays == course.WeekDays {
			// Simple logic: if start times match or overlap (using pure string comparison for MVP "HH:mm")
			// In production, parse Request "10:00" matches Existing "10:00"
			if activeCourse.StartTime == course.StartTime {
				return fmt.Errorf("schedule conflict: student already has a class at %s on %s", course.StartTime, course.WeekDays)
			}
		}
	}

	// 1.3 Validate course capacity (max students)
	if course.MaxStudents > 0 {
		courseEnrollments, err := s.repo.ListByCourse(ctx, enrollment.CourseID)
		if err != nil {
			return fmt.Errorf("failed to validate course capacity: %w", err)
		}

		activeCount := 0
		for _, e := range courseEnrollments {
			if e.Status == "active" {
				activeCount++
			}
		}

		if activeCount >= course.MaxStudents {
			return fmt.Errorf(
				"course capacity reached: %d/%d active enrollments",
				activeCount,
				course.MaxStudents,
			)
		}
	}

	// 2. Generate Enrollment Number if not provided
	if enrollment.EnrollmentNumber == "" {
		enrollment.EnrollmentNumber = fmt.Sprintf("MAT-%d-%d", enrollment.StudentID, time.Now().Unix())
	}

	// 3. Set default dates if not provided
	if enrollment.EnrollmentDate.IsZero() {
		enrollment.EnrollmentDate = time.Now()
	}
	if enrollment.StartDate.IsZero() {
		enrollment.StartDate = time.Now()
	}
	if enrollment.Status == "" {
		enrollment.Status = "active"
	}

	// 4. Salvar na Database local
	if err := s.repo.Create(ctx, enrollment); err != nil {
		return err
	}

	// 5. Convidar aluno para a Turma no Google Classroom
	if s.classroomClient != nil && course.GoogleClassroomID != "" && student.User.Email != "" {
		fmt.Printf("Adding student %s to Google Classroom Course %s...\n", student.User.Email, course.GoogleClassroomID)
		_, err := s.classroomClient.AddStudent(course.GoogleClassroomID, student.User.Email)
		if err != nil {
			fmt.Printf("Warning: Failed to add student to Google Classroom: %v\n", err)
			// Decide if sync fail blocks. Usually not.
		} else {
			fmt.Println("Student added successfully to Google Classroom!")
		}
	}

	return nil
}

func (s *service) GetEnrollment(ctx context.Context, id uint) (*models.Enrollment, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) ListEnrollments(ctx context.Context) ([]models.Enrollment, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) ListByCourse(ctx context.Context, courseID uint) ([]models.Enrollment, error) {
	return s.repo.ListByCourse(ctx, courseID)
}

func (s *service) UpdateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	return s.repo.Update(ctx, enrollment)
}

func (s *service) DeleteEnrollment(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
