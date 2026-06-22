package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/findardi/Wadi/server/internal/access/dto"
	"github.com/findardi/Wadi/server/internal/access/service"
	"github.com/findardi/Wadi/server/internal/platform/middleware"
	"github.com/findardi/Wadi/server/internal/platform/permission"
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

func (h *AccessHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "get permissions success", permission.All)
}

func (h *AccessHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	wID := chi.URLParam(r, "workspaceID")

	var req dto.CreateWorkspaceMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceId = wID

	res, err := h.svc.AddMember(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMemberAlreadyAdd):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusCreated, "add member success", res)
}

func (h *AccessHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	wID := chi.URLParam(r, "workspaceID")

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.AddMembersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceId = wID
	req.InvitedBy = claims.ID

	res, err := h.svc.AddMembers(r.Context(), req)
	if err != nil {
		log.Printf("add members internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "invite members processed", res)
}

func (h *AccessHandler) GetInvitations(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")

	res, err := h.svc.ListInvitations(r.Context(), wID)
	if err != nil {
		log.Printf("list invitations internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "get invitations success", res)
}

func (h *AccessHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")

	res, err := h.svc.GetMembers(r.Context(), wID)
	if err != nil {
		log.Printf("register internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "get members success", res)
}

func (h *AccessHandler) GetMember(w http.ResponseWriter, r *http.Request) {
	mID := chi.URLParam(r, "memberID")

	res, err := h.svc.GetMember(r.Context(), mID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrMemberNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "get member success", res)
}

func (h *AccessHandler) UpdateMember(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	mID := chi.URLParam(r, "memberID")

	var req dto.UpdateMemberRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.MemberID = mID

	res, err := h.svc.UpdateMemberRole(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrMemberNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "update member success", res)
}

func (h *AccessHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	mID := chi.URLParam(r, "memberID")

	if err := h.svc.DeleteMember(r.Context(), mID); err != nil {
		log.Printf("register internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "delete member success", nil)
}
