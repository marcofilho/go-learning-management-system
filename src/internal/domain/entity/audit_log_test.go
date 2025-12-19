package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuditLog(t *testing.T) {
	entityID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		action        AuditAction
		entityID      uuid.UUID
		entityType    string
		payloadBefore string
		payloadAfter  string
		userID        *uuid.UUID
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid audit log",
			action:        AuditActionCourseCreated,
			entityID:      entityID,
			entityType:    "course",
			payloadBefore: "",
			payloadAfter:  `{"title":"Course"}`,
			userID:        &userID,
			wantErr:       false,
		},
		{
			name:          "valid without user ID",
			action:        AuditActionCertificationWebhook,
			entityID:      entityID,
			entityType:    "certification",
			payloadBefore: "",
			payloadAfter:  `{"status":"completed"}`,
			userID:        nil,
			wantErr:       false,
		},
		{
			name:          "empty action",
			action:        AuditAction(""),
			entityID:      entityID,
			entityType:    "course",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        &userID,
			wantErr:       true,
			errMsg:        "action is required",
		},
		{
			name:          "empty resource ID",
			action:        AuditActionCourseCreated,
			entityID:      uuid.Nil,
			entityType:    "course",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        &userID,
			wantErr:       true,
			errMsg:        "entity_id is required",
		},
		// Note: Can't test invalid UUID at this level since uuid.UUID type is always valid
		// Invalid UUID would be caught at parsing level before reaching this constructor
		{
			name:          "empty resource type",
			action:        AuditActionCourseCreated,
			entityID:      entityID,
			entityType:    "",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        &userID,
			wantErr:       true,
			errMsg:        "entity_type is required",
		},
		// Note: nil userID is valid - validation only checks if non-nil UUID is valid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := NewAuditLog(tt.action, tt.entityID, tt.entityType, tt.payloadBefore, tt.payloadAfter, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, log)
			} else {
				require.NoError(t, err)
				require.NotNil(t, log)
				assert.NotEmpty(t, log.ID)
				assert.Equal(t, tt.action, log.Action)
				assert.Equal(t, tt.entityID, log.EntityID)
				assert.Equal(t, tt.entityType, log.EntityType)
			}
		})
	}
}

func TestAuditLog_Validate(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name    string
		log     *AuditLog
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid audit log",
			log: &AuditLog{
				Action:     AuditActionCourseCreated,
				EntityID:   validUUID,
				EntityType: "course",
			},
			wantErr: false,
		},
		{
			name: "empty action",
			log: &AuditLog{
				Action:     AuditAction(""),
				EntityID:   validUUID,
				EntityType: "course",
			},
			wantErr: true,
			errMsg:  "action is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.log.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
