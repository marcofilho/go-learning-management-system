package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCertificationWebhookUseCase_ProcessCertification_Passed(t *testing.T) {
	student, _ := entity.NewUser("student@example.com", "password123", "John", "Student", entity.UserRoleStudent)
	course := &entity.Course{
		ID:    uuid.New(),
		Title: "Test Course",
	}
	enrollment := &entity.Enrollment{
		StudentID: student.ID,
		CourseID:  course.ID,
		Status:    entity.EnrollmentStatusActive,
	}

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockUserRepo.On("GetByID", mock.Anything, student.ID).Return(student, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockEnrollmentRepo.On("GetByStudentAndCourse", mock.Anything, student.ID, course.ID).Return(enrollment, nil)
	mockEnrollmentRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Enrollment")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	uc := NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditRepo)
	payload := CertificationWebhookPayload{
		StudentID:           student.ID,
		CourseID:            course.ID,
		CertificationStatus: "passed",
		Score:               95,
		Timestamp:           time.Now(),
	}

	err := uc.ProcessCertification(context.Background(), payload)

	require.NoError(t, err)
	assert.Equal(t, entity.EnrollmentStatusCompleted, enrollment.Status)
	mockUserRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestCertificationWebhookUseCase_ProcessCertification_Failed(t *testing.T) {
	student, _ := entity.NewUser("student@example.com", "password123", "John", "Student", entity.UserRoleStudent)
	course := &entity.Course{
		ID:    uuid.New(),
		Title: "Test Course",
	}
	enrollment := &entity.Enrollment{
		StudentID: student.ID,
		CourseID:  course.ID,
		Status:    entity.EnrollmentStatusActive,
	}

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockUserRepo.On("GetByID", mock.Anything, student.ID).Return(student, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockEnrollmentRepo.On("GetByStudentAndCourse", mock.Anything, student.ID, course.ID).Return(enrollment, nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	uc := NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditRepo)
	payload := CertificationWebhookPayload{
		StudentID:           student.ID,
		CourseID:            course.ID,
		CertificationStatus: "failed",
		Score:               45,
		Timestamp:           time.Now(),
	}

	err := uc.ProcessCertification(context.Background(), payload)

	require.NoError(t, err)
	assert.Equal(t, entity.EnrollmentStatusActive, enrollment.Status)
	mockUserRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestCertificationWebhookUseCase_ProcessCertification_StudentNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	invalidID := uuid.New()
	mockUserRepo.On("GetByID", mock.Anything, invalidID).Return(nil, entity.ErrNotFound)

	uc := NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditRepo)
	payload := CertificationWebhookPayload{
		StudentID:           invalidID,
		CourseID:            uuid.New(),
		CertificationStatus: "passed",
		Score:               95,
		Timestamp:           time.Now(),
	}

	err := uc.ProcessCertification(context.Background(), payload)

	assert.ErrorIs(t, err, entity.ErrNotFound)
	mockUserRepo.AssertExpectations(t)
}

func TestCertificationWebhookUseCase_ProcessCertification_InactiveEnrollment(t *testing.T) {
	student, _ := entity.NewUser("student@example.com", "password123", "John", "Student", entity.UserRoleStudent)
	course := &entity.Course{
		ID:    uuid.New(),
		Title: "Test Course",
	}
	enrollment := &entity.Enrollment{
		StudentID: student.ID,
		CourseID:  course.ID,
		Status:    entity.EnrollmentStatusDropped,
	}

	mockUserRepo := new(MockUserRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockUserRepo.On("GetByID", mock.Anything, student.ID).Return(student, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockEnrollmentRepo.On("GetByStudentAndCourse", mock.Anything, student.ID, course.ID).Return(enrollment, nil)

	uc := NewCertificationWebhookUseCase(mockUserRepo, mockCourseRepo, mockEnrollmentRepo, mockAuditRepo)
	payload := CertificationWebhookPayload{
		StudentID:           student.ID,
		CourseID:            course.ID,
		CertificationStatus: "passed",
		Score:               95,
		Timestamp:           time.Now(),
	}

	err := uc.ProcessCertification(context.Background(), payload)

	assert.ErrorIs(t, err, entity.ErrInvalidInput)
	mockUserRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
}
