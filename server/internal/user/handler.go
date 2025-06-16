package user

import (
    "common"
    "encoding/json"
    "errors"
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
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        }
    }
}

// GetCurrentUserHandler Handler function for get current user endpoint
func GetCurrentUserHandler(service Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        userClaims, ok := r.Context().Value(common.UsersClaimKey).(*common.UserClaims)
        if !ok {
            http.Error(w, "User claims not found", http.StatusUnauthorized)
            return
        }

        request := dto.GetUserRequest{
            UserId: &userClaims.ID,
        }

        response, err := service.GetUser(r.Context(), &request)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        if err := json.NewEncoder(w).Encode(response); err != nil {
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        }
    }
}

// GetUserHandler Handler function for get user endpoint
func GetUserHandler(service Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        request, err := generateGetUserRequest(r)
        if err != nil {
            handleError(err, w)
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
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        }
    }
}

// GetUsersHandler Handler function for get users endpoint
func GetUsersHandler(service Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        request, err := generateGetUsersRequest(r)
        if err != nil {
            handleError(err, w)
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
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
            http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        }
    }
}

// generateGetUserRequest Populate and return GetUserRequest
func generateGetUserRequest(r *http.Request) (*dto.GetUserRequest, error) {
    query := r.URL.Query()
    var request dto.GetUserRequest

    if userIdStr := query.Get("id"); userIdStr != "" {
        userId, err := strconv.Atoi(userIdStr)
        if err != nil {
            return nil, &common.HTTPError{
                StatusCode: http.StatusBadRequest,
                Message:    "Invalid user id",
            }
        }
        request.UserId = &userId
    }

    if username := query.Get("username"); username != "" {
        request.Username = &username
    }

    if email := query.Get("email"); email != "" {
        request.Email = &email
    }

    return &request, nil
}

// generateGetUsersRequest Populate and return GetUsersRequest
func generateGetUsersRequest(r *http.Request) (*dto.GetUsersRequest, error) {
    query := r.URL.Query()
    var request dto.GetUsersRequest

    if limitStr := query.Get("limit"); limitStr != "" {
        limit, err := strconv.Atoi(limitStr)
        if err != nil {
            return nil, &common.HTTPError{
                StatusCode: http.StatusBadRequest,
                Message:    "Invalid limit",
            }
        }
        request.Limit = &limit
    }

    if offsetStr := query.Get("offset"); offsetStr != "" {
        offset, err := strconv.Atoi(offsetStr)
        if err != nil {
            return nil, &common.HTTPError{
                StatusCode: http.StatusBadRequest,
                Message:    "Invalid offset",
            }
        }
        request.Offset = &offset
    }

    if sortField := query.Get("sort_field"); sortField != "" {
        request.SortField = &sortField
    }

    if sortOrder := query.Get("sort_order"); sortOrder != "" {
        request.SortOrder = &sortOrder
    }

    return &request, nil
}

// generateUpdateUserRequest Populate and return UpdateUserRequest
func generateUpdateUserRequest(r *http.Request) (*dto.UpdateUserRequest, error) {
    var request dto.UpdateUserRequest

    userId, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        return nil, &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "Invalid user id",
        }
    }
    request.UserId = userId

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        return nil, &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "Invalid request body: " + err.Error(),
        }
    }

    return &request, nil
}

// generateDeleteUserRequest Populate and return DeleteUserRequest
func generateDeleteUserRequest(r *http.Request) (*dto.DeleteUserRequest, error) {
    var request dto.DeleteUserRequest

    userId, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        return nil, &common.HTTPError{
            StatusCode: http.StatusBadRequest,
            Message:    "Invalid user id",
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
