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

type CourseHandler struct {
	courseUseCase *usecase.CourseUseCase
}

func NewCourseHandler(courseUseCase *usecase.CourseUseCase) *CourseHandler {
	return &CourseHandler{courseUseCase: courseUseCase}
}

// CreateCourse godoc
// @Summary Create a new course
// @Description Create a new course (instructor only)
// @Tags Courses
// @Accept json
// @Produce json
// @Param request body dto.CreateCourseRequest true "Course details"
// @Success 201 {object} dto.SuccessResponse{data=dto.CourseDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /courses [post]
func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.CreateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	difficultyLevel := entity.DifficultyLevel(req.DifficultyLevel)
	if difficultyLevel == "" {
		difficultyLevel = entity.DifficultyLevelBeginner
	}

	course, err := h.courseUseCase.CreateCourse(r.Context(), req.Title, req.Description, userID, difficultyLevel)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, MapEntityToCourseDTO(course))
}

// GetCourse godoc
// @Summary Get course by ID
// @Description Get course details by ID
// @Tags Courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.CourseDTO}
// @Failure 404 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	course, err := h.courseUseCase.GetCourseByID(r.Context(), courseID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, MapEntityToCourseDTO(course))
}

// ListCourses godoc
// @Summary List courses
// @Description List courses with optional filters
// @Tags Courses
// @Accept json
// @Produce json
// @Param instructor_id query string false "Filter by instructor ID"
// @Param difficulty_level query string false "Filter by difficulty level (beginner/intermediate/advanced)"
// @Param active_only query boolean false "Show only active courses"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.SuccessResponse{data=[]dto.CourseDTO}
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /courses [get]
func (h *CourseHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	instructorID := r.URL.Query().Get("instructor_id")
	difficultyLevel := r.URL.Query().Get("difficulty_level")
	activeOnlyStr := r.URL.Query().Get("active_only")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	filter := &repository.CourseFilter{
		Limit:  10,
		Offset: 0,
	}

	if instructorID != "" {
		filter.InstructorID = &instructorID
	}

	if difficultyLevel != "" {
		dl := entity.DifficultyLevel(difficultyLevel)
		filter.DifficultyLevel = &dl
	}

	if activeOnlyStr == "true" {
		filter.ActiveOnly = true
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	courses, err := h.courseUseCase.ListCourses(r.Context(), filter)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.CourseDTO
	for _, course := range courses {
		response = append(response, MapEntityToCourseDTO(course))
	}

	respondWithJSON(w, http.StatusOK, response)
}

// UpdateCourse godoc
// @Summary Update course
// @Description Update course details (instructor only)
// @Tags Courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID"
// @Param request body dto.UpdateCourseRequest true "Course update details"
// @Success 200 {object} dto.SuccessResponse{data=dto.CourseDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	course, err := h.courseUseCase.GetCourseByID(r.Context(), courseID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	// Only instructor can modify their own course
	if course.InstructorID != userID {
		respondWithError(w, http.StatusForbidden, "Only the course instructor can modify this course")
		return
	}

	var req dto.UpdateCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Title != "" {
		course.Title = req.Title
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.DifficultyLevel != "" {
		course.DifficultyLevel = entity.DifficultyLevel(req.DifficultyLevel)
	}

	if err := h.courseUseCase.UpdateCourse(r.Context(), course, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, MapEntityToCourseDTO(course))
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	// DeleteCourse godoc
	// @Summary Delete course
	// @Description Delete a course (instructor only)
	// @Tags Courses
	// @Accept json
	// @Produce json
	// @Param id path string true "Course ID"
	// @Success 204
	// @Failure 401 {object} dto.ErrorResponse
	// @Failure 403 {object} dto.ErrorResponse
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /courses/{id} [delete]
	vars := mux.Vars(r)
	courseID := vars["id"]

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.courseUseCase.DeleteCourse(r.Context(), courseID, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
