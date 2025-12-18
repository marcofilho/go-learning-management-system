package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnrollment(t *testing.T) {
	student, _ := NewUser("student@example.com", "pass123", "John", "Doe", UserRoleStudent)
	courseID := uuid.New().String()

	enrollment, err := NewEnrollment(student, courseID)
	require.NoError(t, err)
	require.NotNil(t, enrollment)
	assert.Equal(t, student.ID, enrollment.StudentID)
	assert.Equal(t, courseID, enrollment.CourseID)
	assert.Equal(t, EnrollmentStatusActive, enrollment.Status)
}

func TestEnrollment_Complete(t *testing.T) {
	enrollment := &Enrollment{
		StudentID: uuid.New().String(),
		CourseID:  uuid.New().String(),
		Status:    EnrollmentStatusActive,
	}
	enrollment.Complete()
	assert.Equal(t, EnrollmentStatusCompleted, enrollment.Status)
}

func TestEnrollment_Drop(t *testing.T) {
	enrollment := &Enrollment{
		StudentID: uuid.New().String(),
		CourseID:  uuid.New().String(),
		Status:    EnrollmentStatusActive,
	}
	enrollment.Drop()
	assert.Equal(t, EnrollmentStatusDropped, enrollment.Status)
}

func TestEnrollment_IsActive(t *testing.T) {
	enrollment := &Enrollment{Status: EnrollmentStatusActive}
	assert.True(t, enrollment.IsActive())

	enrollment.Status = EnrollmentStatusDropped
	assert.False(t, enrollment.IsActive())
}
