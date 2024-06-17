package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/dibyendu/Authentication-Authorization/lib/errs"
	"github.com/dibyendu/Authentication-Authorization/pkg/dto"
	"github.com/dibyendu/Authentication-Authorization/pkg/middleware"
	"github.com/dibyendu/Authentication-Authorization/pkg/service"
)

type UserHandler struct {
	Service service.UserService
}

// CreateUser handles the HTTP POST request to create a new user.
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User details"
// @Success 200 {object} dto.CreateUserResponse "User created successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Router /create-user [post]
func (h UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		request dto.CreateUserRequest
		ctx     = r.Context()
	)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		response, appError := h.Service.CreateUser(ctx, request)
		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, response)
		}
	}
}
// SignIn handles the HTTP POST request to sign in existing user.
// @Summary signs in user
// @Description sign in a user with the provided details
// @Tags SignInUser
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "login details"
// @Success 200 {object} dto.UserLoginResponse "User login successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Router /user/sign-in [post]
func (h UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		request dto.CreateUserRequest
		ctx     = r.Context()
	)
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		response, appError := h.Service.SignIn(ctx, request)
		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, response)
		}
	}
}
// GetUser handles the HTTP GET request to fetch user details.
// @Security ApiKeyAuth
// @Summary Fetch user details
// @Description Fetch details of a user by their ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id query string true "User ID"
// @Security bearerToken
// @Success 200 {object} dto.GetUserResponse "User details fetched successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Unauthorized access"
// @Failure 404 {object} string "User not found"
// @Router /user/get-user [get]
func (h UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var (
		request dto.GetUserRequest
		ctx     = r.Context()
	)
	if !stringUtils.IsBlank(r.URL.Query().Get("id")) {
		request.Id = r.URL.Query().Get("id")
	}

	userInfo := middleware.GetUserInfo(ctx)
	if strings.EqualFold(strings.ToLower(userInfo.AuthToken), "invalid") {
		writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("unauthorized Access for get user").AsMessage())
		return
	}
	if !strings.EqualFold(strings.ToLower(userInfo.Role), "admin") {
		writeResponse(w, http.StatusUnauthorized, errs.NewValidationError("wrong user it's only for admin").AsMessage())
		return
	}
	response, err := h.Service.GetUser(ctx, request)
	if err != nil {
		if err.Code == http.StatusNoContent {
			writeResponseNoContent(w, err.Code)
		} else {
			writeResponse(w, err.Code, err.AsMessage())
		}
	} else {
		writeResponse(w, http.StatusOK, response)
	}
}
