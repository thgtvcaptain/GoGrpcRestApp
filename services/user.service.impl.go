package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"protoUserManagement/models"
	"strings"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) *UserServiceImpl {
	return &UserServiceImpl{userCollection, ctx}
}

func (us *UserServiceImpl) Login(request *models.LoginRequest) (*models.UserResponse, error) {
	filter, err := bson.Marshal(request)

	if err != nil {
		return nil, err
	}

	result := us.userCollection.FindOne(us.ctx, filter)

	var user *models.UserResponse
	if err = result.Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) ChangePassword(request *models.ChangePasswordRequest) (*models.ChangePasswordResponse, error) {

	filter := bson.M{"username": request.Username}

	var user models.User

	err := us.userCollection.FindOne(us.ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	update := bson.M{"username": request.Username, "password": request.NewPassword}

	if strings.Compare(user.Password, request.OldPassword) == 0 {
		_, err := us.userCollection.UpdateOne(us.ctx, filter, bson.M{"$set": update})
		if err != nil {
			return nil, err
		}
	}

	return &models.ChangePasswordResponse{
		Username: request.Username,
	}, nil
}

func (us *UserServiceImpl) CreateUser(user *models.User) (*models.UserResponse, error) {
	result, err := us.userCollection.InsertOne(us.ctx, user)
	if err != nil {
		return nil, err
	}

	var response *models.UserResponse
	query := bson.M{"_id": result.InsertedID}
	if err = us.userCollection.FindOne(us.ctx, query).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func (us *UserServiceImpl) UpdateUser(request *models.UpdateUserRequest) (*models.UserResponse, error) {
	filter := bson.M{"id": request.Id}
	returnUpdated := options.After

	options := &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnUpdated,
	}
	var user *models.UserResponse

	if err := us.userCollection.FindOneAndUpdate(us.ctx, filter, bson.M{"$set": request}, options).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) DeleteUser(userId uint32) (uint32, error) {

	filter := bson.M{"id": userId}

	result, err := us.userCollection.DeleteOne(us.ctx, filter)
	if err != nil {
		return 0, err
	}

	return uint32(result.DeletedCount), nil
}

func (us *UserServiceImpl) GetUser(userId uint32) (*models.UserResponse, error) {
	query := bson.M{"id": userId}

	var user *models.UserResponse

	if err := us.userCollection.FindOne(us.ctx, query).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) GetAllUsers() ([]*models.UserResponse, error) {
	query := bson.M{}

	cursor, _ := us.userCollection.Find(us.ctx, query)
	defer cursor.Close(us.ctx)

	var userResponseList []*models.UserResponse

	if err := cursor.All(us.ctx, &userResponseList); err != nil {
		return nil, err
	}

	return userResponseList, nil
}
