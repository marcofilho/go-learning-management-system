package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditAction string

const (
	AuditActionCertificationWebhook AuditAction = "certification_webhook"
	AuditActionCourseCreated        AuditAction = "course_created"
	AuditActionCourseUpdated        AuditAction = "course_updated"
	AuditActionCourseDeleted        AuditAction = "course_deleted"
	AuditActionEnrollmentCreated    AuditAction = "enrollment_created"
	AuditActionEnrollmentUpdated    AuditAction = "enrollment_updated"
	AuditActionEnrollmentDeleted    AuditAction = "enrollment_deleted"
	AuditActionLessonVersionCreated AuditAction = "lesson_version_created"
	AuditActionLessonVersionUpdated AuditAction = "lesson_version_updated"
	AuditActionLessonDeleted        AuditAction = "lesson_deleted"
)

type AuditLog struct {
	ID            string      `gorm:"type:uuid;primaryKey" json:"id"`
	Action        AuditAction `gorm:"type:varchar(100);not null" json:"action"`
	ResourceID    string      `gorm:"type:uuid;not null" json:"resource_id"`
	ResourceType  string      `gorm:"type:varchar(50);not null" json:"resource_type"`
	PayloadBefore string      `gorm:"type:jsonb" json:"payload_before,omitempty"`
	PayloadAfter  string      `gorm:"type:jsonb" json:"payload_after,omitempty"`
	UserID        *string     `gorm:"type:uuid" json:"user_id,omitempty"`
	CreatedAt     time.Time   `gorm:"not null" json:"created_at"`

	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *AuditLog) Validate() error {
	if strings.TrimSpace(string(a.Action)) == "" {
		return fmt.Errorf("%w: action is required", ErrFieldRequired)
	}

	if strings.TrimSpace(a.ResourceID) == "" {
		return fmt.Errorf("%w: resource_id is required", ErrFieldRequired)
	}

	if _, err := uuid.Parse(a.ResourceID); err != nil {
		return fmt.Errorf("%w: resource_id must be a valid UUID", ErrInvalidInput)
	}

	if strings.TrimSpace(a.ResourceType) == "" {
		return fmt.Errorf("%w: resource_type is required", ErrFieldRequired)
	}

	if a.UserID != nil {
		if _, err := uuid.Parse(*a.UserID); err != nil {
			return fmt.Errorf("%w: user_id must be a valid UUID", ErrInvalidInput)
		}
	}

	return nil
}

func NewAuditLog(action AuditAction, resourceID, resourceType, payloadBefore, payloadAfter string, userID *string) (*AuditLog, error) {
	auditLog := &AuditLog{
		ID:            uuid.New().String(),
		Action:        action,
		ResourceID:    resourceID,
		ResourceType:  resourceType,
		PayloadBefore: payloadBefore,
		PayloadAfter:  payloadAfter,
		UserID:        userID,
		CreatedAt:     time.Now(),
	}

	if err := auditLog.Validate(); err != nil {
		return nil, err
	}

	return auditLog, nil
}
