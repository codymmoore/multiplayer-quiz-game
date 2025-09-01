package user

import (
    api "common/api/user"
    "context"
    "net/http"
    "testing"
)

func TestValidateCreateUserRequest_Success(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return nil, nil
        },
    }
    request := api.CreateUserRequest{
        Username: ValidUsername,
        Email:    ValidEmail,
        Password: ValidPassword,
    }

    if err := ValidateCreateUserRequest(&request, service, nil); err != nil {
        t.Errorf(`ValidateCreateUserRequest(&request, service, nil) = "%v", expected "<nil>"`, err)
    }
}

func TestValidateCreateUserRequest_InvalidUsername(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return nil, nil
        },
    }
    request := api.CreateUserRequest{
        Username: "%$#@!^&*()",
        Email:    ValidEmail,
        Password: ValidPassword,
    }

    err := ValidateCreateUserRequest(&request, service, nil)
    if err == nil {
        t.Error(`ValidateCreateUserRequest(&request, service, nil) = "<nil>", expected non-nil`)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateCreateUserRequest_DuplicateUsername(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.Username != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }
    request := api.CreateUserRequest{
        Username: ValidUsername,
        Email:    ValidEmail,
        Password: ValidPassword,
    }

    err := ValidateCreateUserRequest(&request, service, nil)
    if err == nil {
        t.Error(`ValidateCreateUserRequest(&request, service, nil) = "<nil>", expected non-nil`)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateCreateUserRequest_InvalidEmail(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return nil, nil
        },
    }
    request := api.CreateUserRequest{
        Username: ValidUsername,
        Email:    "uh",
        Password: ValidPassword,
    }

    err := ValidateCreateUserRequest(&request, service, nil)
    if err == nil {
        t.Error(`ValidateCreateUserRequest(&request, service, nil) = "<nil>", expected non-nil`)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateCreateUserRequest_DuplicateEmail(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.Email != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }
    request := api.CreateUserRequest{
        Username: ValidUsername,
        Email:    ValidEmail,
        Password: ValidPassword,
    }

    err := ValidateCreateUserRequest(&request, service, nil)
    if err == nil {
        t.Error(`ValidateCreateUserRequest(&request, service, nil) = "<nil>", expected non-nil`)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateCreateUserRequest_InvalidPassword(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return nil, nil
        },
    }
    request := api.CreateUserRequest{
        Username: ValidUsername,
        Email:    ValidEmail,
        Password: "invalidPassword",
    }

    err := ValidateCreateUserRequest(&request, service, nil)
    if err == nil {
        t.Error(`ValidateCreateUserRequest(&request, service, nil) = "<nil>", expected non-nil`)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUserRequest_Success(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail

    request := api.GetUserRequest{
        UserId:   &userId,
        Username: &username,
        Email:    &email,
    }

    if err := ValidateGetUserRequest(&request); err != nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "<nil>"`, err)
    }
}

func TestValidateGetUserRequest_MissingParameter(t *testing.T) {
    request := api.GetUserRequest{}
    err := ValidateGetUserRequest(&request)
    if err == nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "ID, username, or email is required"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUserRequest_InvalidUserId(t *testing.T) {
    userId := -1
    request := api.GetUserRequest{
        UserId: &userId,
    }

    err := ValidateGetUserRequest(&request)
    if err == nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "invalid user id"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_Success(t *testing.T) {
    limit := 1
    offset := 1
    sortField := "CreatedAt"
    sortDirection := "desc"

    request := api.GetUsersRequest{
        Limit:         &limit,
        Offset:        &offset,
        SortField:     &sortField,
        SortDirection: &sortDirection,
    }

    if err := ValidateGetUsersRequest(&request); err != nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "<nil>"`, err)
    }
}

func TestValidateGetUsersRequest_InvalidLimit(t *testing.T) {
    limit := -1

    request := api.GetUsersRequest{
        Limit: &limit,
    }

    err := ValidateGetUsersRequest(&request)
    if err == nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "limit must be a positive number"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_InvalidOffset(t *testing.T) {
    offset := -1

    request := api.GetUsersRequest{
        Offset: &offset,
    }

    err := ValidateGetUsersRequest(&request)
    if err == nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "offset must be a positive number"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_InvalidSortField(t *testing.T) {
    sortField := "invalidField"

    request := api.GetUsersRequest{
        SortField: &sortField,
    }

    err := ValidateGetUsersRequest(&request)
    if err == nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "invalid sort field"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateGetUsersRequest_InvalidSortDirection(t *testing.T) {
    sortDirection := "invalidDirection"

    request := api.GetUsersRequest{
        SortDirection: &sortDirection,
    }

    err := ValidateGetUsersRequest(&request)
    if err == nil {
        t.Errorf(`ValidateGetUserRequest(&request) = "%v", expected "invalid sort direction"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_Success(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.UserId != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }

    username := ValidUsername
    email := ValidEmail
    password := ValidPassword

    request := api.UpdateUserRequest{
        UserId:   1,
        Username: &username,
        Email:    &email,
        Password: &password,
    }

    if err := ValidateUpdateUserRequest(&request, service, nil); err != nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "<nil>"`, err)
    }
}

func TestValidateUpdateUserRequest_InvalidUserId(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.UserId != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }

    username := ValidUsername
    email := ValidEmail
    password := ValidPassword

    request := api.UpdateUserRequest{
        UserId:   -1,
        Username: &username,
        Email:    &email,
        Password: &password,
    }

    err := ValidateUpdateUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "invalid user id"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_UserNotFound(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return nil, nil
        },
    }

    username := ValidUsername
    email := ValidEmail
    password := ValidPassword

    request := api.UpdateUserRequest{
        UserId:   1,
        Username: &username,
        Email:    &email,
        Password: &password,
    }

    err := ValidateUpdateUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
    }
    assertHTTPError(t, err, http.StatusNotFound)
}

func TestValidateUpdateUserRequest_InvalidUsername(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.UserId != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }

    username := "%$#@!^&*()"
    request := api.UpdateUserRequest{
        UserId:   1,
        Username: &username,
    }

    err := ValidateUpdateUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_DuplicateUsername(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.UserId != nil {
                return &api.GetUserResponse{}, nil
            }
            if request.Username != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }

    username := ValidUsername
    request := api.UpdateUserRequest{
        UserId:   1,
        Username: &username,
    }

    err := ValidateUpdateUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_InvalidEmail(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.UserId != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }

    email := "uh"
    request := api.UpdateUserRequest{
        UserId: 1,
        Email:  &email,
    }

    err := ValidateUpdateUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_DuplicateEmail(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.UserId != nil {
                return &api.GetUserResponse{}, nil
            }
            if request.Email != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }

    email := ValidEmail
    request := api.UpdateUserRequest{
        UserId: 1,
        Email:  &email,
    }

    err := ValidateUpdateUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateUpdateUserRequest_InvalidPassword(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            if request.UserId != nil {
                return &api.GetUserResponse{}, nil
            }
            return nil, nil
        },
    }

    password := "invalidPassword"
    request := api.UpdateUserRequest{
        UserId:   1,
        Password: &password,
    }

    err := ValidateUpdateUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateUpdateUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateDeleteUserRequest_Success(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return &api.GetUserResponse{}, nil
        },
    }

    request := api.DeleteUserRequest{
        UserId: 1,
    }

    if err := ValidateDeleteUserRequest(&request, service, nil); err != nil {
        t.Errorf(`ValidateDeleteUserRequest(&request, service, nil) = "%v", expected "<nil>"`, err)
    }
}

func TestValidateDeleteUserRequest_InvalidUserId(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return &api.GetUserResponse{}, nil
        },
    }

    request := api.DeleteUserRequest{
        UserId: -1,
    }

    err := ValidateDeleteUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateDeleteUserRequest(&request, service, nil) = "%v", expected "invalid user id"`, err)
    }
    assertHTTPError(t, err, http.StatusBadRequest)
}

func TestValidateDeleteUserRequest_UserNotFound(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return nil, nil
        },
    }

    request := api.DeleteUserRequest{
        UserId: 1,
    }

    err := ValidateDeleteUserRequest(&request, service, nil)
    if err == nil {
        t.Errorf(`ValidateDeleteUserRequest(&request, service, nil) = "%v", expected "user not found"`, err)
    }
    assertHTTPError(t, err, http.StatusNotFound)
}

func TestValidateVerifyUserRequest_Success(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return &api.GetUserResponse{}, nil
        },
    }

    request := api.VerifyUserRequest{
        UserId: 1,
    }

    if err := ValidateVerifyUserRequest(&request, service, nil); err != nil {
        t.Errorf(`ValidateVerifyUserRequest(&request, service, nil) = "%v", expected "<nil>"`, err)
    }
}

func TestValidateVerifyUserRequest_InvalidUserId(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return &api.GetUserResponse{}, nil
        },
    }

    request := api.VerifyUserRequest{
        UserId: -1,
    }

    if err := ValidateVerifyUserRequest(&request, service, nil); err == nil {
        t.Error(`ValidateVerifyUserRequest(&request, service, nil) = "<nil>", expected non-nil"`)
    }
}

func TestValidateVerifyUserRequest_UserNotFound(t *testing.T) {
    service := &mockService{
        getUserFunc: func(context context.Context, request *api.GetUserRequest) (*api.GetUserResponse, error) {
            return &api.GetUserResponse{}, nil
        },
    }

    request := api.VerifyUserRequest{
        UserId: -1,
    }

    if err := ValidateVerifyUserRequest(&request, service, nil); err == nil {
        t.Error(`ValidateVerifyUserRequest(&request, service, nil) = "<nil>", expected non-nil"`)
    }
}
