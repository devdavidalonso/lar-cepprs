package models

import "time"

// EducationalCenter represents an educational center (e.g. CE Prof. Paulo Rossi Severino).
type EducationalCenter struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex"`
	Code      string    `json:"code" gorm:"not null;uniqueIndex"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (EducationalCenter) TableName() string {
	return "educational_centers"
}

// Program represents a program under an educational center (Semear, Voar, Cecor).
type Program struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CenterID  uint      `json:"centerId" gorm:"not null;index"`
	Code      string    `json:"code" gorm:"not null;uniqueIndex"`
	Name      string    `json:"name" gorm:"not null"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Center EducationalCenter `json:"center,omitempty" gorm:"foreignKey:CenterID"`
	// Optional direct links to simplify management screens and permissions
	TeacherPrograms []TeacherProgram `json:"teacherPrograms,omitempty" gorm:"foreignKey:ProgramID"`
}

func (Program) TableName() string {
	return "programs"
}

// StudentProgram links one student to one program (N:N via student + program).
type StudentProgram struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	StudentID uint       `json:"studentId" gorm:"not null;index;uniqueIndex:idx_student_program_unique"`
	ProgramID uint       `json:"programId" gorm:"not null;index;uniqueIndex:idx_student_program_unique"`
	Status    string     `json:"status" gorm:"not null;default:'active'"`
	EntryDate *time.Time `json:"entryDate"`
	ExitDate  *time.Time `json:"exitDate"`
	CreatedAt time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`

	Student Student `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	Program Program `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}

func (StudentProgram) TableName() string {
	return "student_programs"
}

// TeacherProgram links one teacher to one program.
// This relation is important when a teacher can act in multiple programs.
type TeacherProgram struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TeacherID uint      `json:"teacherId" gorm:"not null;index;uniqueIndex:idx_teacher_program_unique"`
	ProgramID uint      `json:"programId" gorm:"not null;index;uniqueIndex:idx_teacher_program_unique"`
	Role      string    `json:"role" gorm:"default:'teacher'"` // teacher, coordinator
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Teacher Teacher `json:"teacher,omitempty" gorm:"foreignKey:TeacherID"`
	Program Program `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}

func (TeacherProgram) TableName() string {
	return "teacher_programs"
}
