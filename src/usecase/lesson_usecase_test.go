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

func (m *MockLessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.LessonVersion, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) GetByModule(ctx context.Context, moduleID uuid.UUID) ([]*entity.LessonVersion, error) {
	args := m.Called(ctx, moduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) GetLatestByModule(ctx context.Context, moduleID uuid.UUID) ([]*entity.LessonVersion, error) {
	args := m.Called(ctx, moduleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) GetNextVersionNumber(ctx context.Context, moduleID uuid.UUID) (int, error) {
	args := m.Called(ctx, moduleID)
	return args.Int(0), args.Error(1)
}

func (m *MockLessonRepository) GetAllVersions(ctx context.Context, lessonID uuid.UUID) ([]*entity.LessonVersion, error) {
	args := m.Called(ctx, lessonID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.LessonVersion), args.Error(1)
}

func (m *MockLessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
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
	lesson := &entity.LessonVersion{
		ID:            uuid.New(),
		ModuleID:      module.ID,
		VersionNumber: 1,
		Content:       "Old Content",
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)

	mockLessonRepo.On("GetByID", mock.Anything, lesson.ID).Return(lesson, nil)
	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockLessonRepo.On("GetNextVersionNumber", mock.Anything, lesson.ModuleID).Return(2, nil)
	mockLessonRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.LessonVersion")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	newLesson := &entity.LessonVersion{
		Content: "New Content",
	}

	err := uc.CreateLessonVersion(context.Background(), lesson.ID, newLesson, instructor.ID)
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
	mockLessonRepo.On("GetLatestByModule", mock.Anything, moduleID).Return(lessons, nil)

	result, err := uc.GetLatestLessonsByModule(context.Background(), moduleID)
	require.NoError(t, err)
	assert.Len(t, result, 2)

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

	mockLessonRepo.On("GetAllVersions", mock.Anything, lessonID).Return(versions, nil)

	result, err := uc.GetAllLessonVersions(context.Background(), lessonID)
	require.NoError(t, err)
	assert.Len(t, result, 2)

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

	mockLessonRepo.On("GetByID", mock.Anything, lessonID).Return(lesson, nil)

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
	lesson := &entity.LessonVersion{
		ID:       uuid.New(),
		ModuleID: module.ID,
	}

	mockLessonRepo := new(MockLessonRepository)
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	uc := NewLessonUseCase(mockLessonRepo, mockModuleRepo, mockCourseRepo, mockAuditRepo)

	mockLessonRepo.On("GetByID", mock.Anything, lesson.ID).Return(lesson, nil)
	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockLessonRepo.On("Delete", mock.Anything, lesson.ID).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	err := uc.DeleteLesson(context.Background(), lesson.ID, instructor.ID)
	require.NoError(t, err)

	mockLessonRepo.AssertExpectations(t)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}
