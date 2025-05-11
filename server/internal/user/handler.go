package user

import (
	"common"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"user/dto"
)

// CreateUserHandler Handler function for create user endpoint
func CreateUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ValidateCreateUserRequest(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := service.CreateUser(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// GetCurrentUserHandler Handler function for get current user endpoint
func GetCurrentUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userClaims, ok := r.Context().Value(USER_CLAIMS_KEY).(*common.UserClaims)
		if !ok {
			http.Error(w, `{"error":"User claims not found"}`, http.StatusUnauthorized)
			return
		}

		request := dto.GetUserRequest{
			UserId: &userClaims.ID,
		}

		response, err := service.GetUser(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// GetUserHandler Handler function for get user endpoint
func GetUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := generateGetUserRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if err := ValidateGetUserRequest(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := service.GetUser(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// GetUsersHandler Handler function for get users endpoint
func GetUsersHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := generateGetUsersRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if err := ValidateGetUsersRequest(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := service.GetUsers(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// UpdateUserHandler Handler function for update user endpoint
func UpdateUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request dto.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ValidateUpdateUserRequest(request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := service.UpdateUser(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// DeleteUserHandler Handler function for delete user endpoint
func DeleteUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := dto.DeleteUserRequest{}
		userIdStr := chi.URLParam(r, "id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		request.UserId = userId
		response, err := service.DeleteUser(r.Context(), request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		}
	}
}

// generateGetUserRequest Populate and return GetUserRequest
func generateGetUserRequest(r *http.Request) (dto.GetUserRequest, error) {
	query := r.URL.Query()
	var request dto.GetUserRequest

	if userIdStr := query.Get("user_id"); userIdStr != "" {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			return dto.GetUserRequest{}, err
		}
		request.UserId = &userId
	}

	if username := query.Get("username"); username != "" {
		request.Username = &username
	}

	if email := query.Get("email"); email != "" {
		request.Email = &email
	}

	return request, nil
}

// generateGetUsersRequest Populate and return GetUsersRequest
func generateGetUsersRequest(r *http.Request) (dto.GetUsersRequest, error) {
	query := r.URL.Query()
	var request dto.GetUsersRequest

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return dto.GetUsersRequest{}, err
		}
		request.Limit = &limit
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return dto.GetUsersRequest{}, err
		}
		request.Offset = &offset
	}

	if sortField := query.Get("sort_field"); sortField != "" {
		request.SortField = &sortField
	}

	if sortOrder := query.Get("sort_order"); sortOrder != "" {
		request.SortOrder = &sortOrder
	}

	return request, nil
}
