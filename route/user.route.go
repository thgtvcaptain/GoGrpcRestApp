package route

import (
	"github.com/gin-gonic/gin"
	"protoUserManagement/controller"
)

func AddUserRoute(router *gin.Engine, userController *controller.UserController) {

	router.GET("/users", userController.GetAllUsers)
	router.GET("/users/:id", userController.GetUserById)
	router.POST("/users", userController.AddUser)
	router.PATCH("/users", userController.UpdateUser)
	router.DELETE("/users/:id", userController.DeleteUser)

	router.POST("/login", userController.Login)
	router.PATCH("/changepassword", userController.ChangePassword)
}
