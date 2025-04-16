package converters

import (
	"poc-fiber/dtos"
	"poc-fiber/model"
)

func ConvertUserToResponse(user model.User) dtos.UserResponse {
	var usrResponse = dtos.UserResponse{
		Uuid:      &user.Uuid,
		LastName:  &user.LastName,
		FirstName: &user.FirstName,
		Login:     &user.Login,
		Email:     &user.Email,
	}
	return usrResponse
}
