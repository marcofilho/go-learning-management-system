package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestModule_Validate(t *testing.T) {
	validCourseID := uuid.New()

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
				CourseID:   uuid.Nil,
				OrderIndex: 0,
			},
			wantErr: true,
			errMsg:  "course_id is required",
		},
		// Note: Can't test invalid UUID at this level since uuid.UUID type is always valid
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
	courseID := uuid.New()
	module := &Module{
		CourseID: courseID,
		Title:    "Test Module",
	}

	assert.True(t, module.BelongsToCourse(courseID))
	assert.False(t, module.BelongsToCourse(uuid.New()))
}

func TestModule_CanBeModifiedBy(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "John", "Instructor", UserRoleInstructor)
	otherInstructor, _ := NewUser("other@example.com", "password123", "Other", "Instructor", UserRoleInstructor)

	course := &Course{
		ID:           uuid.New(),
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

func TestModule_TableName(t *testing.T) {
	var module Module
	assert.Equal(t, "modules", module.TableName())
}
