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

func TestLessonUseCase_CreateLesson(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:       uuid.New(),
		CourseID: course.ID,
		Title:    "Module 1",
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockLessonRepo.On("CreateLesson", mock.Anything, mock.AnythingOfType("*entity.Lesson")).Return(nil)
	mockLessonRepo.On("CreateVersion", mock.Anything, mock.AnythingOfType("*entity.LessonVersion")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)
	lesson := &entity.LessonVersion{
		ModuleID: module.ID,
		Content:  "Lesson 1 content",
	}

	err := uc.CreateLesson(context.Background(), lesson, instructor.ID)

	require.NoError(t, err)
	assert.Equal(t, 1, lesson.VersionNumber)
	assert.NotEqual(t, uuid.Nil, lesson.LessonID)
	mockLessonRepo.AssertExpectations(t)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}

func TestLessonUseCase_CreateLessonVersion(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:       uuid.New(),
		CourseID: course.ID,
		Title:    "Module 1",
	}
	lessonID := uuid.New()
	baseLesson := &entity.Lesson{ID: lessonID, ModuleID: module.ID}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)

	mockLessonRepo.On("GetLessonByID", mock.Anything, lessonID).Return(baseLesson, nil)
	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockLessonRepo.On("GetNextVersionNumber", mock.Anything, lessonID).Return(2, nil)
	mockLessonRepo.On("CreateVersion", mock.Anything, mock.AnythingOfType("*entity.LessonVersion")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	newLesson := &entity.LessonVersion{
		Content: "New Content",
	}

	err := uc.CreateLessonVersion(context.Background(), lessonID, newLesson, instructor.ID)
	require.NoError(t, err)

	mockLessonRepo.AssertExpectations(t)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestLessonUseCase_GetLatestLessonsByModule(t *testing.T) {
	moduleID := uuid.New()
	lessons := []*entity.LessonVersion{
		{ID: uuid.New(), ModuleID: moduleID, Content: "Lesson 1", VersionNumber: 1},
		{ID: uuid.New(), ModuleID: moduleID, Content: "Lesson 2", VersionNumber: 1},
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)

	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(&entity.Module{ID: moduleID}, nil)
	mockLessonRepo.On("GetLatestByModule", mock.Anything, moduleID, 10, 0).Return(lessons, int64(len(lessons)), nil)

	result, total, err := uc.GetLatestLessonsByModule(context.Background(), moduleID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(len(lessons)), total)

	mockLessonRepo.AssertExpectations(t)
}

func TestLessonUseCase_GetAllLessonVersions(t *testing.T) {
	lessonID := uuid.New()
	versions := []*entity.LessonVersion{
		{ID: lessonID, VersionNumber: 1},
		{ID: lessonID, VersionNumber: 2},
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)

	mockLessonRepo.On("GetAllVersions", mock.Anything, lessonID, 10, 0).Return(versions, int64(len(versions)), nil)

	result, total, err := uc.GetAllLessonVersions(context.Background(), lessonID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(len(versions)), total)

	mockLessonRepo.AssertExpectations(t)
}

func TestLessonUseCase_GetLesson(t *testing.T) {
	lessonID := uuid.New()
	lesson := &entity.LessonVersion{ID: lessonID, VersionNumber: 1}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)

	mockLessonRepo.On("GetVersionByID", mock.Anything, lessonID).Return(lesson, nil)

	result, err := uc.GetLesson(context.Background(), lessonID)
	require.NoError(t, err)
	assert.Equal(t, lesson, result)

	mockLessonRepo.AssertExpectations(t)
}

func TestLessonUseCase_DeleteLesson(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:       uuid.New(),
		CourseID: course.ID,
		Title:    "Module 1",
	}
	lesson := &entity.Lesson{
		ID:       uuid.New(),
		ModuleID: module.ID,
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)

	mockLessonRepo.On("GetLessonByID", mock.Anything, lesson.ID).Return(lesson, nil)
	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockLessonRepo.On("DeleteVersionsByLesson", mock.Anything, lesson.ID).Return(nil)
	mockLessonRepo.On("DeleteLesson", mock.Anything, lesson.ID).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	err := uc.DeleteLesson(context.Background(), lesson.ID, instructor.ID)
	require.NoError(t, err)

	mockLessonRepo.AssertExpectations(t)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}
