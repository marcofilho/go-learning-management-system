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

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
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

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["id"]

	userID, ok := r.Context().Value("userID").(string)
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
	vars := mux.Vars(r)
	courseID := vars["id"]

	userID, ok := r.Context().Value("userID").(string)
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
