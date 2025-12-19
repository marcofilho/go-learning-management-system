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

type ModuleHandler struct {
	moduleUseCase *usecase.ModuleUseCase
}

func NewModuleHandler(moduleUseCase *usecase.ModuleUseCase) *ModuleHandler {
	return &ModuleHandler{moduleUseCase: moduleUseCase}
}

func (h *ModuleHandler) CreateModule(w http.ResponseWriter, r *http.Request) {
	// CreateModule godoc
	// @Summary Create a new module
	// @Description Create a new module for a course (instructor only)
	// @Tags Modules
	// @Accept json
	// @Produce json
	// @Param courseId path string true "Course ID"
	// @Param request body dto.CreateModuleRequest true "Module details"
	// @Success 201 {object} dto.SuccessResponse{data=dto.ModuleDTO}
	// @Failure 400 {object} dto.ErrorResponse
	// @Failure 403 {object} dto.ErrorResponse
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /courses/{courseId}/modules [post]
	vars := mux.Vars(r)
	courseID, err := uuid.Parse(vars["courseId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid course ID format")
		return
	}

	var req dto.CreateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	module := &entity.Module{
		CourseID:   courseID,
		Title:      req.Title,
		OrderIndex: req.OrderIndex,
	}

	if err := h.moduleUseCase.CreateModule(r.Context(), module, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	response := dto.ModuleDTO{
		ID:         module.ID.String(),
		CourseID:   module.CourseID.String(),
		Title:      module.Title,
		OrderIndex: module.OrderIndex,
		CreatedAt:  module.CreatedAt,
		UpdatedAt:  module.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *ModuleHandler) GetCourseModules(w http.ResponseWriter, r *http.Request) {
	// GetCourseModules godoc
	// @Summary Get course modules
	// @Description Get all modules for a course
	// @Tags Modules
	// @Accept json
	// @Produce json
	// @Param courseId path string true "Course ID"
	// @Success 200 {object} dto.SuccessResponse{data=[]dto.ModuleDTO}
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /courses/{courseId}/modules [get]
	vars := mux.Vars(r)
	courseID, err := uuid.Parse(vars["courseId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid course ID format")
		return
	}

	modules, err := h.moduleUseCase.GetModulesByCourse(r.Context(), courseID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.ModuleDTO
	for _, module := range modules {
		response = append(response, dto.ModuleDTO{
			ID:         module.ID.String(),
			CourseID:   module.CourseID.String(),
			Title:      module.Title,
			OrderIndex: module.OrderIndex,
			CreatedAt:  module.CreatedAt,
			UpdatedAt:  module.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ModuleHandler) GetModule(w http.ResponseWriter, r *http.Request) {
	// GetModule godoc
	// @Summary Get module
	// @Description Get a module by ID
	// @Tags Modules
	// @Accept json
	// @Produce json
	// @Param id path string true "Module ID"
	// @Success 200 {object} dto.SuccessResponse{data=dto.ModuleDTO}
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /modules/{id} [get]
	vars := mux.Vars(r)
	moduleID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid module ID format")
		return
	}

	module, err := h.moduleUseCase.GetModule(r.Context(), moduleID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	response := dto.ModuleDTO{
		ID:         module.ID.String(),
		CourseID:   module.CourseID.String(),
		Title:      module.Title,
		OrderIndex: module.OrderIndex,
		CreatedAt:  module.CreatedAt,
		UpdatedAt:  module.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ModuleHandler) UpdateModule(w http.ResponseWriter, r *http.Request) {
	// UpdateModule godoc
	// @Summary Update module
	// @Description Update module details (instructor only)
	// @Tags Modules
	// @Accept json
	// @Produce json
	// @Param id path string true "Module ID"
	// @Param request body dto.UpdateModuleRequest true "Module update details"
	// @Success 200 {object} dto.SuccessResponse{data=dto.ModuleDTO}
	// @Failure 400 {object} dto.ErrorResponse
	// @Failure 403 {object} dto.ErrorResponse
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /modules/{id} [put]
	vars := mux.Vars(r)
	moduleID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid module ID format")
		return
	}

	var req dto.UpdateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	module, err := h.moduleUseCase.GetModule(r.Context(), moduleID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	if req.Title != "" {
		module.Title = req.Title
	}
	if req.OrderIndex >= 0 {
		module.OrderIndex = req.OrderIndex
	}

	if err := h.moduleUseCase.UpdateModule(r.Context(), module, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	response := dto.ModuleDTO{
		ID:         module.ID.String(),
		CourseID:   module.CourseID.String(),
		Title:      module.Title,
		OrderIndex: module.OrderIndex,
		CreatedAt:  module.CreatedAt,
		UpdatedAt:  module.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ModuleHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	// DeleteModule godoc
	// @Summary Delete module
	// @Description Delete a module (instructor only)
	// @Tags Modules
	// @Accept json
	// @Produce json
	// @Param id path string true "Module ID"
	// @Success 204
	// @Failure 403 {object} dto.ErrorResponse
	// @Failure 404 {object} dto.ErrorResponse
	// @Failure 500 {object} dto.ErrorResponse
	// @Security BearerAuth
	// @Router /modules/{id} [delete]
	vars := mux.Vars(r)
	moduleID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid module ID format")
		return
	}

	userID, ok := GetUserIDFromContext(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.moduleUseCase.DeleteModule(r.Context(), moduleID, userID); err != nil {
		handleUseCaseError(w, err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
