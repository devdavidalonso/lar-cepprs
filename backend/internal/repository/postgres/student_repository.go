package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/devdavidalonso/cecor/backend/internal/models"
	"github.com/devdavidalonso/cecor/backend/internal/repository"
)

// studentRepository implements the StudentRepository interface using PostgreSQL/GORM
type studentRepository struct {
	db *gorm.DB
}

// NewStudentRepository creates a new implementation of StudentRepository
func NewStudentRepository(db *gorm.DB) repository.StudentRepository {
	return &studentRepository{
		db: db,
	}
}

// FindAll returns all students with pagination and filters
// Parameters:
// - ctx: context for database operations
// - page: page number (starting from 1)
// - pageSize: number of items per page
// - filters: map of filter criteria (name, email, cpf, status, min_age, max_age, course_id)
// Returns:
// - []models.Student: list of student records
// - int64: total count of matching records (before pagination)
// - error: any error encountered during the operation
func (r *studentRepository) FindAll(ctx context.Context, page int, pageSize int, filters map[string]interface{}) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	// Start query with soft delete excluded
	query := r.db.WithContext(ctx).Model(&models.Student{}).Where("students.deleted_at IS NULL")

	// Apply filters
	for key, value := range filters {
		switch key {
		case "name":
			query = query.Joins("JOIN users ON users.id = students.user_id").
				Where("users.name ILIKE ?", fmt.Sprintf("%%%s%%", value))
		case "email":
			// Email is in the users table, so we need to join
			query = query.Joins("JOIN users ON users.id = students.user_id").
				Where("users.email ILIKE ?", fmt.Sprintf("%%%s%%", value))
		case "cpf":
			// CPF is also in the users table
			query = query.Joins("JOIN users ON users.id = students.user_id").
				Where("users.cpf = ?", value)
		case "status":
			query = query.Where("students.status = ?", value)
		case "min_age":
			// Filter by minimum age using birth date
			query = query.Joins("JOIN users ON users.id = students.user_id").
				Where("users.birth_date <= CURRENT_DATE - INTERVAL '? year'", value)
		case "max_age":
			// Filter by maximum age using birth date
			query = query.Joins("JOIN users ON users.id = students.user_id").
				Where("users.birth_date >= CURRENT_DATE - INTERVAL '? year'", value)
		case "course_id":
			// Filter by course requires a join with enrollments
			query = query.Joins("JOIN enrollments ON enrollments.student_id = students.id").
				Where("enrollments.course_id = ? AND enrollments.status IN ('active', 'in_progress')", value)
		case "program_id":
			query = query.Joins("JOIN student_programs sp ON sp.student_id = students.id").
				Where("sp.program_id = ? AND sp.status = 'active'", value)
		}
	}

	// Count total records (for pagination)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting students: %w", err)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Order by name (default)
	query = query.Joins("JOIN users ON users.id = students.user_id").
		Order("users.name ASC")

	// Include associated user data
	query = query.Preload("User")

	// Execute query
	if err := query.Find(&students).Error; err != nil {
		return nil, 0, fmt.Errorf("error fetching students: %w", err)
	}

	// Calculate ages
	for i := range students {
		students[i].User.CalculateAge()
	}

	if err := r.hydrateStudentsProgramIDs(ctx, students); err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

// FindByID finds a student by ID
// Parameters:
// - ctx: context for database operations
// - id: student ID to find
// Returns:
// - *models.Student: pointer to student record (nil if not found)
// - error: any error encountered during the operation
func (r *studentRepository) FindByID(ctx context.Context, id uint) (*models.Student, error) {
	var student models.Student

	result := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&student)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Returns nil without error when not found
		}
		return nil, fmt.Errorf("error finding student by ID: %w", result.Error)
	}

	// Calculate age
	student.User.CalculateAge()

	programIDs, err := r.getStudentProgramIDsByStudentID(ctx, student.ID)
	if err != nil {
		return nil, err
	}
	student.ProgramIDs = programIDs

	return &student, nil
}

// FindByEmail finds a student by email
// Parameters:
// - ctx: context for database operations
// - email: email address to search for
// Returns:
// - *models.Student: pointer to student record (nil if not found)
// - error: any error encountered during the operation
func (r *studentRepository) FindByEmail(ctx context.Context, email string) (*models.Student, error) {
	var student models.Student

	result := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = students.user_id").
		Where("LOWER(users.email) = LOWER(?) AND students.deleted_at IS NULL", email).
		Preload("User").
		First(&student)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Returns nil without error when not found
		}
		return nil, fmt.Errorf("error finding student by email: %w", result.Error)
	}

	// Calculate age
	student.User.CalculateAge()

	programIDs, err := r.getStudentProgramIDsByStudentID(ctx, student.ID)
	if err != nil {
		return nil, err
	}
	student.ProgramIDs = programIDs

	return &student, nil
}

// FindByCPF finds a student by CPF (Brazilian tax ID)
// Parameters:
// - ctx: context for database operations
// - cpf: CPF number (with or without formatting)
// Returns:
// - *models.Student: pointer to student record (nil if not found)
// - error: any error encountered during the operation
func (r *studentRepository) FindByCPF(ctx context.Context, cpf string) (*models.Student, error) {
	var student models.Student

	// Remove CPF formatting
	cleanCPF := strings.ReplaceAll(strings.ReplaceAll(cpf, ".", ""), "-", "")

	result := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = students.user_id").
		Where("users.cpf = ? AND students.deleted_at IS NULL", cleanCPF).
		Preload("User").
		First(&student)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Returns nil without error when not found
		}
		return nil, fmt.Errorf("error finding student by CPF: %w", result.Error)
	}

	// Calculate age
	student.User.CalculateAge()

	programIDs, err := r.getStudentProgramIDsByStudentID(ctx, student.ID)
	if err != nil {
		return nil, err
	}
	student.ProgramIDs = programIDs

	return &student, nil
}

// Create creates a new student and associated user
// Parameters:
// - ctx: context for database operations
// - student: student data to create (must include User data)
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) Create(ctx context.Context, student *models.Student) error {
	// Check if user data is provided
	if student.User.Name == "" && student.User.Email == "" && student.User.CPF == "" {
		return fmt.Errorf("user data is required")
	}

	// Check required fields
	if student.User.Name == "" || student.User.Email == "" || student.User.Phone == "" {
		return fmt.Errorf("required fields not filled")
	}

	// Clean CPF if provided
	if student.User.CPF != "" {
		student.User.CPF = strings.ReplaceAll(strings.ReplaceAll(student.User.CPF, ".", ""), "-", "")
	}

	// Start transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Set profile as 'student' (ProfileID = 3)
		student.User.ProfileID = 3

		// Create user first
		if err := tx.Create(&student.User).Error; err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				if strings.Contains(err.Error(), "email") {
					return fmt.Errorf("a user with this email already exists")
				}
				if strings.Contains(err.Error(), "cpf") {
					return fmt.Errorf("a user with this CPF already exists")
				}
				return fmt.Errorf("uniqueness violation: %w", err)
			}
			return fmt.Errorf("error creating user: %w", err)
		}

		// Associate user ID to student
		student.UserID = student.User.ID

		// Generate registration number if not provided
		if student.RegistrationNumber == "" {
			year := time.Now().Year()
			student.RegistrationNumber = fmt.Sprintf("%d%06d", year, student.User.ID)
		}

		// Set default status if not provided
		if student.Status == "" {
			student.Status = "active"
		}

		// Create student
		if err := tx.Create(student).Error; err != nil {
			return fmt.Errorf("error creating student: %w", err)
		}

		if err := r.replaceStudentProgramsTx(tx, student.ID, student.ProgramIDs); err != nil {
			return err
		}

		return nil
	})
}

// Update updates an existing student and user data
// Parameters:
// - ctx: context for database operations
// - student: student data to update
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) Update(ctx context.Context, student *models.Student) error {
	// Check if student exists
	existing, err := r.FindByID(ctx, student.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("student not found")
	}

	// Start transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// If user data is provided, update it
		// We check if any key field is present to assume user data is being updated
		if student.User.Name != "" || student.User.Email != "" || student.User.Phone != "" || student.User.CPF != "" {
			// Ensure user ID is correct
			student.User.ID = existing.UserID

			// Clean CPF if provided
			if student.User.CPF != "" {
				student.User.CPF = strings.ReplaceAll(strings.ReplaceAll(student.User.CPF, ".", ""), "-", "")
			}

			// Update user
			if err := tx.Model(&models.User{}).Where("id = ?", existing.UserID).Updates(student.User).Error; err != nil {
				if strings.Contains(err.Error(), "unique constraint") {
					if strings.Contains(err.Error(), "email") {
						return fmt.Errorf("another user with this email already exists")
					}
					if strings.Contains(err.Error(), "cpf") {
						return fmt.Errorf("another user with this CPF already exists")
					}
					return fmt.Errorf("uniqueness violation: %w", err)
				}
				return fmt.Errorf("error updating user: %w", err)
			}
		}

		// Keep original UserID and RegistrationNumber
		student.UserID = existing.UserID
		student.RegistrationNumber = existing.RegistrationNumber

		// Update student
		if err := tx.Model(student).Updates(student).Error; err != nil {
			return fmt.Errorf("error updating student: %w", err)
		}

		if student.ProgramIDs != nil {
			if err := r.replaceStudentProgramsTx(tx, student.ID, student.ProgramIDs); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete performs a logical deletion of a student
// Parameters:
// - ctx: context for database operations
// - id: student ID to delete
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) Delete(ctx context.Context, id uint) error {
	// Check if student exists
	existing, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("student not found")
	}

	// Soft delete (updates the deleted_at field)
	result := r.db.WithContext(ctx).Model(&models.Student{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()"))

	if result.Error != nil {
		return fmt.Errorf("error deleting student: %w", result.Error)
	}

	return nil
}

// GetGuardians returns a student's guardians
// Parameters:
// - ctx: context for database operations
// - studentID: ID of the student whose guardians to fetch
// Returns:
// - []models.Guardian: list of guardian records
// - error: any error encountered during the operation
func (r *studentRepository) GetGuardians(ctx context.Context, studentID uint) ([]models.Guardian, error) {
	var guardians []models.Guardian

	result := r.db.WithContext(ctx).
		Where("student_id = ? AND deleted_at IS NULL", studentID).
		Find(&guardians)

	if result.Error != nil {
		return nil, fmt.Errorf("error fetching guardians: %w", result.Error)
	}

	return guardians, nil
}

// AddGuardian adds a guardian to a student
// Parameters:
// - ctx: context for database operations
// - guardian: guardian data to add
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) AddGuardian(ctx context.Context, guardian *models.Guardian) error {
	// Check required fields
	if guardian.Name == "" || guardian.Relationship == "" {
		return fmt.Errorf("required fields not filled")
	}

	// Check if student exists
	student, err := r.FindByID(ctx, guardian.StudentID)
	if err != nil {
		return err
	}
	if student == nil {
		return fmt.Errorf("student not found")
	}

	// Clean CPF if provided
	if guardian.CPF != "" {
		guardian.CPF = strings.ReplaceAll(strings.ReplaceAll(guardian.CPF, ".", ""), "-", "")
	}

	// Insert guardian
	if err := r.db.WithContext(ctx).Create(guardian).Error; err != nil {
		return fmt.Errorf("error adding guardian: %w", err)
	}

	return nil
}

// UpdateGuardian updates an existing guardian
// Parameters:
// - ctx: context for database operations
// - guardian: guardian data to update
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) UpdateGuardian(ctx context.Context, guardian *models.Guardian) error {
	// Check if guardian exists
	var existing models.Guardian
	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", guardian.ID).
		First(&existing)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("guardian not found")
		}
		return fmt.Errorf("error fetching guardian: %w", result.Error)
	}

	// Clean CPF if provided
	if guardian.CPF != "" {
		guardian.CPF = strings.ReplaceAll(strings.ReplaceAll(guardian.CPF, ".", ""), "-", "")
	}

	// Update guardian
	if err := r.db.WithContext(ctx).Model(guardian).Updates(guardian).Error; err != nil {
		return fmt.Errorf("error updating guardian: %w", err)
	}

	return nil
}

// RemoveGuardian performs a logical deletion of a guardian
// Parameters:
// - ctx: context for database operations
// - guardianID: ID of the guardian to remove
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) RemoveGuardian(ctx context.Context, guardianID uint) error {
	// Check if guardian exists
	var existing models.Guardian
	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", guardianID).
		First(&existing)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("guardian not found")
		}
		return fmt.Errorf("error fetching guardian: %w", result.Error)
	}

	// Soft delete
	if err := r.db.WithContext(ctx).Model(&models.Guardian{}).
		Where("id = ?", guardianID).
		Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("error removing guardian: %w", err)
	}

	return nil
}

func (r *studentRepository) getStudentProgramIDsByStudentID(ctx context.Context, studentID uint) ([]uint, error) {
	var ids []uint
	if err := r.db.WithContext(ctx).
		Table("student_programs").
		Select("program_id").
		Where("student_id = ? AND status = 'active'", studentID).
		Order("program_id ASC").
		Scan(&ids).Error; err != nil {
		return nil, fmt.Errorf("error loading student programs: %w", err)
	}
	return ids, nil
}

func (r *studentRepository) hydrateStudentsProgramIDs(ctx context.Context, students []models.Student) error {
	if len(students) == 0 {
		return nil
	}

	studentIDs := make([]uint, 0, len(students))
	for _, s := range students {
		studentIDs = append(studentIDs, s.ID)
	}

	type studentProgramRow struct {
		StudentID uint
		ProgramID uint
	}

	var rows []studentProgramRow
	if err := r.db.WithContext(ctx).
		Table("student_programs").
		Select("student_id, program_id").
		Where("student_id IN ? AND status = 'active'", studentIDs).
		Order("program_id ASC").
		Scan(&rows).Error; err != nil {
		return fmt.Errorf("error loading student programs: %w", err)
	}

	byStudent := make(map[uint][]uint, len(studentIDs))
	for _, row := range rows {
		byStudent[row.StudentID] = append(byStudent[row.StudentID], row.ProgramID)
	}

	for i := range students {
		students[i].ProgramIDs = byStudent[students[i].ID]
	}

	return nil
}

func (r *studentRepository) replaceStudentProgramsTx(tx *gorm.DB, studentID uint, programIDs []uint) error {
	if err := tx.Where("student_id = ?", studentID).Delete(&models.StudentProgram{}).Error; err != nil {
		return fmt.Errorf("error clearing student programs: %w", err)
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

	links := make([]models.StudentProgram, 0, len(unique))
	for programID := range unique {
		links = append(links, models.StudentProgram{
			StudentID: studentID,
			ProgramID: programID,
			Status:    "active",
		})
	}

	if len(links) == 0 {
		return nil
	}

	if err := tx.Create(&links).Error; err != nil {
		return fmt.Errorf("error creating student programs: %w", err)
	}

	return nil
}

// GetDocuments returns a student's documents
// Parameters:
// - ctx: context for database operations
// - studentID: ID of the student whose documents to fetch
// Returns:
// - []models.Document: list of document records
// - error: any error encountered during the operation
func (r *studentRepository) GetDocuments(ctx context.Context, studentID uint) ([]models.Document, error) {
	var documents []models.Document

	result := r.db.WithContext(ctx).
		Where("student_id = ?", studentID).
		Find(&documents)

	if result.Error != nil {
		return nil, fmt.Errorf("error fetching documents: %w", result.Error)
	}

	return documents, nil
}

// AddDocument adds a document to a student
// Parameters:
// - ctx: context for database operations
// - document: document data to add
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) AddDocument(ctx context.Context, document *models.Document) error {
	// Check required fields
	if document.Name == "" || document.Type == "" || document.Path == "" {
		return fmt.Errorf("required fields not filled")
	}

	// Check if student exists
	student, err := r.FindByID(ctx, document.StudentID)
	if err != nil {
		return err
	}
	if student == nil {
		return fmt.Errorf("student not found")
	}

	// Set uploaded_by if not provided
	if document.UploadedByID == 0 {
		// In a real scenario, you would get this from the request context
		userID, ok := ctx.Value("user_id").(uint)
		if ok {
			document.UploadedByID = userID
		}
	}

	// Insert document
	if err := r.db.WithContext(ctx).Create(document).Error; err != nil {
		return fmt.Errorf("error adding document: %w", err)
	}

	return nil
}

// RemoveDocument removes a document
// Parameters:
// - ctx: context for database operations
// - documentID: ID of the document to remove
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) RemoveDocument(ctx context.Context, documentID uint) error {
	// Check if document exists
	var existing models.Document
	result := r.db.WithContext(ctx).
		Where("id = ?", documentID).
		First(&existing)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("document not found")
		}
		return fmt.Errorf("error fetching document: %w", result.Error)
	}

	// Hard delete document
	if err := r.db.WithContext(ctx).Delete(&existing).Error; err != nil {
		return fmt.Errorf("error removing document: %w", err)
	}

	return nil
}

// AddNote adds a note/observation to a student
// Parameters:
// - ctx: context for database operations
// - note: note data to add
// Returns:
// - error: any error encountered during the operation
func (r *studentRepository) AddNote(ctx context.Context, note *models.StudentNote) error {
	// Check required fields
	if note.Content == "" {
		return fmt.Errorf("content field is required")
	}

	// Check if student exists
	student, err := r.FindByID(ctx, note.StudentID)
	if err != nil {
		return err
	}
	if student == nil {
		return fmt.Errorf("student not found")
	}

	// Set author_id if not provided
	if note.AuthorID == 0 {
		// In a real scenario, you would get this from the request context
		userID, ok := ctx.Value("user_id").(uint)
		if ok {
			note.AuthorID = userID
		}
	}

	// Insert note
	if err := r.db.WithContext(ctx).Create(note).Error; err != nil {
		return fmt.Errorf("error adding note: %w", err)
	}

	return nil
}

// GetNotes returns a student's notes/observations
// Parameters:
// - ctx: context for database operations
// - studentID: ID of the student whose notes to fetch
// - includeConfidential: whether to include confidential notes
// Returns:
// - []models.StudentNote: list of note records
// - error: any error encountered during the operation
func (r *studentRepository) GetNotes(ctx context.Context, studentID uint, includeConfidential bool) ([]models.StudentNote, error) {
	var notes []models.StudentNote

	query := r.db.WithContext(ctx).Where("student_id = ?", studentID)

	// Filter confidential notes if necessary
	if !includeConfidential {
		query = query.Where("is_confidential = false")
	}

	// Order by creation date (most recent first)
	query = query.Order("created_at DESC")

	if err := query.Find(&notes).Error; err != nil {
		return nil, fmt.Errorf("error fetching notes: %w", err)
	}

	return notes, nil
}
