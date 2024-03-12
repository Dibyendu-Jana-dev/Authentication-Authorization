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
	if !strings.EqualFold(strings.ToLower(userInfo.Role), "user") {
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
