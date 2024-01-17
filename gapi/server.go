package gapi

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"protoUserManagement/models"
	"protoUserManagement/pb"
	"protoUserManagement/services"
)

type Server struct {
	proto.UnimplementedUserServiceServer
	UserCollection *mongo.Collection
	UserService    services.UserService
}

func NewGrpcUserServer(userCollection *mongo.Collection, userService services.UserService) (*Server, error) {
	userServer := &Server{
		UserCollection: userCollection,
		UserService:    userService,
	}

	return userServer, nil
}

func (s *Server) LoginUser(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {

	data := models.LoginRequest{
		Username: request.GetUsername(),
		Password: request.GetPassword(),
	}

	var user *models.UserResponse

	user, err := s.UserService.Login(&data)
	if err != nil {

		return &proto.LoginResponse{
			Response: &proto.LoginResponse_ErrorResponse{
				ErrorResponse: &proto.ErrorResponse{
					ErrorCode: uint32(codes.Unauthenticated),
					Message:   err.Error(),
				},
			},
		}, nil
	}

	protoUser := &proto.LoginResponse{
		Response: &proto.LoginResponse_UserResponse{
			UserResponse: &proto.UserResponse{
				Id:       user.Id,
				Name:     user.Name,
				Age:      user.Age,
				Username: user.Username,
			},
		},
	}

	return protoUser, nil
}

func (s *Server) ChangePassword(ctx context.Context, request *proto.ChangePasswordRequest) (*proto.ChangePasswordResponse, error) {

	data := models.ChangePasswordRequest{
		Username:    request.GetUsername(),
		OldPassword: request.GetOldPassword(),
		NewPassword: request.GetNewPassword(),
	}

	_, err := s.UserService.ChangePassword(&data)
	if err != nil {
		return &proto.ChangePasswordResponse{
			Response: &proto.ChangePasswordResponse_ErrorResponse{
				ErrorResponse: &proto.ErrorResponse{
					ErrorCode: uint32(codes.Unauthenticated),
					Message:   err.Error(),
				},
			},
		}, nil

	}

	protoResponse := &proto.ChangePasswordResponse{
		Response: &proto.ChangePasswordResponse_ChangePasswordSuccessResponse{
			ChangePasswordSuccessResponse: &proto.ChangePasswordSuccessResponse{
				Username: request.GetUsername(),
			},
		},
	}

	return protoResponse, nil
}

func (s *Server) GetAllUsers(noparam *proto.NoParam, stream proto.UserService_GetAllUsersServer) error {

	var userList []*models.UserResponse

	userList, err := s.UserService.GetAllUsers()
	if err != nil {
		return status.Errorf(codes.Unknown, err.Error())
	}

	for _, user := range userList {
		err := stream.Send(&proto.GetUserResponse{
			Response: &proto.GetUserResponse_UserResponse{
				UserResponse: &proto.UserResponse{
					Id:       user.Id,
					Name:     user.Name,
					Age:      user.Age,
					Username: user.Username,
				},
			},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) ListUsers(stream proto.UserService_ListUsersServer) error {

	var user *models.UserResponse

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			stream.Send(&proto.GetUserResponse{
				Response: &proto.GetUserResponse_ErrorResponse{
					ErrorResponse: &proto.ErrorResponse{
						ErrorCode: uint32(codes.Unknown),
						Message:   err.Error(),
					},
				},
			})
			continue
		}

		user, err = s.UserService.GetUser(req.GetId())
		if err != nil {
			stream.Send(&proto.GetUserResponse{
				Response: &proto.GetUserResponse_ErrorResponse{
					ErrorResponse: &proto.ErrorResponse{
						ErrorCode: uint32(codes.Internal),
						Message:   err.Error(),
					},
				},
			})
			continue
		}

		if err = stream.Send(&proto.GetUserResponse{
			Response: &proto.GetUserResponse_UserResponse{
				UserResponse: &proto.UserResponse{
					Id:       user.Id,
					Name:     user.Name,
					Age:      user.Age,
					Username: user.Username,
				},
			},
		}); err != nil {
			stream.Send(&proto.GetUserResponse{
				Response: &proto.GetUserResponse_ErrorResponse{
					ErrorResponse: &proto.ErrorResponse{
						ErrorCode: uint32(codes.Unavailable),
						Message:   err.Error(),
					},
				},
			})
			continue
		}
	}
}

func (s *Server) GetUser(ctx context.Context, request *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	var user *models.UserResponse

	user, err := s.UserService.GetUser(request.GetId())
	if err != nil {
		return &proto.GetUserResponse{
			Response: &proto.GetUserResponse_ErrorResponse{
				ErrorResponse: &proto.ErrorResponse{
					ErrorCode: uint32(codes.Unavailable),
					Message:   err.Error(),
				},
			},
		}, nil
	}

	return &proto.GetUserResponse{
		Response: &proto.GetUserResponse_UserResponse{
			UserResponse: &proto.UserResponse{
				Id:       user.Id,
				Name:     user.Name,
				Age:      user.Age,
				Username: user.Username,
			},
		},
	}, nil
}

func (s *Server) AddUser(ctx context.Context, request *proto.AddUserRequest) (*proto.GetUserResponse, error) {

	data := &models.User{
		Id:       request.GetId(),
		Name:     request.GetName(),
		Age:      request.GetAge(),
		Username: request.GetUsername(),
		Password: request.GetPassword(),
	}

	user, err := s.UserService.CreateUser(data)
	if err != nil {
		return &proto.GetUserResponse{
			Response: &proto.GetUserResponse_ErrorResponse{
				ErrorResponse: &proto.ErrorResponse{
					ErrorCode: uint32(codes.Unavailable),
					Message:   err.Error(),
				},
			},
		}, err
	}

	return &proto.GetUserResponse{
		Response: &proto.GetUserResponse_UserResponse{
			UserResponse: &proto.UserResponse{
				Id:       user.Id,
				Name:     user.Name,
				Age:      user.Age,
				Username: user.Username,
			},
		},
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, request *proto.UpdateUserRequest) (*proto.GetUserResponse, error) {

	data := &models.UpdateUserRequest{
		Id:       request.GetId(),
		Name:     request.GetName(),
		Age:      request.GetAge(),
		Username: request.GetUsername(),
	}

	user, err := s.UserService.UpdateUser(data)
	if err != nil {
		return &proto.GetUserResponse{
			Response: &proto.GetUserResponse_ErrorResponse{
				ErrorResponse: &proto.ErrorResponse{
					ErrorCode: uint32(codes.Unavailable),
					Message:   err.Error(),
				},
			},
		}, err
	}

	return &proto.GetUserResponse{
		Response: &proto.GetUserResponse_UserResponse{
			UserResponse: &proto.UserResponse{
				Id:       user.Id,
				Name:     user.Name,
				Age:      user.Age,
				Username: user.Username,
			},
		},
	}, nil
}

func (s *Server) DeleteUser(stream proto.UserService_DeleteUserServer) error {

	var DeletedCount uint32
	DeletedIdList := make([]uint32, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.DeleteUserResponse{
				Response: &proto.DeleteUserResponse_DeleteUserSuccessResponse{
					DeleteUserSuccessResponse: &proto.DeleteUserSuccessResponse{
						Id:           DeletedIdList,
						DeletedCount: DeletedCount,
					},
				},
			})
		}
		if err != nil {
			return stream.SendAndClose(&proto.DeleteUserResponse{
				Response: &proto.DeleteUserResponse_ErrorResponse{
					ErrorResponse: &proto.ErrorResponse{
						ErrorCode: uint32(codes.Unknown),
						Message:   err.Error(),
					},
				},
			})
		}

		deleteSuccessCount, err := s.UserService.DeleteUser(req.GetId())
		if err != nil {
			return stream.SendAndClose(&proto.DeleteUserResponse{
				Response: &proto.DeleteUserResponse_ErrorResponse{
					ErrorResponse: &proto.ErrorResponse{
						ErrorCode: uint32(codes.Unknown),
						Message:   err.Error(),
					},
				},
			})
		}
		if deleteSuccessCount != 0 {
			DeletedIdList = append(DeletedIdList, req.GetId())
			DeletedCount++
		}
	}
}
