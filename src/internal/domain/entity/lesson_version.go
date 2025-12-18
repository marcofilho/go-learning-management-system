package entity

import (
        "github.com/google/uuid"
        "strings"
        "time"

        "gorm.io/gorm"
)

type LessonVersion struct {
        ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
        ModuleID      string         `gorm:"type:uuid;not null;index" json:"module_id"`
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
                return ErrFieldRequired
        }

        if _, err := uuid.Parse(l.ModuleID); err != nil {
                return ErrInvalidInput
        }

        if l.VersionNumber < 1 {
                return ErrInvalidInput
        }

        if strings.TrimSpace(l.Content) == "" && strings.TrimSpace(l.VideoURL) == "" {
                return ErrInvalidInput
        }

        return nil
}

func (l *LessonVersion) BelongsToModule(moduleID string) bool {
        return l.ModuleID == moduleID
}

func (l *LessonVersion) CanBeModifiedBy(course *Course, userID string) bool {
        return course.IsOwnedBy(userID)
}
