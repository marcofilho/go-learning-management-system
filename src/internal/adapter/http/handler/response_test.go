package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/handler"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithJSON(t *testing.T) {
	w := httptest.NewRecorder()
	payload := map[string]string{"message": "success"}
	handler.RespondWithJSON(w, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"success"}`, w.Body.String())
}

func TestRespondWithError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("test error")
	handler.RespondWithError(w, http.StatusInternalServerError, err, "Internal Server Error")

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"error":"test error", "message":"Internal Server Error"}`, w.Body.String())
}

func TestMapEntityToUserDTO(t *testing.T) {
	now := time.Now()
	user := &entity.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      entity.UserRoleStudent,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dto := handler.MapEntityToUserDTO(user)

	assert.Equal(t, user.ID.String(), dto.ID)
	assert.Equal(t, user.Email, dto.Email)
	assert.Equal(t, user.FirstName, dto.FirstName)
	assert.Equal(t, user.LastName, dto.LastName)
	assert.Equal(t, string(user.Role), dto.Role)
	assert.Equal(t, user.IsActive, dto.IsActive)
	assert.Equal(t, user.CreatedAt, dto.CreatedAt)
	assert.Equal(t, user.UpdatedAt, dto.UpdatedAt)
}

func TestMapEntityToCourseDTO(t *testing.T) {
	now := time.Now()
	course := &entity.Course{
		ID:              uuid.New(),
		Title:           "Test Course",
		Description:     "Test Description",
		InstructorID:    uuid.New(),
		DifficultyLevel: entity.DifficultyLevelBeginner,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	dto := handler.MapEntityToCourseDTO(course)

	assert.Equal(t, course.ID.String(), dto.ID)
	assert.Equal(t, course.Title, dto.Title)
	assert.Equal(t, course.Description, dto.Description)
	assert.Equal(t, course.InstructorID.String(), dto.InstructorID)
	assert.Equal(t, string(course.DifficultyLevel), dto.DifficultyLevel)
	assert.Equal(t, course.CreatedAt, dto.CreatedAt)
	assert.Equal(t, course.UpdatedAt, dto.UpdatedAt)
}

func TestMapEntityToEnrollmentDTO(t *testing.T) {
	now := time.Now()
	enrollment := &entity.Enrollment{
		StudentID:      uuid.New(),
		CourseID:       uuid.New(),
		EnrollmentDate: now,
		Status:         entity.EnrollmentStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	dto := handler.MapEntityToEnrollmentDTO(enrollment)

	assert.Equal(t, enrollment.StudentID.String(), dto.StudentID)
	assert.Equal(t, enrollment.CourseID.String(), dto.CourseID)
	assert.Equal(t, enrollment.EnrollmentDate, dto.EnrollmentDate)
	assert.Equal(t, string(enrollment.Status), dto.Status)
	assert.Equal(t, enrollment.CreatedAt, dto.CreatedAt)
	assert.Equal(t, enrollment.UpdatedAt, dto.UpdatedAt)
}

func TestHandleUseCaseError(t *testing.T) {
	testCases := []struct {
		name         string
		err          error
		expectedCode int
		expectedBody string
	}{
		{
			name:         "ErrNotFound",
			err:          entity.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"Resource not found"}`,
		},
		{
			name:         "ErrUnauthorized",
			err:          entity.ErrUnauthorized,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"Forbidden"}`,
		},
		{
			name:         "ErrInvalidInput",
			err:          entity.ErrInvalidInput,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid input"}`,
		},
		{
			name:         "ErrInsufficientPermissions",
			err:          entity.ErrInsufficientPermissions,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"Insufficient permissions"}`,
		},
		{
			name:         "ErrAlreadyEnrolled",
			err:          entity.ErrAlreadyEnrolled,
			expectedCode: http.StatusConflict,
			expectedBody: `{"error":"Already enrolled"}`,
		},
		{
			name:         "DefaultError",
			err:          errors.New("some other error"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"Internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.HandleUseCaseErrorForTest(w, tc.err)

			assert.Equal(t, tc.expectedCode, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}
