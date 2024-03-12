package dto

import (
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/lib/utility"
	"github.com/agrison/go-commons-lang/stringUtils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserRequest struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	Id primitive.ObjectID `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Email string `json:"email"`
	Password string `json:"password,omitempty"`
}

type UserLoginResponse struct {
	Id primitive.ObjectID `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func(c CreateUserRequest) Validate() *errs.AppError  {
	if stringUtils.IsBlank(c.Email) {
		return errs.NewValidationError("email cannot be empty")
	}
	if stringUtils.IsBlank(c.Password) {
		return errs.NewValidationError("password cannot be empty")
	} else {
		if ok := utility.IsStrongPassword(c.Password); !ok{
			return errs.NewValidationError("password should contain at least one uppercase character,one digit, one special character and length must be 8 above")
		}
	}
	return nil
}

type GetUserRequest struct{
	Id string `json:"id"`
}
type GetUserResponse struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Email string `json:"email"`
}

func(req GetUserRequest) Validate() *errs.AppError{
	if stringUtils.IsBlank(req.Id){
		return errs.NewValidationError("id can't blank")
	}
	return nil
}