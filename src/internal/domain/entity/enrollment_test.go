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

func TestEnrollment_TableName(t *testing.T) {
	var enrollment Enrollment
	assert.Equal(t, "course_enrollments", enrollment.TableName())
}

func TestEnrollment_Validate(t *testing.T) {
	student, _ := NewUser("student@example.com", "pass123", "John", "Doe", UserRoleStudent)
	tests := []struct {
		name       string
		enrollment *Enrollment
		wantErr    bool
	}{
		{
			name: "EmptyStudentID",
			enrollment: &Enrollment{
				StudentID: "",
				CourseID:  uuid.New().String(),
				Status:    EnrollmentStatusActive,
			},
			wantErr: true,
		},
		{
			name: "InvalidStudentID",
			enrollment: &Enrollment{
				StudentID: "invalid-uuid",
				CourseID:  uuid.New().String(),
				Status:    EnrollmentStatusActive,
			},
			wantErr: true,
		},
		{
			name: "EmptyCourseID",
			enrollment: &Enrollment{
				StudentID: student.ID,
				CourseID:  "",
				Status:    EnrollmentStatusActive,
			},
			wantErr: true,
		},
		{
			name: "InvalidCourseID",
			enrollment: &Enrollment{
				StudentID: student.ID,
				CourseID:  "invalid-uuid",
				Status:    EnrollmentStatusActive,
			},
			wantErr: true,
		},
		{
			name: "InvalidStatus",
			enrollment: &Enrollment{
				StudentID: student.ID,
				CourseID:  uuid.New().String(),
				Status:    "invalid-status",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.enrollment.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnrollment_CanBeModifiedBy(t *testing.T) {
	admin, _ := NewUser("admin@example.com", "pass123", "Admin", "User", UserRoleAdmin)
	student, _ := NewUser("student@example.com", "pass123", "John", "Doe", UserRoleStudent)
	otherStudent, _ := NewUser("other@example.com", "pass123", "Jane", "Doe", UserRoleStudent)

	enrollment := &Enrollment{
		StudentID: student.ID,
	}

	assert.True(t, enrollment.CanBeModifiedBy(admin))
	assert.True(t, enrollment.CanBeModifiedBy(student))
	assert.False(t, enrollment.CanBeModifiedBy(otherStudent))
}

func TestEnrollment_IsCompleted(t *testing.T) {
	enrollment := &Enrollment{Status: EnrollmentStatusCompleted}
	assert.True(t, enrollment.IsCompleted())

	enrollment.Status = EnrollmentStatusActive
	assert.False(t, enrollment.IsCompleted())
}
