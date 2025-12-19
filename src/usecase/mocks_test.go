package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByRole(ctx context.Context, role entity.UserRole) ([]*entity.User, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

// MockCourseRepository for testing
type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) Create(ctx context.Context, course *entity.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Course), args.Error(1)
}

func (m *MockCourseRepository) List(ctx context.Context, filter *repository.CourseFilter) ([]*entity.Course, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Course), args.Error(1)
}

func (m *MockCourseRepository) GetByInstructor(ctx context.Context, instructorID uuid.UUID) ([]*entity.Course, error) {
	args := m.Called(ctx, instructorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Course), args.Error(1)
}

func (m *MockCourseRepository) Update(ctx context.Context, course *entity.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEnrollmentRepository for testing
type MockEnrollmentRepository struct {
	mock.Mock
}

func (m *MockEnrollmentRepository) Create(ctx context.Context, enrollment *entity.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) GetByStudentAndCourse(ctx context.Context, studentID, courseID uuid.UUID) (*entity.Enrollment, error) {
	args := m.Called(ctx, studentID, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) Update(ctx context.Context, enrollment *entity.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) Delete(ctx context.Context, studentID, courseID uuid.UUID) error {
	args := m.Called(ctx, studentID, courseID)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) GetByStudent(ctx context.Context, studentID uuid.UUID, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	args := m.Called(ctx, studentID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) GetByCourse(ctx context.Context, courseID uuid.UUID, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	args := m.Called(ctx, courseID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Enrollment), args.Error(1)
}

// MockAuditLogRepository for testing
type MockAuditLogRepository struct {
	mock.Mock
}

func (m *MockAuditLogRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	args := m.Called(ctx, log)
	return args.Error(0)
}

func (m *MockAuditLogRepository) List(ctx context.Context, limit, offset int) ([]*entity.AuditLog, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuditLog), args.Error(1)
}
