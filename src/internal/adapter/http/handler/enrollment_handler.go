package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/repository"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
)

type EnrollmentHandler struct {
	enrollmentUseCase *usecase.EnrollmentUseCase
}

func NewEnrollmentHandler(enrollmentUseCase *usecase.EnrollmentUseCase) *EnrollmentHandler {
	return &EnrollmentHandler{enrollmentUseCase: enrollmentUseCase}
}

// POST /api/courses/{id}/enroll

func (h *EnrollmentHandler) EnrollInCourse(w http.ResponseWriter, r *http.Request) {
	// EnrollInCourse godoc
	// @Summary Enroll in a course
	// @Description Enroll a student in a course
	// @Tags Enrollments
	// @Accept json
	// @Produce json
	// @Param id path string true "Course ID"
	// @Success 201 {object} dto.SuccessResponse{data=dto.EnrollmentDTO}
	// @Failure 400 {object} dto.ErrorResponse
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 409 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /courses/{id}/enroll [post]
	vars := mux.Vars(r)
	courseID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid course ID format")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.EnrollCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If no body provided, assume self-enrollment
		req.StudentID = userID.String()
	}

	// If StudentID not provided or empty, use self-enrollment
	if req.StudentID == "" {
		req.StudentID = userID.String()
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid student ID format")
		return
	}

	enrollment, err := h.enrollmentUseCase.EnrollStudent(r.Context(), studentID.String(), courseID.String(), userID.String())
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, MapEntityToEnrollmentDTO(enrollment))
}

// GET /api/students/{id}/courses

func (h *EnrollmentHandler) GetStudentCourses(w http.ResponseWriter, r *http.Request) {
	// GetStudentCourses godoc
	// @Summary Get student courses
	// @Description Get all courses a student is enrolled in
	// @Tags Enrollments
	// @Accept json
	// @Produce json
	// @Param id path string true "Student ID"
	// @Param status query string false "Filter by enrollment status (active/dropped/completed)"
	// @Success 200 {object} dto.SuccessResponse{data=[]dto.EnrollmentDTO}
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /students/{id}/courses [get]
	vars := mux.Vars(r)
	studentID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid student ID format")
		return
	}

	filter := parseEnrollmentFilter(r)
	enrollments, err := h.enrollmentUseCase.GetStudentCourses(r.Context(), studentID.String(), filter)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.EnrollmentDTO
	for _, enrollment := range enrollments {
		response = append(response, MapEntityToEnrollmentDTO(enrollment))
	}

	respondWithJSON(w, http.StatusOK, response)
}

// GET /api/courses/{id}/students

func (h *EnrollmentHandler) GetCourseStudents(w http.ResponseWriter, r *http.Request) {
	// GetCourseStudents godoc
	// @Summary Get course students
	// @Description Get all students enrolled in a course
	// @Tags Enrollments
	// @Accept json
	// @Produce json
	// @Param id path string true "Course ID"
	// @Param status query string false "Filter by enrollment status (active/dropped/completed)"
	// @Success 200 {object} dto.SuccessResponse{data=[]dto.EnrollmentDTO}
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /courses/{id}/students [get]
	vars := mux.Vars(r)
	courseID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid course ID format")
		return
	}

	filter := parseEnrollmentFilter(r)
	enrollments, err := h.enrollmentUseCase.GetCourseStudents(r.Context(), courseID.String(), filter)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.EnrollmentDTO
	for _, enrollment := range enrollments {
		response = append(response, MapEntityToEnrollmentDTO(enrollment))
	}

	respondWithJSON(w, http.StatusOK, response)
}

func parseEnrollmentFilter(r *http.Request) *repository.EnrollmentFilter {
	filter := &repository.EnrollmentFilter{
		Limit:  10,
		Offset: 0,
	}

	if status := r.URL.Query().Get("status"); status != "" {
		s := entity.EnrollmentStatus(status)
		filter.Status = &s
	}

	if dateFrom := r.URL.Query().Get("date_from"); dateFrom != "" {
		filter.DateFrom = &dateFrom
	}

	if dateTo := r.URL.Query().Get("date_to"); dateTo != "" {
		filter.DateTo = &dateTo
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	return filter
}
