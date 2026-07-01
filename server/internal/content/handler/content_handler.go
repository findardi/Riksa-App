package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/findardi/Wadi/server/internal/content/dto"
	"github.com/findardi/Wadi/server/internal/content/service"
	"github.com/findardi/Wadi/server/internal/platform/middleware"
	"github.com/findardi/Wadi/server/internal/platform/response"
	"github.com/findardi/Wadi/server/internal/platform/validation"
	"github.com/go-chi/chi/v5"
)

const (
	MaxBodyBytes = 1 << 20
)

type ContentHandler struct {
	svc *service.ContentService
}

func NewContentHandler(svc *service.ContentService) *ContentHandler {
	return &ContentHandler{
		svc: svc,
	}
}

func (h *ContentHandler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	claim, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CreateFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = chi.URLParam(r, "workspaceID")
	req.CreatedBy = claim.ID

	res, err := h.svc.CreateFolder(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrParentNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrParentCrossWorkspace):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrFolderNameTaken):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		default:
			log.Printf("create folder internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusCreated, "create folder success", res)
}

func (h *ContentHandler) MoveFolder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	var req dto.MoveFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	req.FolderID = chi.URLParam(r, "folderID")

	if err := h.svc.MoveFolder(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound), errors.Is(err, service.ErrParentNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrParentCrossWorkspace), errors.Is(err, service.ErrCycle), errors.Is(err, service.ErrFolderTreeTooDeep):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		case errors.Is(err, service.ErrFolderNameTaken):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		default:
			log.Printf("move folder internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "move folder success", nil)
}

func (h *ContentHandler) GetFoldersTree(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")

	res, err := h.svc.GetFoldersTree(r.Context(), wID)
	if err != nil {
		log.Printf("move folder internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "get folders tree success", res)
}

func (h *ContentHandler) RenameFolder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	var req dto.RenameFolderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.FolderID = chi.URLParam(r, "folderID")

	res, err := h.svc.RenameFolder(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNameTaken):
			response.Error(w, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("move folder internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "rename folder success", res)
}

func (h *ContentHandler) DeleteFolder(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.DeleteFolder(r.Context(), chi.URLParam(r, "folderID")); err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("move folder internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "delete folder success", nil)
}
