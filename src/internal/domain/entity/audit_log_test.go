package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuditLog(t *testing.T) {
	resourceID := uuid.New().String()
	userID := uuid.New().String()

	tests := []struct {
		name          string
		action        AuditAction
		resourceID    string
		resourceType  string
		payloadBefore string
		payloadAfter  string
		userID        *string
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid audit log",
			action:        AuditActionCourseCreated,
			resourceID:    resourceID,
			resourceType:  "course",
			payloadBefore: "",
			payloadAfter:  `{"title":"Course"}`,
			userID:        &userID,
			wantErr:       false,
		},
		{
			name:          "valid without user ID",
			action:        AuditActionCertificationWebhook,
			resourceID:    resourceID,
			resourceType:  "certification",
			payloadBefore: "",
			payloadAfter:  `{"status":"completed"}`,
			userID:        nil,
			wantErr:       false,
		},
		{
			name:          "empty action",
			action:        AuditAction(""),
			resourceID:    resourceID,
			resourceType:  "course",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        &userID,
			wantErr:       true,
			errMsg:        "action is required",
		},
		{
			name:          "empty resource ID",
			action:        AuditActionCourseCreated,
			resourceID:    "",
			resourceType:  "course",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        &userID,
			wantErr:       true,
			errMsg:        "resource_id is required",
		},
		{
			name:          "invalid resource UUID",
			action:        AuditActionCourseCreated,
			resourceID:    "not-a-uuid",
			resourceType:  "course",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        &userID,
			wantErr:       true,
			errMsg:        "must be a valid UUID",
		},
		{
			name:          "empty resource type",
			action:        AuditActionCourseCreated,
			resourceID:    resourceID,
			resourceType:  "",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        &userID,
			wantErr:       true,
			errMsg:        "resource_type is required",
		},
		{
			name:          "invalid user UUID",
			action:        AuditActionCourseCreated,
			resourceID:    resourceID,
			resourceType:  "course",
			payloadBefore: "",
			payloadAfter:  `{}`,
			userID:        stringPtr("not-a-uuid"),
			wantErr:       true,
			errMsg:        "must be a valid UUID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, err := NewAuditLog(tt.action, tt.resourceID, tt.resourceType, tt.payloadBefore, tt.payloadAfter, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, log)
			} else {
				require.NoError(t, err)
				require.NotNil(t, log)
				assert.NotEmpty(t, log.ID)
				assert.Equal(t, tt.action, log.Action)
				assert.Equal(t, tt.resourceID, log.ResourceID)
				assert.Equal(t, tt.resourceType, log.ResourceType)
			}
		})
	}
}

func TestAuditLog_Validate(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name    string
		log     *AuditLog
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid audit log",
			log: &AuditLog{
				Action:       AuditActionCourseCreated,
				ResourceID:   validUUID,
				ResourceType: "course",
			},
			wantErr: false,
		},
		{
			name: "empty action",
			log: &AuditLog{
				Action:       AuditAction(""),
				ResourceID:   validUUID,
				ResourceType: "course",
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
