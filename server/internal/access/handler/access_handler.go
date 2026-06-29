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

func (h *AccessHandler) GetMyAccess(w http.ResponseWriter, r *http.Request) {
	ms, ok := middleware.MembershipFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "get my access success", ms)
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
	status := r.URL.Query().Get("status")

	res, err := h.svc.ListInvitations(r.Context(), wID, status)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInvitationStatus):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("list invitations internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
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

func (h *AccessHandler) ResendInvitation(w http.ResponseWriter, r *http.Request) {
	invID := chi.URLParam(r, "invitationID")

	if err := h.svc.ResendInvitation(r.Context(), invID); err != nil {
		switch {
		case errors.Is(err, service.ErrInvitationNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrInvitationNotResendable):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "resend invitation success", nil)
}

func (h *AccessHandler) RevokeInvitation(w http.ResponseWriter, r *http.Request) {
	invID := chi.URLParam(r, "invitationID")

	if err := h.svc.RevokeInvitation(r.Context(), invID); err != nil {
		switch {
		case errors.Is(err, service.ErrInvitationNotRevocable):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "revoke invitation success", nil)
}

func (h *AccessHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	wID := chi.URLParam(r, "workspaceID")

	var req dto.CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = wID

	res, err := h.svc.CreateGroup(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGroupNameTaken):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "create group success", res)
}

func (h *AccessHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")

	res, err := h.svc.GetGroups(r.Context(), wID)
	if err != nil {
		log.Printf("register internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "get groups success", res)
}

func (h *AccessHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	gID := chi.URLParam(r, "groupID")

	if err := h.svc.DeleteGroup(r.Context(), gID); err != nil {
		switch {
		case errors.Is(err, service.ErrGroupNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "delete group success", nil)
}

func (h *AccessHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	gID := chi.URLParam(r, "groupID")

	var req dto.UpdateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.GroupID = gID

	res, err := h.svc.UpdateGroup(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGroupNameTaken):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, service.ErrGroupNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}

		return
	}

	response.Success(w, http.StatusOK, "update group success", res)
}

func (h *AccessHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	gID := chi.URLParam(r, "groupID")

	res, err := h.svc.GetGroupDetail(r.Context(), gID)
	if err != nil {
		log.Printf("register internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "get group detail success", res)
}

func (h *AccessHandler) AssignMember(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	gID := chi.URLParam(r, "groupID")

	var req dto.GroupMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.GroupID = gID

	res, err := h.svc.AssignToGroup(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAssignMemberRole):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("register internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "assign member success", res)
}

func (h *AccessHandler) UnassignMember(w http.ResponseWriter, r *http.Request) {
	gID := chi.URLParam(r, "groupID")
	mID := chi.URLParam(r, "memberID")

	if err := h.svc.UnassignFromGroup(r.Context(), gID, mID); err != nil {
		log.Printf("register internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "unassign member success", nil)
}
