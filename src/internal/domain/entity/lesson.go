package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Lesson represents a lesson thread that groups multiple versions.
// Each lesson belongs to a module; versions carry the module_id for query ease.
type Lesson struct {
	ID        uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ModuleID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"module_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Lesson) TableName() string {
	return "lessons"
}

func (l *Lesson) Validate() error {
	if l.ModuleID == uuid.Nil {
		return fmt.Errorf("%w: module_id is required", ErrFieldRequired)
	}
	return nil
}

func NewLesson(moduleID uuid.UUID) (*Lesson, error) {
	lesson := &Lesson{
		ID:       uuid.New(),
		ModuleID: moduleID,
	}
	if err := lesson.Validate(); err != nil {
		return nil, err
	}
	return lesson, nil
}

// CanBeModifiedBy aligns lesson-level checks with course ownership rules.
func (l *Lesson) CanBeModifiedBy(course *Course, userID uuid.UUID) bool {
	return course.IsOwnedBy(userID)
}
