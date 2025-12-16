package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/middleware"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/domain/entity"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/infrastructure/auth"
	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/usecase"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}
	role := entity.UserRole(req.Role)
	user, err := h.userUseCase.Register(r.Context(), req.Email, req.Password, req.FirstName, req.LastName, role)
	if err != nil {
		if err == entity.ErrDuplicateEntry {
			RespondWithError(w, http.StatusConflict, err, "Email already exists")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err, "Failed to register user")
		return
	}
	userDTO := MapEntityToUserDTO(user)
	RespondWithJSON(w, http.StatusCreated, dto.SuccessResponse{
		Message: "User registered successfully",
		Data:    userDTO,
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}
	token, user, err := h.userUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == entity.ErrInvalidCredentials {
			RespondWithError(w, http.StatusUnauthorized, err, "Invalid credentials")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err, "Failed to login")
		return
	}
	userDTO := MapEntityToUserDTO(user)
	RespondWithJSON(w, http.StatusOK, dto.LoginResponse{
		Token: token,
		User:  userDTO,
	})
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	user, err := h.userUseCase.GetUserByID(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err, "User not found")
		return
	}
	userDTO := MapEntityToUserDTO(user)
	RespondWithJSON(w, http.StatusOK, userDTO)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit := 10
	offset := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}
	users, err := h.userUseCase.ListUsers(r.Context(), limit, offset)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Failed to list users")
		return
	}
	var userDTOs []dto.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, MapEntityToUserDTO(user))
	}
	RespondWithJSON(w, http.StatusOK, userDTOs)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	claims, ok := r.Context().Value(middleware.UserContextKey).(*auth.Claims)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, entity.ErrUnauthorized, "User not authenticated")
		return
	}
	if claims.UserID != id && claims.Role != entity.UserRoleAdmin {
		RespondWithError(w, http.StatusForbidden, entity.ErrInsufficientPermissions, "Insufficient permissions")
		return
	}
	user, err := h.userUseCase.GetUserByID(r.Context(), id)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err, "User not found")
		return
	}
	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if err := h.userUseCase.UpdateUser(r.Context(), user); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Failed to update user")
		return
	}
	userDTO := MapEntityToUserDTO(user)
	RespondWithJSON(w, http.StatusOK, userDTO)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := h.userUseCase.DeleteUser(r.Context(), id); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err, "Failed to delete user")
		return
	}
	RespondWithJSON(w, http.StatusOK, dto.SuccessResponse{
		Message: "User deleted successfully",
	})
}
