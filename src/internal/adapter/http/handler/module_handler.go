package handler

import (
	"encoding/json"
	"net/http"

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
	vars := mux.Vars(r)
	courseID := vars["courseId"]

	var req dto.CreateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
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
		ID:         module.ID,
		CourseID:   module.CourseID,
		Title:      module.Title,
		OrderIndex: module.OrderIndex,
		CreatedAt:  module.CreatedAt,
		UpdatedAt:  module.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *ModuleHandler) GetCourseModules(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID := vars["courseId"]

	modules, err := h.moduleUseCase.GetModulesByCourse(r.Context(), courseID)
	if err != nil {
		handleUseCaseError(w, err)
		return
	}

	var response []dto.ModuleDTO
	for _, module := range modules {
		response = append(response, dto.ModuleDTO{
			ID:         module.ID,
			CourseID:   module.CourseID,
			Title:      module.Title,
			OrderIndex: module.OrderIndex,
			CreatedAt:  module.CreatedAt,
			UpdatedAt:  module.UpdatedAt,
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ModuleHandler) UpdateModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID := vars["id"]

	var req dto.UpdateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, ok := r.Context().Value("userID").(string)
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
		ID:         module.ID,
		CourseID:   module.CourseID,
		Title:      module.Title,
		OrderIndex: module.OrderIndex,
		CreatedAt:  module.CreatedAt,
		UpdatedAt:  module.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *ModuleHandler) DeleteModule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	moduleID := vars["id"]

	userID, ok := r.Context().Value("userID").(string)
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
