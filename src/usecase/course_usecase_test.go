package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCourseUseCase_CreateCourse(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)

	tests := []struct {
		name      string
		title     string
		desc      string
		instID    uuid.UUID
		level     entity.DifficultyLevel
		requestorID uuid.UUID
		requestorRole *entity.UserRole
		setupMock func(*MockCourseRepository, *MockUserRepository, *MockAuditLogRepository)
		wantErr   bool
		checkErr  func(error) bool
	}{
		{
			name:   "successful course creation",
			title:  "Test Course",
			desc:   "Test Description",
			instID: instructor.ID,
			level:  entity.DifficultyLevelBeginner,
			requestorID: instructor.ID,
			requestorRole: func() *entity.UserRole { role := entity.UserRoleInstructor; return &role }(),
			setupMock: func(courseRepo *MockCourseRepository, userRepo *MockUserRepository, auditRepo *MockAuditLogRepository) {
				userRepo.On("GetByID", mock.Anything, instructor.ID).Return(instructor, nil)
				courseRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.Course")).Return(nil)
				auditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "instructor not found",
			title:  "Test Course",
			desc:   "Test Description",
			instID: func() uuid.UUID { id := uuid.New(); return id }(),
			level:  entity.DifficultyLevelBeginner,
			requestorID: func() uuid.UUID { 
				var id uuid.UUID
				return id // This will be overwritten in the test to match instID
			}(),
			requestorRole: func() *entity.UserRole { role := entity.UserRoleInstructor; return &role }(),
			setupMock: func(courseRepo *MockCourseRepository, userRepo *MockUserRepository, auditRepo *MockAuditLogRepository) {
				// Mock will be set up with the actual instID in the test
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, entity.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCourseRepo := new(MockCourseRepository)
			mockUserRepo := new(MockUserRepository)
			mockAuditRepo := new(MockAuditLogRepository)

			requestorID := tt.requestorID
			if tt.name == "instructor not found" {
				requestorID = tt.instID // Make requestorID match instID so authorization passes
				mockUserRepo.On("GetByID", mock.Anything, tt.instID).Return(nil, entity.ErrNotFound)
			} else {
				tt.setupMock(mockCourseRepo, mockUserRepo, mockAuditRepo)
			}

			uc := NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
			course, err := uc.CreateCourse(context.Background(), tt.title, tt.desc, tt.instID, tt.level, requestorID, tt.requestorRole)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.checkErr != nil {
					assert.True(t, tt.checkErr(err))
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, course)
				assert.Equal(t, tt.title, course.Title)
			}

			mockCourseRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			mockAuditRepo.AssertExpectations(t)
		})
	}
}

func TestCourseUseCase_GetCourseByID(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	course := &entity.Course{
		ID:    uuid.New(),
		Title: "Test Course",
	}

	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)

	uc := NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	result, err := uc.GetCourseByID(context.Background(), course.ID)

	require.NoError(t, err)
	assert.Equal(t, course, result)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseUseCase_ListCourses(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	courses := []*entity.Course{
		{ID: uuid.New(), Title: "Course 1"},
		{ID: uuid.New(), Title: "Course 2"},
	}
	filter := &repository.CourseFilter{Limit: 10, Offset: 0}

	mockCourseRepo.On("List", mock.Anything, filter).Return(courses, int64(len(courses)), nil)

	uc := NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	result, total, err := uc.ListCourses(context.Background(), filter)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(len(courses)), total)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseUseCase_GetCoursesByInstructor(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	instructorID := uuid.New()
	courses := []*entity.Course{
		{ID: uuid.New(), Title: "Course 1", InstructorID: instructorID},
		{ID: uuid.New(), Title: "Course 2", InstructorID: instructorID},
	}

	mockCourseRepo.On("GetByInstructor", mock.Anything, instructorID).Return(courses, nil)

	uc := NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	result, err := uc.GetCoursesByInstructor(context.Background(), instructorID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	mockCourseRepo.AssertExpectations(t)
}

func TestCourseUseCase_UpdateCourse(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:              uuid.New(),
		Title:           "Updated Course",
		Description:     "Test description",
		InstructorID:    instructor.ID,
		DifficultyLevel: entity.DifficultyLevelBeginner,
	}

	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockUserRepo.On("GetByID", mock.Anything, instructor.ID).Return(instructor, nil)
	mockCourseRepo.On("Update", mock.Anything, course).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	uc := NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	err := uc.UpdateCourse(context.Background(), course, instructor.ID)

	require.NoError(t, err)
	mockCourseRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestCourseUseCase_DeleteCourse(t *testing.T) {
	instructor, _ := entity.NewUser("instructor@example.com", "password123", "John", "Instructor", entity.UserRoleInstructor)
	course := &entity.Course{
		ID:           uuid.New(),
		Title:        "Test Course",
		InstructorID: instructor.ID,
	}

	mockCourseRepo := new(MockCourseRepository)
	mockUserRepo := new(MockUserRepository)
	mockAuditRepo := new(MockAuditLogRepository)

	mockCourseRepo.On("GetByID", mock.Anything, course.ID).Return(course, nil)
	mockUserRepo.On("GetByID", mock.Anything, instructor.ID).Return(instructor, nil)
	mockCourseRepo.On("Delete", mock.Anything, course.ID).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.AuditLog")).Return(nil)

	uc := NewCourseUseCase(mockCourseRepo, mockUserRepo, mockAuditRepo)
	err := uc.DeleteCourse(context.Background(), course.ID, instructor.ID)

	require.NoError(t, err)
	mockCourseRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}
