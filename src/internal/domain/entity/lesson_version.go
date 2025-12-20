package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonVersion struct {
	ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	LessonID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"lesson_id"`
	ModuleID      uuid.UUID      `gorm:"-" json:"module_id"`
	VersionNumber int            `gorm:"not null;default:1" json:"version_number"`
	Content       string         `gorm:"type:text" json:"content"`
	VideoURL      string         `gorm:"type:varchar(500)" json:"video_url,omitempty"`
	AttachmentURL string         `gorm:"-" json:"attachment_url,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"-" json:"-"`

	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnDelete:CASCADE" json:"-"`
	Module Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"-"`
}

func (LessonVersion) TableName() string {
	return "lesson_versions"
}

func (l *LessonVersion) Validate() error {
	if l.LessonID == uuid.Nil {
		return fmt.Errorf("%w: lesson_id is required", ErrFieldRequired)
	}

	if l.ModuleID == uuid.Nil {
		return fmt.Errorf("%w: module_id is required", ErrFieldRequired)
	}

	if l.VersionNumber < 1 {
		return fmt.Errorf("%w: version_number must be at least 1", ErrInvalidInput)
	}

	if strings.TrimSpace(l.Content) == "" && strings.TrimSpace(l.VideoURL) == "" {
		return fmt.Errorf("%w: either content or video_url must be provided", ErrInvalidInput)
	}

	return nil
}

func NewLessonVersion(title, content string, moduleID uuid.UUID, version int) (*LessonVersion, error) {
	lesson := &LessonVersion{
		ID:            uuid.New(),
		ModuleID:      moduleID,
		VersionNumber: version,
		Content:       content,
	}
	if err := lesson.Validate(); err != nil {
		return nil, err
	}
	return lesson, nil
}

func (l *LessonVersion) BelongsToModule(moduleID uuid.UUID) bool {
	return l.ModuleID == moduleID
}

func (l *LessonVersion) CanBeModifiedBy(course *Course, userID uuid.UUID) bool {
	return course.IsOwnedBy(userID)
}
