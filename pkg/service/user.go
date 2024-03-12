package service

import (
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/utility"
	"github.com/dibyendu/Authentication-Authorization/pkg/domain"
	"github.com/dibyendu/Authentication-Authorization/pkg/dto"
	"github.com/dibyendu/Authentication-Authorization/pkg/middleware"
	"context"
)

type DefaultUserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) DefaultUserService {
	return DefaultUserService{
		repo: repo,
	}
}

type UserService interface {
	CreateUser(ctx context.Context, request dto.CreateUserRequest) (*dto.CreateUserResponse, *errs.AppError)
	SignIn(ctx context.Context, request dto.CreateUserRequest) (*dto.UserLoginResponse, *errs.AppError)
	GetUser(ctx context.Context, request dto.GetUserRequest) (*dto.GetUserResponse, *errs.AppError)
}

func (s DefaultUserService) CreateUser(ctx context.Context, request dto.CreateUserRequest) (*dto.CreateUserResponse, *errs.AppError) {
	var (
		response    *dto.CreateUserResponse
		validateErr *errs.AppError
	)
	//useAccess := middleware.GetUserInfo(ctx)
	//if ok := strings.EqualFold(strings.ToLower(useAccess.Role),"admin"); !ok {
	//	return nil, errs.NewValidationError("userr cannot be empty")
	//}
	validateErr = request.Validate()
	if validateErr != nil {
		return nil, validateErr
	}
	req := domain.CreateUserRequest(request)
	data, err := s.repo.CreateUser(ctx, req)

	if err != nil {
		return nil, err
	}
	response = data.ToDto()
	return response, nil
}

func (s DefaultUserService) SignIn(ctx context.Context, request dto.CreateUserRequest) (*dto.UserLoginResponse, *errs.AppError) {
	var (
		validateErr *errs.AppError
	)
	validateErr = request.Validate()
	if validateErr != nil {
		return nil, validateErr
	}
	req := domain.CreateUserRequest(request)
	data, err1 := s.repo.IsEmailExists(ctx, req.Email)
	if err1 != nil {
		return nil, err1
	}
	//matching the password
	err := utility.VerifyPassword(data.Password, request.Password)
	if err != nil {
		return nil, errs.NewValidationError("password is not valid")
	}
	token, err := middleware.GenerateJWT(data.Email, data.Role)
	if err != nil {
		return nil, errs.NewValidationError("token is not valid")
	}
	return &dto.UserLoginResponse{
		Id:    data.Id,
		Name:  data.Name,
		Role:  data.Role,
		Email: data.Email,
		Token: token,
	}, nil
}

func(s DefaultUserService) GetUser(ctx context.Context, request dto.GetUserRequest)(*dto.GetUserResponse, *errs.AppError){
	var(
		validateErr *errs.AppError
	)

	validateErr = request.Validate()
	if validateErr != nil {
		return nil, validateErr
	}
	req := domain.GetUserRequest(request)
	data, err := s.repo.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	response := data.ToDto()
	return &response, nil
}