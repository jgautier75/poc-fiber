package dtos

type CreateUserRequest struct {
	LastName  *string `json:"last_name" validate:"required,max=50"`
	FirstName *string `json:"first_name" validate:"required,max=50"`
	Login     *string `json:"login" validate:"required,max=50"`
	Email     *string `json:"email" validate:"required,max=50"`
}

type UserResponse struct {
	Uuid      *string `json:"uuid"`
	LastName  *string `json:"last_name"`
	FirstName *string `json:"first_name"`
	Login     *string `json:"login"`
	Email     *string `json:"email"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}
