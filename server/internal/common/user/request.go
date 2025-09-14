package user

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
	Limit         *int    `json:"limit"`
	Offset        *int    `json:"offset"`
	SortField     *string `json:"sortField"`
	SortDirection *string `json:"sortDirection"`
}

type UpdateUserRequest struct {
	UserId   int
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type VerifyUserRequest struct {
	UserId int
}

type DeleteUserRequest struct {
	UserId int `json:"userId"`
}
