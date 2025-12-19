package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Module struct {
	ID         uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CourseID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"course_id"`
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

	if m.CourseID == uuid.Nil {
		return fmt.Errorf("%w: course_id is required", ErrFieldRequired)
	}

	if m.OrderIndex < 0 {
		return fmt.Errorf("%w: order_index must be non-negative", ErrInvalidInput)
	}

	return nil
}

func NewModule(title string, courseID uuid.UUID) (*Module, error) {
	module := &Module{
		ID:       uuid.New(),
		Title:    title,
		CourseID: courseID,
	}
	if err := module.Validate(); err != nil {
		return nil, err
	}
	return module, nil
}

func (m *Module) BelongsToCourse(courseID uuid.UUID) bool {
	return m.CourseID == courseID
}

func (m *Module) CanBeModifiedBy(course *Course, userID uuid.UUID) bool {
	return course.IsOwnedBy(userID)
}
