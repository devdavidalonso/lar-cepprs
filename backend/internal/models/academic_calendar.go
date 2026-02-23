package models

import "time"

// AcademicCalendar represents a recess, holiday, or special event in the program calendar.
type AcademicCalendar struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProgramID *uint     `json:"programId" gorm:"index"` // Optional: if null, applies to all
	Name      string    `json:"name" gorm:"not null"`
	StartDate time.Time `json:"startDate" gorm:"not null"`
	EndDate   time.Time `json:"endDate" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null;default:'recess'"` // recess, holiday, event
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Program *Program `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}

func (AcademicCalendar) TableName() string {
	return "academic_calendars"
}
