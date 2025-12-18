package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestModule_Validate(t *testing.T) {
	validCourseID := uuid.New().String()

	tests := []struct {
		name    string
		module  *Module
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid module",
			module: &Module{
				Title:      "Introduction",
				CourseID:   validCourseID,
				OrderIndex: 0,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			module: &Module{
				Title:      "",
				CourseID:   validCourseID,
				OrderIndex: 0,
			},
			wantErr: true,
			errMsg:  "title is required",
		},
		{
			name: "title too long",
			module: &Module{
				Title:      string(make([]byte, 256)),
				CourseID:   validCourseID,
				OrderIndex: 0,
			},
			wantErr: true,
			errMsg:  "must not exceed 255 characters",
		},
		{
			name: "empty course ID",
			module: &Module{
				Title:      "Introduction",
				CourseID:   "",
				OrderIndex: 0,
			},
			wantErr: true,
			errMsg:  "course_id is required",
		},
		{
			name: "invalid course UUID",
			module: &Module{
				Title:      "Introduction",
				CourseID:   "not-a-uuid",
				OrderIndex: 0,
			},
			wantErr: true,
			errMsg:  "must be a valid UUID",
		},
		{
			name: "negative order index",
			module: &Module{
				Title:      "Introduction",
				CourseID:   validCourseID,
				OrderIndex: -1,
			},
			wantErr: true,
			errMsg:  "must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.module.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestModule_BelongsToCourse(t *testing.T) {
	courseID := uuid.New().String()
	module := &Module{
		CourseID: courseID,
		Title:    "Test Module",
	}

	assert.True(t, module.BelongsToCourse(courseID))
	assert.False(t, module.BelongsToCourse(uuid.New().String()))
}

func TestModule_CanBeModifiedBy(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "John", "Instructor", UserRoleInstructor)
	otherInstructor, _ := NewUser("other@example.com", "password123", "Other", "Instructor", UserRoleInstructor)

	course := &Course{
		ID:           uuid.New().String(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}

	module := &Module{
		CourseID: course.ID,
		Title:    "Test Module",
	}

	assert.True(t, module.CanBeModifiedBy(course, instructor.ID))
	assert.False(t, module.CanBeModifiedBy(course, otherInstructor.ID))
}
