package entity

import (
	"fmt"
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
	ID          uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Action      AuditAction    `gorm:"type:varchar(50);not null" json:"action"`
	EntityID    uuid.UUID      `gorm:"type:uuid;not null" json:"entity_id"`
	EntityType  string         `gorm:"type:varchar(50);not null" json:"entity_type"`
	OldValues   string         `gorm:"type:jsonb" json:"old_values,omitempty"`
	NewValues   string         `gorm:"type:jsonb" json:"new_values,omitempty"`
	UserID      *uuid.UUID     `gorm:"type:uuid" json:"user_id,omitempty"`
	UserRole    *string        `gorm:"type:varchar(50)" json:"user_role,omitempty"`
	PerformedAt time.Time      `gorm:"autoCreateTime" json:"performed_at"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func (a *AuditLog) Validate() error {
	if a.Action == "" {
		return fmt.Errorf("%w: action is required", ErrFieldRequired)
	}

	if a.EntityID == uuid.Nil {
		return fmt.Errorf("%w: entity_id is required", ErrFieldRequired)
	}

	if a.EntityType == "" {
		return fmt.Errorf("%w: entity_type is required", ErrFieldRequired)
	}

	if a.UserID != nil {
		if _, err := uuid.Parse((*a.UserID).String()); err != nil {
			return fmt.Errorf("%w: user_id must be a valid UUID", ErrInvalidInput)
		}
	}

	return nil
}

func NewAuditLog(action AuditAction, entityID uuid.UUID, entityType, oldValues, newValues string, userID *uuid.UUID) (*AuditLog, error) {
	auditLog := &AuditLog{
		ID:         uuid.New(),
		Action:     action,
		EntityID:   entityID,
		EntityType: entityType,
		OldValues:  oldValues,
		NewValues:  newValues,
		UserID:     userID,
	}

	if err := auditLog.Validate(); err != nil {
		return nil, err
	}

	return auditLog, nil
}
