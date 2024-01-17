package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"protoUserManagement/models"
	"protoUserManagement/services"
	"strconv"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) (*UserController, error) {
	return &UserController{userService}, nil
}

func (uc *UserController) Login(ctx *gin.Context) {

	var loginRequest *models.LoginRequest
	if err := ctx.BindJSON(&loginRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse, err := uc.UserService.Login(loginRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, userResponse)
}

func (uc *UserController) ChangePassword(ctx *gin.Context) {

	var changePassRequest *models.ChangePasswordRequest
	if err := ctx.BindJSON(&changePassRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	changePassRes, err := uc.UserService.ChangePassword(changePassRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, changePassRes)
}

func (uc *UserController) GetAllUsers(ctx *gin.Context) {

	userResponseList, err := uc.UserService.GetAllUsers()
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, userResponseList)
}

func (uc *UserController) GetUserById(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userResponse, err := uc.UserService.GetUser(uint32(id))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, userResponse)
}

func (uc *UserController) AddUser(ctx *gin.Context) {

	var user *models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse, err := uc.UserService.CreateUser(user)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, userResponse)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {

	var updateRequest *models.UpdateUserRequest
	if err := ctx.BindJSON(&updateRequest); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResponse, err := uc.UserService.UpdateUser(updateRequest)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, userResponse)
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deletedCount, err := uc.UserService.DeleteUser(uint32(id))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, map[string]any{"deletedCount": deletedCount})

}
