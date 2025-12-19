package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
)

type LessonHandler struct {
	lessonUseCase *usecase.LessonUseCase
}

func NewLessonHandler(lessonUseCase *usecase.LessonUseCase) *LessonHandler {
	return &LessonHandler{lessonUseCase: lessonUseCase}
}

// CreateLesson godoc
// @Summary Create a new lesson
// @Description Create a new lesson version 1 for a module (instructor only)
// @Tags Lessons
// @Accept json
// @Produce json
// @Param moduleId path string true "Module ID"
// @Param request body dto.CreateLessonRequest true "Lesson details"
// @Success 201 {object} dto.SuccessResponse{data=dto.LessonVersionDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /modules/{moduleId}/lessons [post]
func (h *LessonHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := uuid.Parse(vars["moduleId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid module ID format")
		return
	}

	var req dto.CreateLessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lesson := &entity.LessonVersion{
		ModuleID:      moduleID,
		Content:       req.Content,
		VideoURL:      req.VideoURL,
		AttachmentURL: req.AttachmentURL,
	}

	if err := h.lessonUseCase.CreateLesson(r.Context(), lesson, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	response := dto.LessonVersionDTO{
		ID:            lesson.ID.String(),
		LessonID:      lesson.LessonID.String(),
		ModuleID:      lesson.ModuleID.String(),
		VersionNumber: lesson.VersionNumber,
		Content:       lesson.Content,
		VideoURL:      lesson.VideoURL,
		AttachmentURL: lesson.AttachmentURL,
		CreatedAt:     lesson.CreatedAt,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// CreateLessonVersion godoc
// @Summary Create a new lesson version
// @Description Create a new version of an existing lesson (instructor only)
// @Tags Lessons
// @Accept json
// @Produce json
// @Param lessonId path string true "Lesson ID (any version)"
// @Param request body dto.CreateLessonVersionRequest true "Lesson version details"
// @Success 201 {object} dto.SuccessResponse{data=dto.LessonVersionDTO}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /lessons/{lessonId}/version [post]
func (h *LessonHandler) CreateLessonVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lessonID, err := uuid.Parse(vars["lessonId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lesson ID format")
		return
	}

	var req dto.CreateLessonVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	newVersion := &entity.LessonVersion{
		Content:       req.Content,
		VideoURL:      req.VideoURL,
		AttachmentURL: req.AttachmentURL,
	}

	if err := h.lessonUseCase.CreateLessonVersion(r.Context(), lessonID, newVersion, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	response := dto.LessonVersionDTO{
		ID:            newVersion.ID.String(),
		LessonID:      newVersion.LessonID.String(),
		ModuleID:      newVersion.ModuleID.String(),
		VersionNumber: newVersion.VersionNumber,
		Content:       newVersion.Content,
		VideoURL:      newVersion.VideoURL,
		AttachmentURL: newVersion.AttachmentURL,
		CreatedAt:     newVersion.CreatedAt,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// GetModuleLessons godoc
// @Summary Get module lessons
// @Description Get latest version of all lessons in a module
// @Tags Lessons
// @Accept json
// @Produce json
// @Param moduleId path string true "Module ID"
// @Success 200 {object} dto.SuccessResponse{data=[]dto.LessonVersionDTO}
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /modules/{moduleId}/lessons [get]
func (h *LessonHandler) GetModuleLessons(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID, err := uuid.Parse(vars["moduleId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid module ID format")
		return
	}

	pg := parsePagination(r)
	lessons, total, err := h.lessonUseCase.GetLatestLessonsByModule(r.Context(), moduleID, pg.limit, pg.offset)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.LessonVersionDTO
	for _, lesson := range lessons {
		response = append(response, dto.LessonVersionDTO{
			ID:            lesson.ID.String(),
			ModuleID:      lesson.ModuleID.String(),
			VersionNumber: lesson.VersionNumber,
			Content:       lesson.Content,
			VideoURL:      lesson.VideoURL,
			AttachmentURL: lesson.AttachmentURL,
			CreatedAt:     lesson.CreatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, paginatedResponse{
		Data:       response,
		Pagination: buildPagination(pg.page, pg.pageSize, total),
	})
}

// GetAllLessonVersions godoc
// @Summary Get all lesson versions
// @Description Get all versions of a specific lesson
// @Tags Lessons
// @Accept json
// @Produce json
// @Param lessonId path string true "Lesson ID (any version)"
// @Success 200 {object} dto.SuccessResponse{data=[]dto.LessonVersionDTO}
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /lessons/{lessonId}/all-versions [get]
func (h *LessonHandler) GetAllLessonVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lessonID, err := uuid.Parse(vars["lessonId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lesson ID format")
		return
	}

	pg := parsePagination(r)
	versions, total, err := h.lessonUseCase.GetAllLessonVersions(r.Context(), lessonID, pg.limit, pg.offset)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.LessonVersionDTO
	for _, version := range versions {
		response = append(response, dto.LessonVersionDTO{
			ID:            version.ID.String(),
			LessonID:      version.LessonID.String(),
			ModuleID:      version.ModuleID.String(),
			VersionNumber: version.VersionNumber,
			Content:       version.Content,
			VideoURL:      version.VideoURL,
			AttachmentURL: version.AttachmentURL,
			CreatedAt:     version.CreatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, paginatedResponse{
		Data:       response,
		Pagination: buildPagination(pg.page, pg.pageSize, total),
	})
}

// DeleteLesson godoc
// @Summary Delete lesson
// @Description Delete all versions of a lesson (instructor only)
// @Tags Lessons
// @Accept json
// @Produce json
// @Param lessonId path string true "Lesson ID (any version)"
// @Success 204
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /lessons/{lessonId} [delete]
func (h *LessonHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lessonID, err := uuid.Parse(vars["lessonId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid lesson ID format")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.lessonUseCase.DeleteLesson(r.Context(), lessonID, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
