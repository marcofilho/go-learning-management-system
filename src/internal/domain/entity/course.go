package entity

import (
	"time"

	"gorm.io/gorm"
)

type DifficultyLevel string

const (
	DifficultyLevelBeginner     DifficultyLevel = "beginner"
	DifficultyLevelIntermediate DifficultyLevel = "intermediate"
	DifficultyLevelAdvanced     DifficultyLevel = "advanced"
)

type Course struct {
	ID              string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title           string          `gorm:"not null" json:"title"`
	Description     string          `gorm:"type:text" json:"description"`
	InstructorID    string          `gorm:"type:uuid;not null;index" json:"instructor_id"`
	DifficultyLevel DifficultyLevel `gorm:"type:varchar(20);not null;default:'beginner'" json:"difficulty_level"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`

	Instructor User     `gorm:"foreignKey:InstructorID;constraint:OnDelete:CASCADE" json:"-"`
	Modules    []Module `gorm:"foreignKey:CourseID" json:"modules,omitempty"`
}

func (Course) TableName() string {
	return "courses"
}
