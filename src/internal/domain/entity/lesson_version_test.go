package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLessonVersion_Validate(t *testing.T) {
	validModuleID := uuid.New().String()

	tests := []struct {
		name    string
		lesson  *LessonVersion
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid lesson with content",
			lesson: &LessonVersion{
				ModuleID:      validModuleID,
				VersionNumber: 1,
				Content:       "Lesson content here",
			},
			wantErr: false,
		},
		{
			name: "valid lesson with video URL",
			lesson: &LessonVersion{
				ModuleID:      validModuleID,
				VersionNumber: 1,
				VideoURL:      "https://example.com/video.mp4",
			},
			wantErr: false,
		},
		{
			name: "valid lesson with both content and video",
			lesson: &LessonVersion{
				ModuleID:      validModuleID,
				VersionNumber: 1,
				Content:       "Lesson content",
				VideoURL:      "https://example.com/video.mp4",
			},
			wantErr: false,
		},
		{
			name: "empty module ID",
			lesson: &LessonVersion{
				ModuleID:      "",
				VersionNumber: 1,
				Content:       "Content",
			},
			wantErr: true,
			errMsg:  "module_id is required",
		},
		{
			name: "invalid module UUID",
			lesson: &LessonVersion{
				ModuleID:      "not-a-uuid",
				VersionNumber: 1,
				Content:       "Content",
			},
			wantErr: true,
			errMsg:  "must be a valid UUID",
		},
		{
			name: "version number less than 1",
			lesson: &LessonVersion{
				ModuleID:      validModuleID,
				VersionNumber: 0,
				Content:       "Content",
			},
			wantErr: true,
			errMsg:  "must be at least 1",
		},
		{
			name: "no content or video URL",
			lesson: &LessonVersion{
				ModuleID:      validModuleID,
				VersionNumber: 1,
				Content:       "",
				VideoURL:      "",
			},
			wantErr: true,
			errMsg:  "either content or video_url must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.lesson.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLessonVersion_BelongsToModule(t *testing.T) {
	moduleID := uuid.New().String()
	lesson := &LessonVersion{
		ModuleID:      moduleID,
		VersionNumber: 1,
		Content:       "Test content",
	}

	assert.True(t, lesson.BelongsToModule(moduleID))
	assert.False(t, lesson.BelongsToModule(uuid.New().String()))
}

func TestLessonVersion_CanBeModifiedBy(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "John", "Instructor", UserRoleInstructor)
	otherInstructor, _ := NewUser("other@example.com", "password123", "Other", "Instructor", UserRoleInstructor)

	course := &Course{
		ID:           uuid.New().String(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}

	lesson := &LessonVersion{
		ModuleID:      uuid.New().String(),
		VersionNumber: 1,
		Content:       "Test content",
	}

	assert.True(t, lesson.CanBeModifiedBy(course, instructor.ID))
	assert.False(t, lesson.CanBeModifiedBy(course, otherInstructor.ID))
}
