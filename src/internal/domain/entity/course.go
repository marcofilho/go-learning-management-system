package entity

import (
        "github.com/google/uuid"
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

func NewCourse(title, description, instructorID string, difficultyLevel DifficultyLevel, instructor *User) (*Course, error) {
    if title == "" {
		return nil, ErrInvalidInput
	}
	
	if !instructor.IsInstructor() {
                return nil, ErrInsufficientPermissions
        }

        if difficultyLevel == "" {
                difficultyLevel = DifficultyLevelBeginner
        }

        return &Course{
                ID:              uuid.New().String(),
                Title:           title,
                Description:     description,
                InstructorID:    instructorID,
                DifficultyLevel: difficultyLevel,
                CreatedAt:       time.Now(),
                UpdatedAt:       time.Now(),
        }, nil
}

func (c *Course) CanBeModifiedBy(user *User) bool {
        return user.IsAdmin() || c.InstructorID == user.ID
}

func (c *Course) IsOwnedBy(instructorID string) bool {
        return c.InstructorID == instructorID
}
