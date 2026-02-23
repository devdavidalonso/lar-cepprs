package models

import (
	"time"
)

// Location represents a physical location (classroom, lab, etc.)
type Location struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CenterID  *uint     `json:"centerId" gorm:"index"`
	ProgramID *uint     `json:"programId" gorm:"index"`
	Name      string    `json:"name" gorm:"not null"`
	Capacity  int       `json:"capacity"`
	Resources string    `json:"resources"` // Text field
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Center  *EducationalCenter `json:"center,omitempty" gorm:"foreignKey:CenterID"`
	Program *Program           `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}

// TableName defines the table name in the database
func (Location) TableName() string {
	return "locations"
}
