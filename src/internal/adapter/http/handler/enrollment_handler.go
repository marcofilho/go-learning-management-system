package handler

import (
	"encoding/json"
	"net/http"

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
func (h *EnrollmentHandler) EnrollInCourse(w http.ResponseWriter, r *http.Request) {
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
		req.StudentID = userID.String()
	}

	if req.StudentID == "" {
		req.StudentID = userID.String()
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid student ID format")
		return
	}

	enrollment, err := h.enrollmentUseCase.EnrollStudent(r.Context(), studentID, courseID, userID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, MapEntityToEnrollmentDTO(enrollment))
}

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
func (h *EnrollmentHandler) GetStudentCourses(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid student ID format")
		return
	}

	pg := parsePagination(r)
	filter := parseEnrollmentFilter(r, pg)
	enrollments, total, err := h.enrollmentUseCase.GetStudentCourses(r.Context(), studentID, filter)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.EnrollmentDTO
	for _, enrollment := range enrollments {
		response = append(response, MapEntityToEnrollmentDTO(enrollment))
	}

	respondWithJSON(w, http.StatusOK, paginatedResponse{
		Data:       response,
		Pagination: buildPagination(pg.page, pg.pageSize, total),
	})
}

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
func (h *EnrollmentHandler) GetCourseStudents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid course ID format")
		return
	}

	pg := parsePagination(r)
	filter := parseEnrollmentFilter(r, pg)
	enrollments, total, err := h.enrollmentUseCase.GetCourseStudents(r.Context(), courseID, filter)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.EnrollmentDTO
	for _, enrollment := range enrollments {
		response = append(response, MapEntityToEnrollmentDTO(enrollment))
	}

	respondWithJSON(w, http.StatusOK, paginatedResponse{
		Data:       response,
		Pagination: buildPagination(pg.page, pg.pageSize, total),
	})
}

// UpdateEnrollmentStatus godoc
// @Summary Update enrollment status
// @Description Update the status of an enrollment (student can update own, admin can update any)
// @Tags Enrollments
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID"
// @Param request body dto.UpdateEnrollmentStatusRequest true "Update enrollment status request"
// @Success 200 {object} dto.SuccessResponse{data=dto.EnrollmentDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /enrollments/{courseId}/status [put]
func (h *EnrollmentHandler) UpdateEnrollmentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := uuid.Parse(vars["courseId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid course ID format")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.UpdateEnrollmentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if !validateAndRespond(w, &req) {
		return
	}

	studentID := userID
	if req.StudentID != "" {
		parsedStudentID, err := uuid.Parse(req.StudentID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid student ID format")
			return
		}
		studentID = parsedStudentID
	}

	status := entity.EnrollmentStatus(req.Status)

	updatedEnrollment, err := h.enrollmentUseCase.UpdateEnrollmentStatus(r.Context(), studentID, courseID, status, userID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, MapEntityToEnrollmentDTO(updatedEnrollment))
}

// DropEnrollment godoc
// @Summary Drop enrollment
// @Description Drop an enrollment (set status to dropped). Student can drop own, admin can drop any.
// @Tags Enrollments
// @Accept json
// @Produce json
// @Param courseId path string true "Course ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.EnrollmentDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /enrollments/{courseId} [delete]
func (h *EnrollmentHandler) DropEnrollment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := uuid.Parse(vars["courseId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid course ID format")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	studentID := userID

	if studentIDParam := r.URL.Query().Get("student_id"); studentIDParam != "" {
		parsedStudentID, err := uuid.Parse(studentIDParam)
		if err == nil {
			studentID = parsedStudentID
		}
	}

	updatedEnrollment, err := h.enrollmentUseCase.DropEnrollment(r.Context(), studentID, courseID, userID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, MapEntityToEnrollmentDTO(updatedEnrollment))
}

func parseEnrollmentFilter(r *http.Request, pg paginationParams) *repository.EnrollmentFilter {
	filter := &repository.EnrollmentFilter{
		Limit:  pg.limit,
		Offset: pg.offset,
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

	return filter
}
