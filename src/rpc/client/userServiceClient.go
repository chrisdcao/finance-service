package client

import (
	"context"
	"errors"
	pb "finance-service/rpc/protos/user/generated_files" // Adjust the import path
	"finance-service/utils"
)

type UserServiceClient struct {
	client pb.UserServiceClient
}

func NewUserServiceClient(client pb.UserServiceClient) *UserServiceClient {
	return &UserServiceClient{
		client: client,
	}
}

func (c *UserServiceClient) GetUser(uuid string) (*pb.User, error) {
	req := &pb.UserRequest{Uuid: uuid}
	res, err := c.client.GetUser(context.Background(), req)
	if err != nil {
		utils.Logger().Println("Error getting user from user service:", err)
		return nil, errors.New("user is not registered")
	}
	return res.User, nil
}
