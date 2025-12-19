package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Module struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CourseID   string         `gorm:"type:uuid;not null;index" json:"course_id"`
	Title      string         `gorm:"not null" json:"title"`
	OrderIndex int            `gorm:"not null;default:0" json:"order_index"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Course  Course          `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE" json:"-"`
	Lessons []LessonVersion `gorm:"foreignKey:ModuleID" json:"lessons,omitempty"`
}

func (Module) TableName() string {
	return "modules"
}

func (m *Module) Validate() error {
	if strings.TrimSpace(m.Title) == "" {
		return fmt.Errorf("%w: title is required", ErrFieldRequired)
	}

	if len(m.Title) > 255 {
		return fmt.Errorf("%w: title must not exceed 255 characters", ErrInvalidInput)
	}

	if strings.TrimSpace(m.CourseID) == "" {
		return fmt.Errorf("%w: course_id is required", ErrFieldRequired)
	}

	if _, err := uuid.Parse(m.CourseID); err != nil {
		return fmt.Errorf("%w: course_id must be a valid UUID", ErrInvalidInput)
	}

	if m.OrderIndex < 0 {
		return fmt.Errorf("%w: order_index must be non-negative", ErrInvalidInput)
	}

	return nil
}

func (m *Module) BelongsToCourse(courseID string) bool {
	return m.CourseID == courseID
}

func (m *Module) CanBeModifiedBy(course *Course, userID string) bool {
	return course.IsOwnedBy(userID)
}

func NewModule(title, courseID string) (*Module, error) {
	module := &Module{
		Title:    title,
		CourseID: courseID,
	}
	if err := module.Validate(); err != nil {
		return nil, err
	}
	return module, nil
}
