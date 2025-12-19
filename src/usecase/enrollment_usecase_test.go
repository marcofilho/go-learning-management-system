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
		ID:    uuid.New(),
		Title: "Test Course",
	}

	tests := []struct {
		name        string
		studentID   uuid.UUID
		courseID    uuid.UUID
		requestorID uuid.UUID
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
			courseID:    uuid.New(),
			requestorID: student.ID,
			setupMocks: func(enrollRepo *MockEnrollmentRepository, courseRepo *MockCourseRepository, userRepo *MockUserRepository, auditRepo *MockAuditLogRepository) {
				courseRepo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, entity.ErrNotFound)
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
	studentID := uuid.New()
	student, _ := entity.NewUser("student@example.com", "password123", "John", "Student", entity.UserRoleStudent)
	enrollments := []*entity.Enrollment{
		{StudentID: student.ID, CourseID: uuid.New()},
		{StudentID: student.ID, CourseID: uuid.New()},
	}

	mockEnrollRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockEnrollRepo.On("GetByStudent", mock.Anything, studentID, mock.AnythingOfType("*repository.EnrollmentFilter")).Return(enrollments, int64(len(enrollments)), nil)

	uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	result, total, err := uc.GetStudentCourses(context.Background(), studentID, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(len(enrollments)), total)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentUseCase_GetCourseStudents(t *testing.T) {
	courseID := uuid.New()
	enrollments := []*entity.Enrollment{
		{StudentID: uuid.New(), CourseID: courseID},
		{StudentID: uuid.New(), CourseID: courseID},
	}

	mockEnrollRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockEnrollRepo.On("GetByCourse", mock.Anything, courseID, mock.AnythingOfType("*repository.EnrollmentFilter")).Return(enrollments, int64(len(enrollments)), nil)

	uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)
	result, total, err := uc.GetCourseStudents(context.Background(), courseID, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(len(enrollments)), total)
	mockEnrollRepo.AssertExpectations(t)
}

func TestEnrollmentUseCase_UpdateEnrollmentStatus(t *testing.T) {
	admin, _ := entity.NewUser("admin@example.com", "password123", "Admin", "User", entity.UserRoleAdmin)
	student, _ := entity.NewUser("student@example.com", "password123", "Student", "User", entity.UserRoleStudent)
	enrollment := &entity.Enrollment{
		StudentID: student.ID,
		CourseID:  uuid.New(),
		Status:    entity.EnrollmentStatusActive,
	}

	mockEnrollRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)

	mockUserRepo.On("GetByID", mock.Anything, admin.ID).Return(admin, nil)
	enrollmentCopy := *enrollment
	enrollmentCopy.Status = entity.EnrollmentStatusCompleted
	mockEnrollRepo.On("GetByStudentAndCourse", mock.Anything, enrollment.StudentID, enrollment.CourseID).Return(enrollment, nil)
	mockEnrollRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	updatedEnrollment, err := uc.UpdateEnrollmentStatus(context.Background(), enrollment.StudentID, enrollment.CourseID, entity.EnrollmentStatusCompleted, admin.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedEnrollment)
	require.Equal(t, entity.EnrollmentStatusCompleted, updatedEnrollment.Status)

	mockEnrollRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestEnrollmentUseCase_DropEnrollment(t *testing.T) {
	admin, _ := entity.NewUser("admin@example.com", "password123", "Admin", "User", entity.UserRoleAdmin)
	student, _ := entity.NewUser("student@example.com", "password123", "Student", "User", entity.UserRoleStudent)
	enrollment := &entity.Enrollment{
		StudentID: student.ID,
		CourseID:  uuid.New(),
		Status:    entity.EnrollmentStatusActive,
	}

	mockEnrollRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)

	mockUserRepo.On("GetByID", mock.Anything, admin.ID).Return(admin, nil)
	mockEnrollRepo.On("GetByStudentAndCourse", mock.Anything, enrollment.StudentID, enrollment.CourseID).Return(enrollment, nil)
	mockEnrollRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	updatedEnrollment, err := uc.DropEnrollment(context.Background(), enrollment.StudentID, enrollment.CourseID, admin.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedEnrollment)
	require.Equal(t, entity.EnrollmentStatusDropped, updatedEnrollment.Status)

	mockEnrollRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestEnrollmentUseCase_CompleteEnrollment(t *testing.T) {
	admin, _ := entity.NewUser("admin@example.com", "password123", "Admin", "User", entity.UserRoleAdmin)
	student, _ := entity.NewUser("student@example.com", "password123", "Student", "User", entity.UserRoleStudent)
	enrollment := &entity.Enrollment{
		StudentID: student.ID,
		CourseID:  uuid.New(),
		Status:    entity.EnrollmentStatusActive,
	}

	mockEnrollRepo := new(MockEnrollmentRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewEnrollmentUseCase(mockEnrollRepo, mockCourseRepo, mockUserRepo, mockAuditRepo)

	mockUserRepo.On("GetByID", mock.Anything, admin.ID).Return(admin, nil)
	mockEnrollRepo.On("GetByStudentAndCourse", mock.Anything, enrollment.StudentID, enrollment.CourseID).Return(enrollment, nil)
	mockEnrollRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	updatedEnrollment, err := uc.CompleteEnrollment(context.Background(), enrollment.StudentID, enrollment.CourseID, admin.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedEnrollment)
	require.Equal(t, entity.EnrollmentStatusCompleted, updatedEnrollment.Status)

	mockEnrollRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}
