package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEnrollmentUseCase_EnrollStudent(t *testing.T) {
	student, _ := entity.NewUser("student@example.com", "password123", "John", "Student", entity.UserRoleStudent)
	admin, _ := entity.NewUser("admin@example.com", "password123", "Admin", "User", entity.UserRoleAdmin)
	course := &entity.Course{
		ID:    uuid.New().String(),
		Title: "Test Course",
	}

	tests := []struct {
		name        string
		studentID   string
		courseID    string
		requestorID string
		setupMocks  func(*MockEnrollmentRepository, *MockCourseRepository, *MockUserRepository, *MockAuditLogRepository)
		wantErr     bool
	}{
		{
			name:        "successful enrollment by self",
			studentID:   student.ID,
			courseID:    course.ID,
			requestorID: student.ID,
			setupMocks: func(enrollRepo *MockEnrollmentRepository, courseRepo *MockCourseRepository, userRepo *MockUserRepository, auditRepo *MockAuditLogRepository) {
				courseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
				userRepo.On("GetByID", mock.Anything, student.ID).Return(student, nil)
				enrollRepo.On("GetByStudentAndCourse", mock.Anything, student.ID, course.ID).Return(nil, entity.ErrNotFound)
				enrollRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
				auditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "successful enrollment by admin",
			studentID:   student.ID,
			courseID:    course.ID,
			requestorID: admin.ID,
			setupMocks: func(enrollRepo *MockEnrollmentRepository, courseRepo *MockCourseRepository, userRepo *MockUserRepository, auditRepo *MockAuditLogRepository) {
				courseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
				userRepo.On("GetByID", mock.Anything, student.ID).Return(student, nil)
				userRepo.On("GetByID", mock.Anything, admin.ID).Return(admin, nil)
				enrollRepo.On("GetByStudentAndCourse", mock.Anything, student.ID, course.ID).Return(nil, entity.ErrNotFound)
				enrollRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
				auditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:        "course not found",
			studentID:   student.ID,
			courseID:    "invalid-course",
			requestorID: student.ID,
			setupMocks: func(enrollRepo *MockEnrollmentRepository, courseRepo *MockCourseRepository, userRepo *MockUserRepository, auditRepo *MockAuditLogRepository) {
				courseRepo.On("GetByID", mock.Anything, "invalid-course").Return(nil, entity.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEnrollRepo := new(MockEnrollmentRepository)
			mockCourseRepo := new(MockCourseRepository)
			mockUserRepo := new(MockUserRepository)
			mockAuditRepo := new(MockAuditLogRepository)

			tt.setupMocks(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)

			uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
			enrollment, err := uc.EnrollStudent(context.Background(), tt.studentID, tt.courseID, tt.requestorID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, enrollment)
				assert.Equal(t, tt.studentID, enrollment.StudentID)
				assert.Equal(t, tt.courseID, enrollment.CourseID)
			}

			mockEnrollRepo.AssertExpectations(t)
			mockCourseRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockAuditRepo.AssertExpectations(t)
		})
	}
}

func TestEnrollmentUseCase_GetStudentCourses(t *testing.T) {
	studentID := uuid.New().String()
	student, _ := entity.NewUser("student@example.com", "password123", "John", "Student", entity.UserRoleStudent)
	enrollments := []*entity.Enrollment{
		{StudentID: student.ID, CourseID: "course-1"},
		{StudentID: student.ID, CourseID: "course-2"},
	}

	mockEnrollRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockEnrollRepo.On("GetByStudent", mock.Anything, studentID, mock.AnythingOfType("*repository.EnrollmentFilter")).Return(enrollments, nil)

	uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	result, err := uc.GetStudentCourses(context.Background(), studentID, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentUseCase_GetCourseStudents(t *testing.T) {
	courseID := uuid.New().String()
	enrollments := []*entity.Enrollment{
		{StudentID: "student-1", CourseID: courseID},
		{StudentID: "student-2", CourseID: courseID},
	}

	mockEnrollRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockEnrollRepo.On("GetByCourse", mock.Anything, courseID, mock.AnythingOfType("*repository.EnrollmentFilter")).Return(enrollments, nil)

	uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	result, err := uc.GetCourseStudents(context.Background(), courseID, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockEnrollRepo.AssertExpectations(t)
}
