package user

import (
	"common"
	api "common/api/user"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// CreateUserHandler Handler function for create user endpoint
func CreateUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request api.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			handleError(err, w)
			return
		}

		if err := ValidateCreateUserRequest(&request, service, r.Context()); err != nil {
			handleError(err, w)
			return
		}

		response, err := service.CreateUser(r.Context(), &request)
		if err != nil {
			handleError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// GetCurrentUserHandler Handler function for get current user endpoint
func GetCurrentUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userClaims, ok := r.Context().Value(common.UserClaimsCtxKey).(*common.UserClaims)
		if !ok {
			http.Error(w, "user claims not found", http.StatusUnauthorized)
			return
		}

		request := api.GetUserRequest{
			UserId: &userClaims.ID,
		}

		response, err := service.GetUser(r.Context(), &request)
		if err != nil {
			http.Error(w, "unable to retrieve user", http.StatusInternalServerError)
			return
		} else if response == nil {
			http.Error(w, "user not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// GetUserHandler Handler function for get user endpoint
func GetUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := generateGetUserRequest(r)
		if err != nil {
			handleError(err, w)
			return
		}

		if err := ValidateGetUserRequest(request); err != nil {
			handleError(err, w)
			return
		}

		response, err := service.GetUser(r.Context(), request)
		if err != nil {
			handleError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// GetUsersHandler Handler function for get users endpoint
func GetUsersHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := generateGetUsersRequest(r)
		if err != nil {
			handleError(err, w)
			return
		}

		if err := ValidateGetUsersRequest(request); err != nil {
			handleError(err, w)
			return
		}

		response, err := service.GetUsers(r.Context(), request)
		if err != nil {
			handleError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// UpdateUserHandler Handler function for update user endpoint
func UpdateUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := generateUpdateUserRequest(r)
		if err != nil {
			handleError(err, w)
			return
		}

		if err := ValidateUpdateUserRequest(request, service, r.Context()); err != nil {
			handleError(err, w)
			return
		}

		response, err := service.UpdateUser(r.Context(), request)
		if err != nil {
			handleError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// DeleteUserHandler Handler function for delete user endpoint
func DeleteUserHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := generateDeleteUserRequest(r)
		if err != nil {
			handleError(err, w)
			return
		}

		if err := ValidateDeleteUserRequest(request, service, r.Context()); err != nil {
			handleError(err, w)
			return
		}

		response, err := service.DeleteUser(r.Context(), request)
		if err != nil {
			handleError(err, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

// generateGetUserRequest Populate and return GetUserRequest
func generateGetUserRequest(r *http.Request) (*api.GetUserRequest, error) {
	query := r.URL.Query()
	var request api.GetUserRequest

	if userIdStr := query.Get(api.UserIdKey); userIdStr != "" {
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			return nil, &common.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid user id",
			}
		}
		request.UserId = &userId
	}

	if username := query.Get(api.UsernameKey); username != "" {
		request.Username = &username
	}

	if email := query.Get(api.EmailKey); email != "" {
		request.Email = &email
	}

	return &request, nil
}

// generateGetUsersRequest Populate and return GetUsersRequest
func generateGetUsersRequest(r *http.Request) (*api.GetUsersRequest, error) {
	query := r.URL.Query()
	var request api.GetUsersRequest

	if limitStr := query.Get(api.LimitKey); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, &common.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid limit",
			}
		}
		request.Limit = &limit
	}

	if offsetStr := query.Get(api.OffsetKey); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, &common.HTTPError{
				StatusCode: http.StatusBadRequest,
				Message:    "invalid offset",
			}
		}
		request.Offset = &offset
	}

	if sortField := query.Get(api.SortFieldKey); sortField != "" {
		request.SortField = &sortField
	}

	if sortOrder := query.Get(api.SortDirectionKey); sortOrder != "" {
		request.SortDirection = &sortOrder
	}

	return &request, nil
}

// generateUpdateUserRequest Populate and return UpdateUserRequest
func generateUpdateUserRequest(r *http.Request) (*api.UpdateUserRequest, error) {
	var request api.UpdateUserRequest

	userId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return nil, &common.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid user id",
		}
	}
	request.UserId = userId

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &common.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid request body: " + err.Error(),
		}
	}

	return &request, nil
}

// generateDeleteUserRequest Populate and return DeleteUserRequest
func generateDeleteUserRequest(r *http.Request) (*api.DeleteUserRequest, error) {
	var request api.DeleteUserRequest

	userId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return nil, &common.HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid user id",
		}
	}
	request.UserId = userId

	return &request, nil
}

// handleError Write the appropriate response given an error
func handleError(err error, w http.ResponseWriter) {
	var httpErr *common.HTTPError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Error(), httpErr.StatusCode)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
