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

// MockModuleRepository for testing
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
	otherInstructor, _ := entity.NewUser("other@example.com", "password123", "Other", "Instructor", entity.UserRoleInstructor)

	course := &entity.Course{
		ID:           uuid.New().String(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}

	tests := []struct {
		name        string
		module      *entity.Module
		instructorID string
		setupMock   func(*MockModuleRepository, *MockCourseRepository)
		wantErr     bool
	}{
		{
			name: "successful module creation",
			module: &entity.Module{
				CourseID:   course.ID,
				Title:      "Module 1",
				OrderIndex: 0,
			},
			instructorID: instructor.ID,
			setupMock: func(moduleRepo *MockModuleRepository, courseRepo *MockCourseRepository) {
				courseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
				moduleRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Module")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "unauthorized instructor",
			module: &entity.Module{
				CourseID:   course.ID,
				Title:      "Module 1",
				OrderIndex: 0,
			},
			instructorID: otherInstructor.ID,
			setupMock: func(moduleRepo *MockModuleRepository, courseRepo *MockCourseRepository) {
				courseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
			},
			wantErr: true,
		},
		{
			name: "course not found",
			module: &entity.Module{
				CourseID:   uuid.New().String(),
				Title:      "Module 1",
				OrderIndex: 0,
			},
			instructorID: instructor.ID,
			setupMock: func(moduleRepo *MockModuleRepository, courseRepo *MockCourseRepository) {
				courseRepo.On("GetByID", mock.Anything, mock.Anything).Return(nil, entity.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockModuleRepo := new(MockModuleRepository)
			mockCourseRepo := new(MockCourseRepository)

			tt.setupMock(mockModuleRepo, mockCourseRepo)

			uc := NewModuleUseCase(mockModuleRepo, mockCourseRepo)
			err := uc.CreateModule(context.Background(), tt.module, tt.instructorID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockModuleRepo.AssertExpectations(t)
			mockCourseRepo.AssertExpectations(t)
		})
	}
}

func TestModuleUseCase_GetModulesByCourse(t *testing.T) {
	courseID := uuid.New().String()
	modules := []*entity.Module{
		{ID: "1", CourseID: courseID, Title: "Module 1", OrderIndex: 0},
		{ID: "2", CourseID: courseID, Title: "Module 2", OrderIndex: 1},
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
	mockCourseRepo.AssertExpectations(t)
}

func TestModuleUseCase_GetModule(t *testing.T) {
	module := &entity.Module{
		ID:         uuid.New().String(),
		CourseID:   uuid.New().String(),
		Title:      "Test Module",
		OrderIndex: 0,
	}

	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)

	mockModuleRepo.On("GetByID", mock.Anything, module.ID).Return(module, nil)

	uc := NewModuleUseCase(mockModuleRepo, mockCourseRepo)
	result, err := uc.GetModule(context.Background(), module.ID)

	require.NoError(t, err)
	assert.Equal(t, module, result)
	mockModuleRepo.AssertExpectations(t)
}

func TestModuleUseCase_DeleteModule(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New().String(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}
	module := &entity.Module{
		ID:         uuid.New().String(),
		CourseID:   course.ID,
		Title:      "Test Module",
		OrderIndex: 0,
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
