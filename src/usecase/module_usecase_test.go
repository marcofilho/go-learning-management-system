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

func (m *MockModuleRepository) GetByID(ctx context.Context, id string) (*entity.Module, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Module), args.Error(1)
}

func (m *MockModuleRepository) GetByCourse(ctx context.Context, courseID string) ([]*entity.Module, error) {
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

func (m *MockModuleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestModuleUseCase_CreateModule(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New().String(),
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
	courseID := uuid.New().String()
	modules := []*entity.Module{
		{ID: "1", CourseID: courseID, Title: "Module 1"},
		{ID: "2", CourseID: courseID, Title: "Module 2"},
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
