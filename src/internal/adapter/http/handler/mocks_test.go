package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/stretchr/testify/mock"
)

// User

type MockUserRepository struct{ mock.Mock }

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

func (m *MockUserRepository) GetByRole(ctx context.Context, role entity.UserRole) ([]*entity.User, error) {
	args := m.Called(ctx, role)
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

// Token

type MockTokenProvider struct{ mock.Mock }

func (m *MockTokenProvider) GenerateToken(user *entity.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockTokenProvider) ValidateToken(token string) (*auth.Claims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Claims), args.Error(1)
}

// Course

type MockCourseRepository struct{ mock.Mock }

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

func (m *MockCourseRepository) List(ctx context.Context, filter *repository.CourseFilter) ([]*entity.Course, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.Course), args.Get(1).(int64), args.Error(2)
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

// AuditLog

type MockAuditLogRepository struct{ mock.Mock }

func (m *MockAuditLogRepository) Create(ctx context.Context, auditLog *entity.AuditLog) error {
	args := m.Called(ctx, auditLog)
	return args.Error(0)
}

func (m *MockAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) List(ctx context.Context, limit, offset int) ([]*entity.AuditLog, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.AuditLog), args.Get(1).(int64), args.Error(2)
}

// Enrollment

type MockEnrollmentRepository struct{ mock.Mock }

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

func (m *MockEnrollmentRepository) GetByStudent(ctx context.Context, studentID uuid.UUID, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, int64, error) {
	args := m.Called(ctx, studentID, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.Enrollment), args.Get(1).(int64), args.Error(2)
}

func (m *MockEnrollmentRepository) GetByCourse(ctx context.Context, courseID uuid.UUID, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, int64, error) {
	args := m.Called(ctx, courseID, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.Enrollment), args.Get(1).(int64), args.Error(2)
}

func (m *MockEnrollmentRepository) Update(ctx context.Context, enrollment *entity.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) Delete(ctx context.Context, studentID, courseID uuid.UUID) error {
	args := m.Called(ctx, studentID, courseID)
	return args.Error(0)
}

// Module

type MockModuleRepository struct{ mock.Mock }

func (m *MockModuleRepository) Create(ctx context.Context, module *entity.Module) error {
	args := m.Called(ctx, module)
	return args.Error(0)
}

func (m *MockModuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Module, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Module), args.Error(1)
}

func (m *MockModuleRepository) GetByCourse(ctx context.Context, courseID uuid.UUID) ([]*entity.Module, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Module), args.Error(1)
}

func (m *MockModuleRepository) Update(ctx context.Context, module *entity.Module) error {
	args := m.Called(ctx, module)
	return args.Error(0)
}

func (m *MockModuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Lesson

type MockLessonRepository struct{ mock.Mock }

func (m *MockLessonRepository) CreateLesson(ctx context.Context, lesson *entity.Lesson) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *MockLessonRepository) GetLessonByID(ctx context.Context, id uuid.UUID) (*entity.Lesson, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Lesson), args.Error(1)
}

func (m *MockLessonRepository) DeleteLesson(ctx context.Context, lessonID uuid.UUID) error {
	args := m.Called(ctx, lessonID)
	return args.Error(0)
}

func (m *MockLessonRepository) CreateVersion(ctx context.Context, lesson *entity.LessonVersion) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *MockLessonRepository) GetVersionByID(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) GetLatestByModule(ctx context.Context, moduleID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	args := m.Called(ctx, moduleID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Get(1).(int64), args.Error(2)
}

func (m *MockLessonRepository) GetAllVersions(ctx context.Context, lessonID uuid.UUID, limit, offset int) ([]*entity.LessonVersion, int64, error) {
	args := m.Called(ctx, lessonID, limit, offset)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Get(1).(int64), args.Error(2)
}

func (m *MockLessonRepository) GetNextVersionNumber(ctx context.Context, lessonID uuid.UUID) (int, error) {
	args := m.Called(ctx, lessonID)
	return args.Int(0), args.Error(1)
}

func (m *MockLessonRepository) DeleteVersionsByLesson(ctx context.Context, lessonID uuid.UUID) error {
	args := m.Called(ctx, lessonID)
	return args.Error(0)
}
