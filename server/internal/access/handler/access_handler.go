package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/findardi/Wadi/server/internal/access/dto"
	"github.com/findardi/Wadi/server/internal/access/service"
	"github.com/findardi/Wadi/server/internal/platform/response"
	"github.com/findardi/Wadi/server/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

const (
	MaxBodyBytes = 1 << 20
)

type AccessHandler struct {
	svc *service.AccessService
}

func NewAccessHandler(svc *service.AccessService) *AccessHandler {
	return &AccessHandler{
		svc: svc,
	}
}

func (h *AccessHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	wID := chi.URLParam(r, "workspaceID")

	var req dto.CreateWorkspaceRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = wID

	res, err := h.svc.InsertRole(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRoleNameTaken):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "create role success", res)
}

func (h *AccessHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")

	res, err := h.svc.GetRoles(r.Context(), wID)
	if err != nil {
		log.Printf("register internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "get roles success", res)
}

func (h *AccessHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	rID := chi.URLParam(r, "roleID")

	res, err := h.svc.GetRole(r.Context(), rID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRoleNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "get role success", res)
}

func (h *AccessHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	wID := chi.URLParam(r, "workspaceID")
	rID := chi.URLParam(r, "roleID")

	var req dto.UpdateWorkspaceRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = wID
	req.RoleID = rID

	res, err := h.svc.UpdateRole(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRoleNameTaken):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, service.ErrRoleNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "update role success", res)
}

func (h *AccessHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	rID := chi.URLParam(r, "roleID")

	if err := h.svc.DeleteRole(r.Context(), rID); err != nil {
		switch {
		case errors.Is(err, service.ErrRoleInUse):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "delete role success", nil)
}
