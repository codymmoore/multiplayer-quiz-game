package user

import (
    "common"
    "context"
    "errors"
    "fmt"
    "github.com/go-chi/chi/v5"
    "os"
    "testing"
    "time"
    "user/db/generated"
    "user/dto"
)

const (
    BASE_URL_KEY = "BASE_URL"
    MockUrl      = "https://mock-url"
)

func TestService_CreateUser_Success(t *testing.T) {
    userId := 1
    mockQuerier := &mockQuerier{
        createUserFunc: func(context context.Context, arg db.CreateUserParams) (db.User, error) {
            return db.User{
                ID:           int32(userId),
                Username:     ValidUsername,
                Email:        ValidEmail,
                PasswordHash: ValidPassword,
                IsVerified:   true,
                CreatedAt:    time.Now(),
                UpdatedAt:    time.Now(),
            }, nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    request := dto.CreateUserRequest{
        Username: ValidUsername,
        Email:    ValidEmail,
        Password: ValidPassword,
    }
    response, err := service.CreateUser(nil, &request)
    if err != nil {
        t.Errorf(`service.CreateUser(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.CreateUser(nil, request) response = "<nil>", expected non-nil`)
        return
    }
    if response.UserId != userId {
        t.Errorf(`response.UserId = %v, expected %v`, response.UserId, userId)
    }
}

func TestService_CreateUser_QueryFailure(t *testing.T) {
    mockQuerier := &mockQuerier{
        createUserFunc: func(context context.Context, arg db.CreateUserParams) (db.User, error) {
            return db.User{}, errors.New("")
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    request := dto.CreateUserRequest{
        Username: ValidUsername,
        Email:    ValidEmail,
        Password: ValidPassword,
    }
    if _, err := service.CreateUser(nil, &request); err == nil {
        t.Error(`service.CreateUser(nil, request) error = "<nil>", expected non-nil`)
    }
}

func TestService_GetUser_Success(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: ValidPassword,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        getUserFunc: func(context context.Context, arg db.GetUserParams) (db.User, error) {
            return mockUser, nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    request := dto.GetUserRequest{
        UserId:   &userId,
        Username: &username,
        Email:    &email,
    }
    response, err := service.GetUser(nil, &request)
    if err != nil {
        t.Errorf(`service.GetUser(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.GetUser(nil, request) response = "<nil>", expected non-nil`)
        return
    }
    assertUserEqualToDB(t, response, &mockUser)
}

func TestService_GetUser_QueryFailure(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    mockQuerier := &mockQuerier{
        getUserFunc: func(context context.Context, arg db.GetUserParams) (db.User, error) {
            return db.User{}, errors.New("")
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    request := dto.GetUserRequest{
        UserId:   &userId,
        Username: &username,
        Email:    &email,
    }
    if _, err := service.GetUser(nil, &request); err == nil {
        t.Error(`service.GetUser(nil, request) error = "<nil>", expected non-nil`)
    }
}

func TestService_GetUsers_Success(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: ValidPassword,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        getUsersFunc: func(context context.Context, arg db.GetUsersParams) ([]db.User, error) {
            return []db.User{
                mockUser,
            }, nil
        },
        countUsersFunc: func(ctx context.Context) (int64, error) {
            return int64(10), nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    limit := 1
    offset := 1
    sortField := "CreatedAt"
    sortDirection := "ASC"
    request := dto.GetUsersRequest{
        Limit:         &limit,
        Offset:        &offset,
        SortField:     &sortField,
        SortDirection: &sortDirection,
    }

    os.Setenv(BASE_URL_KEY, MockUrl)

    rctx := chi.NewRouteContext()
    rctx.RoutePatterns = []string{"/users/{userID}"}
    ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

    response, err := service.GetUsers(ctx, &request)
    if err != nil {
        t.Errorf(`service.GetUsers(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.GetUsers(nil, request) response = "<nil>", expected non-nil`)
        return
    }
    assertUserEqualToDB(t, &response.Users[0], &mockUser)

    routeUrl, _ := common.GetRouteUrl(ctx)
    expectedPrevLink := fmt.Sprintf(
        "%s?limit=%d&offset=%d",
        routeUrl,
        limit,
        offset-limit,
    )
    if response.PrevLink == nil {
        t.Errorf(`response.PrevLink = "<nil>", expected "%v"`, expectedPrevLink)
    } else if *response.PrevLink != expectedPrevLink {
        t.Errorf(`response.PrevLink = "%v", expected "%v"`, *response.PrevLink, expectedPrevLink)
    }

    expectedNextLink := fmt.Sprintf(
        "%s?limit=%d&offset=%d",
        routeUrl,
        limit,
        offset+limit,
    )
    if response.NextLink == nil {
        t.Errorf(`response.NextLink = "<nil>", expected "%v"`, expectedNextLink)
    } else if *response.NextLink != expectedNextLink {
        t.Errorf(`response.NextLink = "%v", expected "%v"`, *response.NextLink, expectedNextLink)
    }
}

func TestService_GetUsers_DefaultLimit(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: ValidPassword,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        getUsersFunc: func(context context.Context, arg db.GetUsersParams) ([]db.User, error) {
            return []db.User{
                mockUser,
            }, nil
        },
        countUsersFunc: func(ctx context.Context) (int64, error) {
            return int64(DefaultUsersPageLimit + 1), nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    offset := 1
    request := dto.GetUsersRequest{
        Offset: &offset,
    }

    os.Setenv(BASE_URL_KEY, MockUrl)

    rctx := chi.NewRouteContext()
    rctx.RoutePatterns = []string{"/users/{userID}"}
    ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

    response, err := service.GetUsers(ctx, &request)
    if err != nil {
        t.Errorf(`service.GetUsers(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.GetUsers(nil, request) response = "<nil>", expected non-nil`)
        return
    }
    assertUserEqualToDB(t, &response.Users[0], &mockUser)

    routeUrl, _ := common.GetRouteUrl(ctx)
    expectedPrevLink := fmt.Sprintf(
        "%s?limit=%d&offset=%d",
        routeUrl,
        DefaultUsersPageLimit,
        0,
    )
    if response.PrevLink == nil {
        t.Errorf(`response.PrevLink = "<nil>", expected "%v"`, expectedPrevLink)
    } else if *response.PrevLink != expectedPrevLink {
        t.Errorf(`response.PrevLink = "%v", expected "%v"`, *response.PrevLink, expectedPrevLink)
    }

    if response.NextLink != nil {
        t.Errorf(`response.NextLink = "%v", expected "<nil>"`, *response.NextLink)
    }
}

func TestService_GetUsers_DefaultOffset(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: ValidPassword,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        getUsersFunc: func(context context.Context, arg db.GetUsersParams) ([]db.User, error) {
            return []db.User{
                mockUser,
            }, nil
        },
        countUsersFunc: func(ctx context.Context) (int64, error) {
            return int64(10), nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    limit := 1
    offset := 0
    request := dto.GetUsersRequest{
        Limit: &limit,
    }

    os.Setenv(BASE_URL_KEY, MockUrl)

    rctx := chi.NewRouteContext()
    rctx.RoutePatterns = []string{"/users/{userID}"}
    ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

    response, err := service.GetUsers(ctx, &request)
    if err != nil {
        t.Errorf(`service.GetUsers(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.GetUsers(nil, request) response = "<nil>", expected non-nil`)
        return
    }
    assertUserEqualToDB(t, &response.Users[0], &mockUser)

    routeUrl, _ := common.GetRouteUrl(ctx)
    if response.PrevLink != nil {
        t.Errorf(`response.PrevLink = "%v", expected "<nil>"`, *response.PrevLink)
    }

    expectedNextLink := fmt.Sprintf(
        "%s?limit=%d&offset=%d",
        routeUrl,
        limit,
        offset+limit,
    )
    if response.NextLink == nil {
        t.Errorf(`response.NextLink = "<nil>", expected "%v"`, expectedNextLink)
    } else if *response.NextLink != expectedNextLink {
        t.Errorf(`response.NextLink = "%v", expected "%v"`, *response.NextLink, expectedNextLink)
    }
}

func TestService_GetUsers_CountUsersFail(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: ValidPassword,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        getUsersFunc: func(context context.Context, arg db.GetUsersParams) ([]db.User, error) {
            return []db.User{
                mockUser,
            }, nil
        },
        countUsersFunc: func(ctx context.Context) (int64, error) {
            return 0, errors.New("")
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    sortField := "CreatedAt"
    sortDirection := "ASC"
    request := dto.GetUsersRequest{
        SortField:     &sortField,
        SortDirection: &sortDirection,
    }

    os.Setenv(BASE_URL_KEY, MockUrl)

    rctx := chi.NewRouteContext()
    rctx.RoutePatterns = []string{"/users/{userID}"}
    ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

    _, err := service.GetUsers(ctx, &request)
    if err == nil {
        t.Error(`service.GetUsers(nil, request) error = "<nil>", expected non-nil`)
    }
}

func TestService_GetUsers_QueryFailure(t *testing.T) {
    mockQuerier := &mockQuerier{
        getUsersFunc: func(context context.Context, arg db.GetUsersParams) ([]db.User, error) {
            return nil, errors.New("")
        },
        countUsersFunc: func(ctx context.Context) (int64, error) {
            return int64(10), nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    sortField := "CreatedAt"
    sortDirection := "ASC"
    request := dto.GetUsersRequest{
        SortField:     &sortField,
        SortDirection: &sortDirection,
    }

    os.Setenv(BASE_URL_KEY, MockUrl)

    rctx := chi.NewRouteContext()
    rctx.RoutePatterns = []string{"/users/{userID}"}
    ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

    _, err := service.GetUsers(ctx, &request)
    if err == nil {
        t.Error(`service.GetUsers(nil, request) error = "<nil>", expected non-nil`)
    }
}

func TestService_GetUsers_GetRouteUrlFailure(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: ValidPassword,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        getUsersFunc: func(context context.Context, arg db.GetUsersParams) ([]db.User, error) {
            return []db.User{
                mockUser,
            }, nil
        },
        countUsersFunc: func(ctx context.Context) (int64, error) {
            return int64(10), nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    sortField := "CreatedAt"
    sortDirection := "ASC"
    request := dto.GetUsersRequest{
        SortField:     &sortField,
        SortDirection: &sortDirection,
    }

    os.Unsetenv(BASE_URL_KEY)

    rctx := chi.NewRouteContext()
    rctx.RoutePatterns = []string{"/users/{userID}"}
    ctx := context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

    _, err := service.GetUsers(ctx, &request)
    if err == nil {
        t.Error(`service.GetUsers(nil, request) error = "<nil>", expected non-nil`)
    }
}

func TestService_UpdateUser_Success(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    password := ValidPassword
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: password,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        updateUserFunc: func(context context.Context, arg db.UpdateUserParams) (db.User, error) {
            return mockUser, nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    request := dto.UpdateUserRequest{
        UserId:   userId,
        Username: &username,
        Email:    &email,
        Password: &password,
    }
    response, err := service.UpdateUser(nil, &request)
    if err != nil {
        t.Errorf(`service.UpdateUser(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.GetUser(nil, request) response = "<nil>", expected non-nil`)
        return
    }
    assertUserEqualToDB(t, response, &mockUser)
}

func TestService_UpdateUser_NoChange(t *testing.T) {
    userId := 1
    username := ValidUsername
    email := ValidEmail
    password := ValidPassword
    mockUser := db.User{
        ID:           int32(userId),
        Username:     username,
        Email:        email,
        PasswordHash: password,
        IsVerified:   true,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    mockQuerier := &mockQuerier{
        updateUserFunc: func(context context.Context, arg db.UpdateUserParams) (db.User, error) {
            return mockUser, nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    request := dto.UpdateUserRequest{
        UserId: userId,
    }
    response, err := service.UpdateUser(nil, &request)
    if err != nil {
        t.Errorf(`service.UpdateUser(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.GetUser(nil, request) response = "<nil>", expected non-nil`)
        return
    }
    assertUserEqualToDB(t, response, &mockUser)
}

func TestService_UpdateUser_QueryFailure(t *testing.T) {
    userId := 1
    mockQuerier := &mockQuerier{
        updateUserFunc: func(context context.Context, arg db.UpdateUserParams) (db.User, error) {
            return db.User{}, errors.New("")
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    request := dto.UpdateUserRequest{
        UserId: userId,
    }
    _, err := service.UpdateUser(nil, &request)
    if err == nil {
        t.Error(`service.UpdateUser(nil, request) error = "<nil>", expected non-nil`)
    }
}

func TestService_DeleteUser_Success(t *testing.T) {
    mockQuerier := &mockQuerier{
        deleteUserFunc: func(context context.Context, arg int32) error {
            return nil
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    userId := 1
    request := dto.DeleteUserRequest{
        UserId: userId,
    }
    response, err := service.DeleteUser(nil, &request)
    if err != nil {
        t.Errorf(`service.DeleteUser(nil, request) error = "%v", expected "<nil>"`, err)
    }
    if response == nil {
        t.Error(`service.GetUser(nil, request) response = "<nil>", expected non-nil`)
    }
}

func TestService_DeleteUser_QueryFailure(t *testing.T) {
    mockQuerier := &mockQuerier{
        deleteUserFunc: func(context context.Context, arg int32) error {
            return errors.New("")
        },
    }
    service := ServiceImpl{
        Queries: mockQuerier,
    }

    userId := 1
    request := dto.DeleteUserRequest{
        UserId: userId,
    }
    _, err := service.DeleteUser(nil, &request)
    if err == nil {
        t.Error(`service.DeleteUser(nil, request) error = "<nil>", expected non-nil`)
    }
}

type mockQuerier struct {
    countUsersFunc func(ctx context.Context) (int64, error)
    createUserFunc func(ctx context.Context, arg db.CreateUserParams) (db.User, error)
    deleteUserFunc func(ctx context.Context, id int32) error
    getUserFunc    func(ctx context.Context, arg db.GetUserParams) (db.User, error)
    getUsersFunc   func(ctx context.Context, arg db.GetUsersParams) ([]db.User, error)
    updateUserFunc func(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
}

func (q *mockQuerier) CountUsers(ctx context.Context) (int64, error) {
    return q.countUsersFunc(ctx)
}

func (q *mockQuerier) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
    return q.createUserFunc(ctx, arg)
}

func (q *mockQuerier) DeleteUser(ctx context.Context, id int32) error {
    return q.deleteUserFunc(ctx, id)
}

func (q *mockQuerier) GetUser(ctx context.Context, arg db.GetUserParams) (db.User, error) {
    return q.getUserFunc(ctx, arg)
}

func (q *mockQuerier) GetUsers(ctx context.Context, arg db.GetUsersParams) ([]db.User, error) {
    return q.getUsersFunc(ctx, arg)
}

func (q *mockQuerier) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
    return q.updateUserFunc(ctx, arg)
}

func assertUserEqualToDB(t *testing.T, actual *dto.User, expected *db.User) {
    if actual.UserId != int(expected.ID) {
        t.Errorf(`actual.UserId = "%d", expected "%d"`, actual.UserId, expected.ID)
    }

    if actual.Username != expected.Username {
        t.Errorf(`actual.Username = "%s", expected "%s"`, actual.Username, expected.Username)
    }

    if actual.Email != expected.Email {
        t.Errorf(`actual.Email = "%s", expected "%s"`, actual.Email, expected.Email)
    }

    if actual.PasswordHash != expected.PasswordHash {
        t.Errorf(`actual.PasswordHash = "%s", expected "%s"`, actual.PasswordHash, expected.PasswordHash)
    }

    if actual.IsVerified != true {
        t.Errorf(`actual.IsVerified = "%t"`, actual.IsVerified)
    }

    if !actual.CreatedAt.Equal(expected.CreatedAt) {
        t.Errorf(`actual.CreatedAt = "%s", expected "%s"`, actual.CreatedAt, expected.CreatedAt)
    }

    if !actual.UpdatedAt.Equal(expected.UpdatedAt) {
        t.Errorf(`actual.UpdatedAt = "%s", expected "%s"`, actual.UpdatedAt, expected.UpdatedAt)
    }
}
