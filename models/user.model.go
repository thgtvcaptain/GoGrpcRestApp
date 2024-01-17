package models

type User struct {
	Id       uint32 `json:"id" bson:"id" binding:"required"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Age      uint32 `json:"age" bson:"age" binding:"required"`
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Id       uint32 `json:"id" bson:"id" binding:"required"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Age      uint32 `json:"age,omitempty" bson:"age,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username" bson:"username" binding:"required"`
	OldPassword string `json:"oldpassword" bson:"oldpassword" binding:"required"`
	NewPassword string `json:"newpassword" bson:"newpassword" binding:"required"`
}

type ChangePasswordResponse struct {
	Username string `json:"username" bson:"username" binding:"required"`
}

type UserResponse struct {
	Id       uint32 `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Age      uint32 `json:"age" bson:"age"`
	Username string `json:"username" bson:"username"`
}
