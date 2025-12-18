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
	course := &Course{ID: uuid.New().String(), InstructorID: instructor.ID}

	assert.True(t, course.CanBeModifiedBy(instructor))
	assert.True(t, course.CanBeModifiedBy(admin))
}
