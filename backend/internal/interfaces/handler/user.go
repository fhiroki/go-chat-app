package handler

import "github.com/fhiroki/chat/internal/domain/user"

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
