package models

import (
	"time"
)

// Course represents a course in the system
type Course struct {
	ID                  uint       `json:"id" gorm:"primaryKey"`
	ProgramID           *uint      `json:"programId" gorm:"index"`
	Name                string     `json:"name" gorm:"not null"`
	ShortDescription    string     `json:"shortDescription"`
	CoverImage          string     `json:"coverImage"` // URL to course banner
	DetailedDescription string     `json:"detailedDescription" gorm:"type:text"`
	Workload            int        `json:"workload" gorm:"not null"`
	GoogleClassroomURL  string     `json:"googleClassroomUrl"` // Link para a turma no Google Classroom
	GoogleClassroomID   string     `json:"googleClassroomId"`  // ID da turma na API do Google
	MaxStudents         int        `json:"maxStudents" gorm:"not null"`
	Prerequisites       string     `json:"prerequisites"`
	DifficultyLevel     string     `json:"difficultyLevel"`
	TargetAudience      string     `json:"targetAudience"`
	Tags                string     `json:"tags" gorm:"type:json"`
	WeekDays            string     `json:"weekDays" gorm:"not null"` // E.g., "1,3,5" for Monday, Wednesday, Friday
	StartTime           string     `json:"startTime" gorm:"not null"`
	EndTime             string     `json:"endTime" gorm:"not null"`
	Schedule            string     `json:"schedule"` // Human-readable schedule (e.g., "Saturdays 09:00-11:00")
	Duration            int        `json:"duration" gorm:"not null"` // In weeks
	StartDate           time.Time  `json:"startDate"`
	EndDate             time.Time  `json:"endDate"`
	Status              string     `json:"status" gorm:"not null;default:'active'"`
	CreatedAt           time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt           time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt           *time.Time `json:"deletedAt" gorm:"index"`

	// Associations
	ClassSessions  []ClassSession  `json:"classSessions,omitempty" gorm:"foreignKey:CourseID"`
	TeacherCourses []TeacherCourse `json:"teacherCourses,omitempty" gorm:"foreignKey:CourseID"`
	Program        *Program        `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}

// TableName defines the table name in the database
func (Course) TableName() string {
	return "courses"
}

// TeacherCourse represents the association between a teacher and a course
type TeacherCourse struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	TeacherID uint       `json:"teacherId" gorm:"not null;index"` // References teachers.id
	CourseID  uint       `json:"courseId" gorm:"not null;index"`
	Role      string     `json:"role" gorm:"not null"` // primary, assistant, substitute
	StartDate time.Time  `json:"startDate" gorm:"not null"`
	EndDate   *time.Time `json:"endDate"`
	Active    bool       `json:"active" gorm:"default:true"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`

	// Associations
	Teacher Teacher `json:"teacher,omitempty" gorm:"foreignKey:TeacherID"`
	Course  Course  `json:"course,omitempty" gorm:"foreignKey:CourseID"`
}

// TableName defines the table name in the database
func (TeacherCourse) TableName() string {
	return "teacher_courses"
}

// Enrollment represents a student's enrollment in a course
type Enrollment struct {
	ID                     uint       `json:"id" gorm:"primaryKey"`
	StudentID              uint       `json:"studentId" gorm:"not null;index"`
	CourseID               uint       `json:"courseId" gorm:"not null;index"`
	EnrollmentNumber       string     `json:"enrollmentNumber" gorm:"not null;unique"`
	Status                 string     `json:"status" gorm:"not null;default:'active'"` // active, in_progress, locked, completed, cancelled
	StartDate              time.Time  `json:"startDate" gorm:"not null"`
	EndDate                *time.Time `json:"endDate"`
	EnrollmentDate         time.Time  `json:"enrollmentDate" gorm:"not null"`
	CancellationReason     string     `json:"cancellationReason"`
	AgreementURL           string     `json:"agreementUrl"`
	GoogleInvitationStatus string     `json:"googleInvitationStatus" gorm:"default:'not_sent'"` // not_sent, pending, accepted
	CreatedAt              time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt              time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt              *time.Time `json:"deletedAt" gorm:"index"`
}

// TableName defines the table name in the database
func (Enrollment) TableName() string {
	return "enrollments"
}

// EnrollmentHistory records changes to an enrollment's status
type EnrollmentHistory struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	EnrollmentID   uint      `json:"enrollmentId" gorm:"not null;index"`
	PreviousStatus string    `json:"previousStatus"`
	NewStatus      string    `json:"newStatus" gorm:"not null"`
	Notes          string    `json:"notes"`
	ChangeDate     time.Time `json:"changeDate" gorm:"not null"`
	UserID         uint      `json:"userId" gorm:"not null"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

// TableName defines the table name in the database
func (EnrollmentHistory) TableName() string {
	return "enrollment_history"
}

// WaitingList represents a student on a course's waiting list
type WaitingList struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	StudentID    uint       `json:"studentId" gorm:"not null;index"`
	CourseID     uint       `json:"courseId" gorm:"not null;index"`
	RegisterDate time.Time  `json:"registerDate" gorm:"not null"`
	Priority     int        `json:"priority" gorm:"default:0"`
	Notes        string     `json:"notes"`
	Status       string     `json:"status" gorm:"not null;default:'waiting'"` // waiting, called, withdrawn
	CallDate     *time.Time `json:"callDate"`
	ResponseDate *time.Time `json:"responseDate"`
	Response     string     `json:"response"` // accepted, declined
	CreatedAt    time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName defines the table name in the database
func (WaitingList) TableName() string {
	return "waiting_list"
}

// Certificate represents a certificate issued for an enrollment
type Certificate struct {
	ID               uint       `json:"id" gorm:"primaryKey"`
	EnrollmentID     uint       `json:"enrollmentId" gorm:"not null;index"`
	StudentID        uint       `json:"studentId" gorm:"not null;index"`
	CourseID         uint       `json:"courseId" gorm:"not null;index"`
	Type             string     `json:"type" gorm:"not null"` // completion, participation, in_progress
	IssueDate        time.Time  `json:"issueDate" gorm:"not null"`
	ExpiryDate       *time.Time `json:"expiryDate"`
	CertificateURL   string     `json:"certificateUrl"`
	VerificationCode string     `json:"verificationCode" gorm:"unique"`
	QRCodeURL        string     `json:"qrCodeUrl"`
	Status           string     `json:"status" gorm:"not null;default:'active'"` // active, revoked, expired
	RevocationReason string     `json:"revocationReason"`
	CreatedByID      uint       `json:"createdById" gorm:"not null"`
	CreatedAt        time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName defines the table name in the database
func (Certificate) TableName() string {
	return "certificates"
}

// CertificateTemplate represents a certificate template
type CertificateTemplate struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"not null"` // completion, participation, in_progress
	HTMLContent string    `json:"htmlContent" gorm:"type:text;not null"`
	CSSStyles   string    `json:"cssStyles" gorm:"type:text"`
	IsDefault   bool      `json:"isDefault" gorm:"default:false"`
	IsActive    bool      `json:"isActive" gorm:"default:true"`
	CreatedByID uint      `json:"createdById" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName defines the table name in the database
func (CertificateTemplate) TableName() string {
	return "certificate_templates"
}

// Registration represents a student enrollment (legacy version for compatibility)
type Registration struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	StudentID uint       `json:"studentId" gorm:"not null;index"`
	CourseID  uint       `json:"courseId" gorm:"not null;index"`
	Status    string     `json:"status" gorm:"not null"` // active, completed, canceled
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
	Student   Student    `json:"student" gorm:"foreignKey:StudentID"`
	Course    Course     `json:"course" gorm:"foreignKey:CourseID"`
}
