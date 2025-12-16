package handler

import (
	"encoding/json"
	"net/http"

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
		ID:        user.ID,
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
		ID:              course.ID,
		Title:           course.Title,
		Description:     course.Description,
		InstructorID:    course.InstructorID,
		DifficultyLevel: string(course.DifficultyLevel),
		CreatedAt:       course.CreatedAt,
		UpdatedAt:       course.UpdatedAt,
	}
}

func MapEntityToEnrollmentDTO(enrollment *entity.Enrollment) dto.EnrollmentDTO {
	return dto.EnrollmentDTO{
		StudentID:      enrollment.StudentID,
		CourseID:       enrollment.CourseID,
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
	switch err {
	case entity.ErrNotFound:
		respondWithError(w, http.StatusNotFound, "Resource not found")
	case entity.ErrUnauthorized:
		respondWithError(w, http.StatusForbidden, "Forbidden")
	case entity.ErrInvalidInput:
		respondWithError(w, http.StatusBadRequest, err.Error())
	case entity.ErrInsufficientPermissions:
		respondWithError(w, http.StatusForbidden, "Insufficient permissions")
	case entity.ErrAlreadyEnrolled:
		respondWithError(w, http.StatusConflict, "Already enrolled")
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}
