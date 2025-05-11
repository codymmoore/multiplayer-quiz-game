package dto

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetUserRequest struct {
	UserId   *int    `json:"userId"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type GetUsersRequest struct {
	Limit     *int    `json:"limit"`
	Offset    *int    `json:"offset"`
	SortField *string `json:"sortField"`
	SortOrder *string `json:"sortOrder"`
}

type UpdateUserRequest struct {
	UserId   int     `json:"userId"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type DeleteUserRequest struct {
	UserId int `json:"userId"`
}
