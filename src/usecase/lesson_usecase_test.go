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

type MockLessonRepository struct {
	mock.Mock
}

func (m *MockLessonRepository) Create(ctx context.Context, lesson *entity.LessonVersion) error {
	args := m.Called(ctx, lesson)
	return args.Error(0)
}

func (m *MockLessonRepository) GetByID(ctx context.Context, id string) (*entity.LessonVersion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) GetByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error) {
	args := m.Called(ctx, moduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) GetLatestByModule(ctx context.Context, moduleID string) ([]*entity.LessonVersion, error) {
	args := m.Called(ctx, moduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) GetNextVersionNumber(ctx context.Context, moduleID string) (int, error) {
	args := m.Called(ctx, moduleID)
	return args.Int(0), args.Error(1)
}

func (m *MockLessonRepository) GetAllVersions(ctx context.Context, lessonID string) ([]*entity.LessonVersion, error) {
	args := m.Called(ctx, lessonID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestLessonUseCase_CreateLesson(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New().String(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:       uuid.New().String(),
		CourseID: course.ID,
		Title:    "Module 1",
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockLessonRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.LessonVersion")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	lesson := &entity.LessonVersion{
		ModuleID: module.ID,
		Content:  "Lesson 1 content",
	}

	err := uc.CreateLesson(context.Background(), lesson, instructor.ID)

	require.NoError(t, err)
	assert.Equal(t, 1, lesson.VersionNumber)
	mockLessonRepo.AssertExpectations(t)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestLessonUseCase_CreateLesson_Unauthorized(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	otherInstructor, _ := entity.NewUser("other@example.com", "password123", "Jane", "Smith", entity.UserRoleInstructor)

	course := &entity.Course{
		ID:           uuid.New().String(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:       uuid.New().String(),
		CourseID: course.ID,
		Title:    "Module 1",
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	lesson := &entity.LessonVersion{
		ModuleID: module.ID,
		Content:  "Lesson 1 content",
	}

	err := uc.CreateLesson(context.Background(), lesson, otherInstructor.ID)

	assert.ErrorIs(t, err, entity.ErrUnauthorized)
}

func TestLessonUseCase_GetLessonsByModule(t *testing.T) {
	moduleID := uuid.New().String()
	module := &entity.Module{
		ID:    moduleID,
		Title: "Test Module",
	}
	lessons := []*entity.LessonVersion{
		{ID: "1", ModuleID: moduleID, Content: "Lesson 1", VersionNumber: 1},
		{ID: "2", ModuleID: moduleID, Content: "Lesson 2", VersionNumber: 1},
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)
	mockLessonRepo.On("GetLatestByModule", mock.Anything, moduleID).Return(lessons, nil)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	result, err := uc.GetLatestLessonsByModule(context.Background(), moduleID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockModuleRepo.AssertExpectations(t)
	mockLessonRepo.AssertExpectations(t)
}
