package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, code int, err error, message string) {
	RespondWithJSON(w, code, dto.NewErrorResponse(err, message))
}

func MapEntityToUserDTO(user *entity.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func MapEntityToCourseDTO(course *entity.Course) dto.CourseDTO {
	return dto.CourseDTO{
		ID:              course.ID.String(),
		Title:           course.Title,
		Description:     course.Description,
		InstructorID:    course.InstructorID.String(),
		DifficultyLevel: string(course.DifficultyLevel),
		CreatedAt:       course.CreatedAt,
		UpdatedAt:       course.UpdatedAt,
	}
}

func MapEntityToEnrollmentDTO(enrollment *entity.Enrollment) dto.EnrollmentDTO {
	return dto.EnrollmentDTO{
		StudentID:      enrollment.StudentID.String(),
		CourseID:       enrollment.CourseID.String(),
		EnrollmentDate: enrollment.EnrollmentDate,
		Status:         string(enrollment.Status),
		CreatedAt:      enrollment.CreatedAt,
		UpdatedAt:      enrollment.UpdatedAt,
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func handleUseCaseError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if errors.Is(err, entity.ErrNotFound) {
		respondWithError(w, http.StatusNotFound, "Resource not found")
		return
	}
	if errors.Is(err, entity.ErrInvalidInput) || errors.Is(err, entity.ErrFieldRequired) {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, entity.ErrInsufficientPermissions) {
		respondWithError(w, http.StatusForbidden, "Insufficient permissions")
		return
	}
	if errors.Is(err, entity.ErrAlreadyEnrolled) {
		respondWithError(w, http.StatusConflict, "Already enrolled")
		return
	}
	if errors.Is(err, entity.ErrUnauthorized) {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	if isUUIDParsingError(err) {
		respondWithError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}

// isUUIDParsingError checks if an error is a PostgreSQL UUID parsing error
func isUUIDParsingError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "invalid input syntax for type uuid") ||
		strings.Contains(errStr, "22P02") ||
		strings.Contains(errStr, "invalid UUID format")
}
