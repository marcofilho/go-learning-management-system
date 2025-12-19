package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonVersion struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ModuleID      string         `gorm:"type:uuid;not null;index" json:"module_id"`
	Title         string         `gorm:"not null" json:"title"`
	VersionNumber int            `gorm:"not null;default:1" json:"version_number"`
	Content       string         `gorm:"type:text" json:"content"`
	VideoURL      string         `gorm:"type:varchar(500)" json:"video_url,omitempty"`
	AttachmentURL string         `gorm:"type:varchar(500)" json:"attachment_url,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Module Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"-"`
}

func (LessonVersion) TableName() string {
	return "lesson_versions"
}

func (l *LessonVersion) Validate() error {
	if strings.TrimSpace(l.ModuleID) == "" {
		return fmt.Errorf("%w: module_id is required", ErrFieldRequired)
	}

	if _, err := uuid.Parse(l.ModuleID); err != nil {
		return fmt.Errorf("%w: module_id must be a valid UUID", ErrInvalidInput)
	}

	if l.VersionNumber < 1 {
		return fmt.Errorf("%w: version_number must be at least 1", ErrInvalidInput)
	}

	if strings.TrimSpace(l.Content) == "" && strings.TrimSpace(l.VideoURL) == "" {
		return fmt.Errorf("%w: either content or video_url must be provided", ErrInvalidInput)
	}

	return nil
}

func (l *LessonVersion) BelongsToModule(moduleID string) bool {
	return l.ModuleID == moduleID
}

func (l *LessonVersion) CanBeModifiedBy(course *Course, userID string) bool {
	return course.IsOwnedBy(userID)
}

func NewLessonVersion(title, content, moduleID string, version int) (*LessonVersion, error) {
	lesson := &LessonVersion{
		ModuleID:      moduleID,
		VersionNumber: version,
		Content:       content,
	}
	if err := lesson.Validate(); err != nil {
		return nil, err
	}
	return lesson, nil
}
