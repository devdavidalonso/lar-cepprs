package models

import (
	"time"
)

// Teacher represents a teacher in the system
type Teacher struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	UserID         uint       `json:"userId" gorm:"not null;uniqueIndex"`
	Specialization string     `json:"specialization"`
	Bio            string     `json:"bio"`
	Phone          string     `json:"phone"`
	Active         bool       `json:"active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`

	// Associations
	User         User                  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Courses      []Course              `json:"courses,omitempty" gorm:"many2many:teacher_courses;"`
	Skills       []TeacherSkill        `json:"skills,omitempty" gorm:"foreignKey:TeacherID"`
	Availability []TeacherAvailability `json:"availability,omitempty" gorm:"foreignKey:TeacherID"`
	Programs     []TeacherProgram      `json:"programs,omitempty" gorm:"foreignKey:TeacherID"`
}

// TableName defines the table name in the database
func (Teacher) TableName() string {
	return "teachers"
}
