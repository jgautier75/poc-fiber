package dtos

type CreateUserRequest struct {
	LastName  *string `json:"last_name"`
	FirstName *string `json:"first_name"`
	Login     *string `json:"login"`
	Email     *string `json:"email"`
}
