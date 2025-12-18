package entity

import (
        "github.com/google/uuid"
        "strings"
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

func (c *Course) Validate() error {
        if strings.TrimSpace(c.Title) == "" {
                return ErrFieldRequired
        }

        if len(c.Title) > 255 {
                return ErrInvalidInput
        }

        if strings.TrimSpace(c.Description) == "" {
                return ErrFieldRequired
        }

        if strings.TrimSpace(c.InstructorID) == "" {
                return ErrFieldRequired
        }

        if _, err := uuid.Parse(c.InstructorID); err != nil {
                return ErrInvalidInput
        }

        if c.DifficultyLevel != DifficultyLevelBeginner &&
                c.DifficultyLevel != DifficultyLevelIntermediate &&
                c.DifficultyLevel != DifficultyLevelAdvanced {
                return ErrInvalidInput
        }

        return nil
}

func NewCourse(title, description, instructorID string, difficultyLevel DifficultyLevel, instructor *User) (*Course, error) {
        if !instructor.IsInstructor() {
                return nil, ErrInsufficientPermissions
        }

        if difficultyLevel == "" {
                difficultyLevel = DifficultyLevelBeginner
        }

        course := &Course{
                ID:              uuid.New().String(),
                Title:           title,
                Description:     description,
                InstructorID:    instructorID,
                DifficultyLevel: difficultyLevel,
                CreatedAt:       time.Now(),
                UpdatedAt:       time.Now(),
        }

        if err := course.Validate(); err != nil {
                return nil, err
        }

        return course, nil
}

func (c *Course) CanBeModifiedBy(user *User) bool {
        return user.IsAdmin() || c.InstructorID == user.ID
}

func (c *Course) IsOwnedBy(instructorID string) bool {
        return c.InstructorID == instructorID
}
