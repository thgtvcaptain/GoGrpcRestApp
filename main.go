package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"protoUserManagement/controller"
	"protoUserManagement/gapi"
	proto "protoUserManagement/pb"
	"protoUserManagement/route"
	"protoUserManagement/services"
)

var mongoclient *mongo.Client

func init() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Not found .env file")
	}

	if err := connect_to_mongodb(); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	userCollection := mongoclient.Database("Company").Collection("User")
	userServiceImpl := services.NewUserService(userCollection, context.TODO())
	userController, _ := controller.NewUserController(userServiceImpl)

	router := gin.Default()

	route.AddUserRoute(router, userController)

	go router.Run()

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterUserServiceServer(grpcServer,
		&gapi.Server{
			UserCollection: userCollection,
			UserService:    userServiceImpl,
		})
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)

}

func connect_to_mongodb() error {

	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	clusterurl := os.Getenv("MONGODB_URL")
	connectionType := os.Getenv("CONNECTION_TYPE")

	uri := connectionType + "://" + username + ":" + password + "@" + clusterurl + "/?retryWrites=true&w=majority"

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	mongoclient = client
	return err
}
