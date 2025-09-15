package user

import (
	"common/test"
	"fmt"
	"strings"
	"testing"
)

func TestClient_createGetUserQueryString(t *testing.T) {
	userId := 1
	username := test.ValidUsername
	email := test.ValidEmail
	expectedQueryString := fmt.Sprintf(
		"%s=%s&%s=%d&%s=%s",
		EmailKey,
		strings.ReplaceAll(email, "@", "%40"),
		UserIdKey,
		userId,
		UsernameKey,
		username,
	)

	request := &GetUserRequest{
		UserId:   &userId,
		Username: &username,
		Email:    &email,
	}
	queryString, err := createGetUserQueryString(request)
	if err != nil {
		t.Errorf("failed to create query string: %v", err)
	}
	if queryString != expectedQueryString {
		t.Errorf(`queryString: "%s", expected "%s"`, queryString, expectedQueryString)
	}
}

func TestClient_createGetUsersQueryString(t *testing.T) {
	limit := 1
	offset := 2
	sortField := "CreatedAt"
	sortDirection := "asc"
	expectedQueryString := fmt.Sprintf(
		"%s=%d&%s=%d&%s=%s&%s=%s",
		LimitKey,
		limit,
		OffsetKey,
		offset,
		SortDirectionKey,
		sortDirection,
		SortFieldKey,
		sortField,
	)

	request := &GetUsersRequest{
		Limit:         &limit,
		Offset:        &offset,
		SortField:     &sortField,
		SortDirection: &sortDirection,
	}
	queryString, err := createGetUsersQueryString(request)
	if err != nil {
		t.Errorf("failed to create query string: %v", err)
	}
	if queryString != expectedQueryString {
		t.Errorf(`queryString: "%s", expected "%s"`, queryString, expectedQueryString)
	}
}
