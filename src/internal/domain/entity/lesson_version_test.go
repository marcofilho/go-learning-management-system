package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLessonVersion_Validate(t *testing.T) {
	validModuleID := uuid.New()
	validLessonID := uuid.New()

	tests := []struct {
		name    string
		lesson  *LessonVersion
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid lesson with content",
			lesson: &LessonVersion{
				LessonID:      validLessonID,
				ModuleID:      validModuleID,
				VersionNumber: 1,
				Content:       "Lesson content here",
			},
			wantErr: false,
		},
		{
			name: "valid lesson with video URL",
			lesson: &LessonVersion{
				LessonID:      validLessonID,
				ModuleID:      validModuleID,
				VersionNumber: 1,
				VideoURL:      "https://example.com/video.mp4",
			},
			wantErr: false,
		},
		{
			name: "valid lesson with both content and video",
			lesson: &LessonVersion{
				LessonID:      validLessonID,
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
				LessonID:      validLessonID,
				ModuleID:      uuid.Nil,
				VersionNumber: 1,
				Content:       "Content",
			},
			wantErr: true,
			errMsg:  "module_id is required",
		},
		// Note: Can't test invalid UUID at this level since uuid.UUID type is always valid
		{
			name: "version number less than 1",
			lesson: &LessonVersion{
				LessonID:      validLessonID,
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
				LessonID:      validLessonID,
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
	moduleID := uuid.New()
	lesson := &LessonVersion{
		LessonID:      uuid.New(),
		ModuleID:      moduleID,
		VersionNumber: 1,
		Content:       "Test content",
	}

	assert.True(t, lesson.BelongsToModule(moduleID))
	assert.False(t, lesson.BelongsToModule(uuid.New()))
}

func TestLessonVersion_CanBeModifiedBy(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "John", "Instructor", UserRoleInstructor)
	otherInstructor, _ := NewUser("other@example.com", "password123", "Other", "Instructor", UserRoleInstructor)

	course := &Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}

	lesson := &LessonVersion{
		LessonID:      uuid.New(),
		ModuleID:      uuid.New(),
		VersionNumber: 1,
		Content:       "Test content",
	}

	assert.True(t, lesson.CanBeModifiedBy(course, instructor.ID))
	assert.False(t, lesson.CanBeModifiedBy(course, otherInstructor.ID))
}

func TestLessonVersion_TableName(t *testing.T) {
	var lessonVersion LessonVersion
	assert.Equal(t, "lesson_versions", lessonVersion.TableName())
}
