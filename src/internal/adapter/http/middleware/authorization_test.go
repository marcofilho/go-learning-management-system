package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) GetByID(ctx context.Context, id string) (*entity.Course, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Course), args.Error(1)
}

func (m *MockCourseRepository) Create(ctx context.Context, course *entity.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) Update(ctx context.Context, course *entity.Course) error {
	args := m.Called(ctx, course)
	return args.Error(0)
}

func (m *MockCourseRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCourseRepository) List(ctx context.Context, filter *repository.CourseFilter) ([]*entity.Course, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Course), args.Error(1)
}

func (m *MockCourseRepository) GetByInstructor(ctx context.Context, instructorID string) ([]*entity.Course, error) {
	args := m.Called(ctx, instructorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Course), args.Error(1)
}

func TestRequireCourseOwnership(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	middleware := RequireCourseOwnership(mockCourseRepo)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	courseID := uuid.New().String()
	instructorID := uuid.New().String()
	otherUserID := uuid.New().String()
	adminID := uuid.New().String()

	course := &entity.Course{ID: courseID, InstructorID: instructorID}

	// Test case 1: Admin user
	claims := &auth.Claims{UserID: adminID, Role: entity.UserRoleAdmin}
	ctx := context.WithValue(context.Background(), UserContextKey, claims)
	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr := httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test case 2: Course instructor
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil).Once()
	claims = &auth.Claims{UserID: instructorID, Role: entity.UserRoleInstructor}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)

	// Test case 3: Non-instructor user
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil).Once()
	claims = &auth.Claims{UserID: otherUserID, Role: entity.UserRoleStudent}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockCourseRepo.AssertExpectations(t)

	// Test case 4: Course not found
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(nil, entity.ErrNotFound).Once()
	claims = &auth.Claims{UserID: instructorID, Role: entity.UserRoleInstructor}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockCourseRepo.AssertExpectations(t)
}

type MockEnrollmentRepository struct {
	mock.Mock
}

func (m *MockEnrollmentRepository) GetByStudentAndCourse(ctx context.Context, studentID, courseID string) (*entity.Enrollment, error) {
	args := m.Called(ctx, studentID, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) Create(ctx context.Context, enrollment *entity.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) Update(ctx context.Context, enrollment *entity.Enrollment) error {
	args := m.Called(ctx, enrollment)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) Delete(ctx context.Context, studentID, courseID string) error {
	args := m.Called(ctx, studentID, courseID)
	return args.Error(0)
}

func (m *MockEnrollmentRepository) GetByStudent(ctx context.Context, studentID string, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	args := m.Called(ctx, studentID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Enrollment), args.Error(1)
}

func (m *MockEnrollmentRepository) GetByCourse(ctx context.Context, courseID string, filter *repository.EnrollmentFilter) ([]*entity.Enrollment, error) {
	args := m.Called(ctx, courseID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Enrollment), args.Error(1)
}

func TestRequireEnrollment(t *testing.T) {
	mockCourseRepo := new(MockCourseRepository)
	mockEnrollmentRepo := new(MockEnrollmentRepository)
	middleware := RequireEnrollment(mockCourseRepo, mockEnrollmentRepo)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	courseID := uuid.New().String()
	studentID := uuid.New().String()
	instructorID := uuid.New().String()
	adminID := uuid.New().String()

	course := &entity.Course{ID: courseID, InstructorID: instructorID}
	enrollment := &entity.Enrollment{StudentID: studentID, CourseID: courseID, Status: entity.EnrollmentStatusActive}

	// Test case 1: Admin user
	claims := &auth.Claims{UserID: adminID, Role: entity.UserRoleAdmin}
	ctx := context.WithValue(context.Background(), UserContextKey, claims)
	req := httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr := httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test case 2: Course instructor
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil).Once()
	claims = &auth.Claims{UserID: instructorID, Role: entity.UserRoleInstructor}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)

	// Test case 3: Enrolled student
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil).Once()
	mockEnrollmentRepo.On("GetByStudentAndCourse", mock.Anything, studentID, courseID).Return(enrollment, nil).Once()
	claims = &auth.Claims{UserID: studentID, Role: entity.UserRoleStudent}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)

	// Test case 4: Not enrolled student
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil).Once()
	mockEnrollmentRepo.On("GetByStudentAndCourse", mock.Anything, studentID, courseID).Return(nil, entity.ErrNotFound).Once()
	claims = &auth.Claims{UserID: studentID, Role: entity.UserRoleStudent}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/courses/"+courseID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": courseID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockCourseRepo.AssertExpectations(t)
	mockEnrollmentRepo.AssertExpectations(t)
}

type MockModuleRepository struct {
	mock.Mock
}

func (m *MockModuleRepository) GetByID(ctx context.Context, id string) (*entity.Module, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Module), args.Error(1)
}

func (m *MockModuleRepository) Create(ctx context.Context, module *entity.Module) error {
	args := m.Called(ctx, module)
	return args.Error(0)
}

func (m *MockModuleRepository) Update(ctx context.Context, module *entity.Module) error {
	args := m.Called(ctx, module)
	return args.Error(0)
}

func (m *MockModuleRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockModuleRepository) GetNextOrder(ctx context.Context, courseID string) (int, error) {
	args := m.Called(ctx, courseID)
	return args.Int(0), args.Error(1)
}

func (m *MockModuleRepository) GetByCourse(ctx context.Context, courseID string) ([]*entity.Module, error) {
	args := m.Called(ctx, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Module), args.Error(1)
}

func TestRequireModuleOwnership(t *testing.T) {
	mockModuleRepo := new(MockModuleRepository)
	mockCourseRepo := new(MockCourseRepository)
	middleware := RequireModuleOwnership(mockModuleRepo, mockCourseRepo)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	moduleID := uuid.New().String()
	courseID := uuid.New().String()
	instructorID := uuid.New().String()
	otherUserID := uuid.New().String()
	adminID := uuid.New().String()

	module := &entity.Module{ID: moduleID, CourseID: courseID}
	course := &entity.Course{ID: courseID, InstructorID: instructorID}

	// Test case 1: Admin user
	claims := &auth.Claims{UserID: adminID, Role: entity.UserRoleAdmin}
	ctx := context.WithValue(context.Background(), UserContextKey, claims)
	req := httptest.NewRequest(http.MethodGet, "/modules/"+moduleID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": moduleID})
	rr := httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test case 2: Course instructor
	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil).Once()
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil).Once()
	claims = &auth.Claims{UserID: instructorID, Role: entity.UserRoleInstructor}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/modules/"+moduleID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": moduleID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)

	// Test case 3: Non-instructor user
	mockModuleRepo.On("GetByID", mock.Anything, moduleID).Return(module, nil).Once()
	mockCourseRepo.On("GetByID", mock.Anything, courseID).Return(course, nil).Once()
	claims = &auth.Claims{UserID: otherUserID, Role: entity.UserRoleStudent}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/modules/"+moduleID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": moduleID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	mockModuleRepo.AssertExpectations(t)
	mockCourseRepo.AssertExpectations(t)
}

func TestRequireAdminOrSelf(t *testing.T) {
	middleware := RequireAdminOrSelf()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	userID := uuid.New().String()
	otherUserID := uuid.New().String()
	adminID := uuid.New().String()

	// Test case 1: Admin user
	claims := &auth.Claims{UserID: adminID, Role: entity.UserRoleAdmin}
	ctx := context.WithValue(context.Background(), UserContextKey, claims)
	req := httptest.NewRequest(http.MethodGet, "/users/"+otherUserID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": otherUserID})
	rr := httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test case 2: Accessing own resource
	claims = &auth.Claims{UserID: userID, Role: entity.UserRoleStudent}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/users/"+userID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": userID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test case 3: Accessing other user's resource
	claims = &auth.Claims{UserID: otherUserID, Role: entity.UserRoleStudent}
	ctx = context.WithValue(context.Background(), UserContextKey, claims)
	req = httptest.NewRequest(http.MethodGet, "/users/"+userID, nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": userID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Test case 4: No user in context
	req = httptest.NewRequest(http.MethodGet, "/users/"+userID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": userID})
	rr = httptest.NewRecorder()
	middleware(handler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetUserFromContext(t *testing.T) {
	// Test case 1: User in context
	claims := &auth.Claims{UserID: uuid.New().String(), Role: entity.UserRoleStudent}
	ctx := context.WithValue(context.Background(), UserContextKey, claims)
	retrievedClaims, ok := GetUserFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, claims, retrievedClaims)

	// Test case 2: No user in context
	ctx = context.Background()
	retrievedClaims, ok = GetUserFromContext(ctx)
	assert.False(t, ok)
	assert.Nil(t, retrievedClaims)

	// Test case 3: Wrong type in context
	ctx = context.WithValue(context.Background(), UserContextKey, "not a claims object")
	retrievedClaims, ok = GetUserFromContext(ctx)
	assert.False(t, ok)
	assert.Nil(t, retrievedClaims)
}
