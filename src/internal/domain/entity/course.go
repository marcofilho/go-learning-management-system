package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DifficultyLevel string

const (
	DifficultyLevelBeginner     DifficultyLevel = "beginner"
	DifficultyLevelIntermediate DifficultyLevel = "intermediate"
	DifficultyLevelAdvanced     DifficultyLevel = "advanced"
)

type Course struct {
	ID              uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title           string          `gorm:"not null" json:"title"`
	Description     string          `gorm:"type:text" json:"description"`
	InstructorID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"instructor_id"`
	DifficultyLevel DifficultyLevel `gorm:"type:varchar(20);not null" json:"difficulty_level"`
	CreatedAt       time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`

	Instructor *User    `gorm:"foreignKey:InstructorID;constraint:OnDelete:CASCADE" json:"-"`
	Modules    []Module `gorm:"foreignKey:CourseID" json:"modules,omitempty"`
}

func (Course) TableName() string {
	return "courses"
}

func (c *Course) Validate() error {
	if _, err := uuid.Parse(c.ID.String()); err != nil && c.ID != uuid.Nil {
		return fmt.Errorf("%w: id must be a valid UUID", ErrInvalidInput)
	}

	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrFieldRequired)
	}

	if len(c.Title) > 255 {
		return fmt.Errorf("%w: title must not exceed 255 characters", ErrInvalidInput)
	}

	if c.InstructorID == uuid.Nil {
		return fmt.Errorf("%w: instructor_id is required", ErrFieldRequired)
	}

	if c.DifficultyLevel != DifficultyLevelBeginner &&
		c.DifficultyLevel != DifficultyLevelIntermediate &&
		c.DifficultyLevel != DifficultyLevelAdvanced {
		return fmt.Errorf("%w: difficulty_level must be one of: beginner, intermediate, advanced", ErrInvalidInput)
	}

	return nil
}

func NewCourse(title, description string, instructorID uuid.UUID, difficultyLevel DifficultyLevel, instructor *User) (*Course, error) {
	course := &Course{
		ID:              uuid.New(),
		Title:           title,
		Description:     description,
		InstructorID:    instructorID,
		DifficultyLevel: difficultyLevel,
		Instructor:      instructor,
	}

	if err := course.Validate(); err != nil {
		return nil, err
	}

	return course, nil
}

func (c *Course) IsOwnedBy(userID uuid.UUID) bool {
	return c.InstructorID == userID
}

func (c *Course) CanBeModifiedBy(userID uuid.UUID, userRole UserRole) bool {
	return userRole == UserRoleAdmin || c.InstructorID == userID
}
