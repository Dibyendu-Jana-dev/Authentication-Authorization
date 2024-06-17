package domain

import (
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/pkg/dto"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	CreateUser(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, *errs.AppError)
	SignIn(ctx context.Context, request CreateUserRequest) (*CreateUserResponse, *errs.AppError)
	IsEmailExists(ctx context.Context, email string) (*CreateUserResponse, *errs.AppError)
	GetUser(ctx context.Context, req GetUserRequest) (*GetUserResponse, *errs.AppError)
}

type CreateUserRequest struct {
	Name     string `bson:"name"`
	Role     string `bson:"role"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type CreateUserResponse struct {
	Id       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Role     string             `bson:"role"`
	Email    string             `bson:"email"`
	Password string             `bson:"password,omitempty"`
}

func (r CreateUserResponse) ToDto() *dto.CreateUserResponse{
	return &dto.CreateUserResponse{
		Id:       r.Id,
		Name:     r.Name,
		Role:     r.Role,
		Email:    r.Email,
		Password: r.Password,
	}
}

type GetUserRequest struct{
	Id string `json:"id"`
}

type GetUserResponse struct{
	Id primitive.ObjectID `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Role string `json:"role" bson:"role"`
	Email string `json:"email" bson:"email"`
}

func (r GetUserResponse)ToDto() dto.GetUserResponse {
	return dto.GetUserResponse{
		Id: r.Id.Hex(),
		Name: r.Name,
		Role: r.Role,
		Email: r.Email,
	}
}