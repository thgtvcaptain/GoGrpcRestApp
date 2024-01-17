package services

import "protoUserManagement/models"

type UserService interface {
	Login(request *models.LoginRequest) (*models.UserResponse, error)
	ChangePassword(request *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error)
	CreateUser(user *models.User) (*models.UserResponse, error)
	UpdateUser(request *models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(userId uint32) (uint32, error)
	GetUser(userId uint32) (*models.UserResponse, error)
	GetAllUsers() ([]*models.UserResponse, error)
}
