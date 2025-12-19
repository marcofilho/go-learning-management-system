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

type MockModuleRepository struct {
	mock.Mock
}

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

func TestModuleUseCase_CreateModule(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}

	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockModuleRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Module")).Return(nil)

	uc := NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	module := &entity.Module{
		CourseID:   course.ID,
		Title:      "Module 1",
		OrderIndex: 0,
	}
	err := uc.CreateModule(context.Background(), module, instructor.ID)

	require.NoError(t, err)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}

func TestModuleUseCase_GetModulesByCourse(t *testing.T) {
	courseID := uuid.New()
	modules := []*entity.Module{
		{ID: uuid.New(), CourseID: courseID, Title: "Module 1"},
		{ID: uuid.New(), CourseID: courseID, Title: "Module 2"},
	}

	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(&entity.Course{ID: courseID}, nil)
	mockModuleRepo.On("GetByCourse", mock.Anything, courseID).Return(modules, nil)

	uc := NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	result, err := uc.GetModulesByCourse(context.Background(), courseID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockModuleRepo.AssertExpectations(t)
}

func TestModuleUseCase_GetModule(t *testing.T) {
	moduleID := uuid.New()
	module := &entity.Module{ID: moduleID, Title: "Module 1"}

	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil)

	uc := NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	result, err := uc.GetModule(context.Background(), moduleID)

	require.NoError(t, err)
	assert.Equal(t, module, result)
	mockModuleRepo.AssertExpectations(t)
}

func TestModuleUseCase_UpdateModule(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:       uuid.New(),
		CourseID: course.ID,
		Title:    "Old Title",
	}

	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockModuleRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.Module")).Return(nil)

	uc := NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	updatedModule := &entity.Module{
		ID:       module.ID,
		CourseID: course.ID,
		Title:    "New Title",
	}
	err := uc.UpdateModule(context.Background(), updatedModule, instructor.ID)

	require.NoError(t, err)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}

func TestModuleUseCase_DeleteModule(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:       uuid.New(),
		CourseID: course.ID,
	}

	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)
	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockModuleRepo.On("Delete", mock.Anything, module.ID).Return(nil)

	uc := NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	err := uc.DeleteModule(context.Background(), module.ID, instructor.ID)

	require.NoError(t, err)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}
