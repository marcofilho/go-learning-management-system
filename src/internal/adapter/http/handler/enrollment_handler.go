package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	vars := mux.Vars(r)
	courseID := vars["id"]

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.EnrollCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If no body provided, assume self-enrollment
		req.StudentID = userID
	}

	// If StudentID not provided or empty, use self-enrollment
	if req.StudentID == "" {
		req.StudentID = userID
	}

	enrollment, err := h.enrollmentUseCase.EnrollStudent(r.Context(), req.StudentID, courseID, userID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, MapEntityToEnrollmentDTO(enrollment))
}

// GET /api/students/{id}/courses
func (h *EnrollmentHandler) GetStudentCourses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	filter := parseEnrollmentFilter(r)
	enrollments, err := h.enrollmentUseCase.GetStudentCourses(r.Context(), studentID, filter)
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
	vars := mux.Vars(r)
	courseID := vars["id"]

	filter := parseEnrollmentFilter(r)
	enrollments, err := h.enrollmentUseCase.GetCourseStudents(r.Context(), courseID, filter)
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
