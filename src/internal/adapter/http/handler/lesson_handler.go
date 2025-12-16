package handler

import (
	"encoding/json"
	"net/http"

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

func (h *LessonHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID := vars["moduleId"]

	var req dto.CreateLessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
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
		ID:            lesson.ID,
		ModuleID:      lesson.ModuleID,
		VersionNumber: lesson.VersionNumber,
		Content:       lesson.Content,
		VideoURL:      lesson.VideoURL,
		AttachmentURL: lesson.AttachmentURL,
		CreatedAt:     lesson.CreatedAt,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *LessonHandler) CreateLessonVersion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lessonID := vars["lessonId"]

	var req dto.CreateLessonVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
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
		ID:            newVersion.ID,
		ModuleID:      newVersion.ModuleID,
		VersionNumber: newVersion.VersionNumber,
		Content:       newVersion.Content,
		VideoURL:      newVersion.VideoURL,
		AttachmentURL: newVersion.AttachmentURL,
		CreatedAt:     newVersion.CreatedAt,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *LessonHandler) GetModuleLessons(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID := vars["moduleId"]

	lessons, err := h.lessonUseCase.GetLatestLessonsByModule(r.Context(), moduleID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.LessonVersionDTO
	for _, lesson := range lessons {
		response = append(response, dto.LessonVersionDTO{
			ID:            lesson.ID,
			ModuleID:      lesson.ModuleID,
			VersionNumber: lesson.VersionNumber,
			Content:       lesson.Content,
			VideoURL:      lesson.VideoURL,
			AttachmentURL: lesson.AttachmentURL,
			CreatedAt:     lesson.CreatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *LessonHandler) GetAllLessonVersions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lessonID := vars["lessonId"]

	versions, err := h.lessonUseCase.GetAllLessonVersions(r.Context(), lessonID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.LessonVersionDTO
	for _, version := range versions {
		response = append(response, dto.LessonVersionDTO{
			ID:            version.ID,
			ModuleID:      version.ModuleID,
			VersionNumber: version.VersionNumber,
			Content:       version.Content,
			VideoURL:      version.VideoURL,
			AttachmentURL: version.AttachmentURL,
			CreatedAt:     version.CreatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *LessonHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lessonID := vars["lessonId"]

	userID, ok := r.Context().Value("userID").(string)
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
