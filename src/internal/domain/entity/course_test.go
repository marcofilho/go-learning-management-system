package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCourse(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "Inst", "Ructor", UserRoleInstructor)
	course, err := NewCourse("Test Course", "Description", instructor.ID, DifficultyLevelBeginner, instructor)
	require.NoError(t, err)
	require.NotNil(t, course)
	assert.Equal(t, "Test Course", course.Title)
}

func TestCourse_CanBeModifiedBy(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "I", "T", UserRoleInstructor)
	admin, _ := NewUser("admin@example.com", "password123", "A", "D", UserRoleAdmin)
	course := &Course{ID: uuid.New(), InstructorID: instructor.ID}

	assert.True(t, course.CanBeModifiedBy(instructor.ID, instructor.Role))
	assert.True(t, course.CanBeModifiedBy(admin.ID, admin.Role))
}

func TestNewCourse_Validation(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "I", "T", UserRoleInstructor)
	student, _ := NewUser("student@example.com", "password123", "S", "T", UserRoleStudent)

	_, err := NewCourse("", "Description", instructor.ID, DifficultyLevelBeginner, instructor)
	assert.Error(t, err, "should fail with empty title")

	_, err = NewCourse("Title", "Description", instructor.ID, DifficultyLevelBeginner, student)
	assert.Error(t, err, "should fail if creator is not an instructor or admin")
}

func TestCourse_TableName(t *testing.T) {
	var course Course
	assert.Equal(t, "courses", course.TableName())
}

func TestCourse_Validate(t *testing.T) {
	instructor, _ := NewUser("instructor@example.com", "password123", "I", "T", UserRoleInstructor)
	tests := []struct {
		name    string
		course  *Course
		wantErr bool
	}{
		{
			name: "EmptyTitle",
			course: &Course{
				Title:        "",
				Description:  "Description",
				InstructorID: instructor.ID,
			},
			wantErr: true,
		},
		{
			name: "EmptyDescription",
			course: &Course{
				Title:        "Title",
				Description:  "",
				InstructorID: instructor.ID,
			},
			wantErr: true,
		},
		{
			name: "EmptyInstructorID",
			course: &Course{
				Title:        "Title",
				Description:  "Description",
				InstructorID: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "InvalidInstructorID",
			course: &Course{
				Title:        "Title",
				Description:  "Description",
				InstructorID: uuid.Nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.course.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
