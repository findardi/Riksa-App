package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/findardi/Wadi/server/internal/auth/dto"
	"github.com/findardi/Wadi/server/internal/auth/service"
	"github.com/findardi/Wadi/server/internal/platform/response"
	"github.com/findardi/Wadi/server/internal/platform/validation"
)

const (
	MaxBodyBytes = 1 << 20
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	res, err := h.svc.RegisterUser(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailUnique), errors.Is(err, service.ErrUsernameUnique):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusCreated, "success registered account", res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	res, err := h.svc.LoginUser(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Error(w, http.StatusUnauthorized, err.Error(), nil)
		default:
			log.Printf("login internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "login success", res)
}
