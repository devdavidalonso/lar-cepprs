package models

import (
	"time"

	"gorm.io/gorm"
)

// StudentStatus represents the possible statuses of a student
type StudentStatus string

const (
	StudentStatusActive    StudentStatus = "active"
	StudentStatusInactive  StudentStatus = "inactive"
	StudentStatusSuspended StudentStatus = "suspended"
)

// Student represents a student in the system
// It specializes the User entity by having a reference to a User
type Student struct {
	ID                 uint          `json:"id" gorm:"primaryKey"`
	UserID             uint          `json:"userId" gorm:"not null;unique;index"`                         // Reference to User entity
	RegistrationNumber string        `json:"registrationNumber" gorm:"size:10;uniqueIndex;not null"`      // Auto-generated
	Status             StudentStatus `json:"status" gorm:"type:student_status;not null;default:'active'"` // ENUM: active, inactive, suspended
	SpecialNeeds       string        `json:"specialNeeds"`
	MedicalInfo        string        `json:"medicalInfo"`
	SocialMedia        *string       `json:"socialMedia" gorm:"type:json"`
	Notes              string        `json:"notes"`
	CreatedAt          time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          time.Time     `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt          *time.Time    `json:"deletedAt" gorm:"index"`
	ProgramIDs         []uint        `json:"programIds,omitempty" gorm:"-"`

	// Associations
	User         User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Guardians    []Guardian    `json:"guardians,omitempty" gorm:"foreignKey:StudentID"`
	UserContacts []UserContact `json:"userContacts,omitempty" gorm:"foreignKey:StudentID"`
	Documents    []Document    `json:"documents,omitempty" gorm:"foreignKey:StudentID"`
	StudentNotes []StudentNote `json:"studentNotes,omitempty" gorm:"foreignKey:StudentID"`
	Enrollments  []Enrollment  `json:"enrollments,omitempty" gorm:"foreignKey:StudentID"`
	Programs     []StudentProgram `json:"programs,omitempty" gorm:"foreignKey:StudentID"`
}

// TableName defines the table name in the database
func (Student) TableName() string {
	return "students"
}

// BeforeCreate hook - clears registration_number to force auto-generation by database
func (s *Student) BeforeCreate(tx *gorm.DB) error {
	// Clear registration_number to force auto-generation by database trigger
	s.RegistrationNumber = ""
	return nil
}
