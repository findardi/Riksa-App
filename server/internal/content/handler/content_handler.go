package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/findardi/Riksa-App/server/internal/content/dto"
	"github.com/findardi/Riksa-App/server/internal/content/service"
	"github.com/findardi/Riksa-App/server/internal/platform/middleware"
	"github.com/findardi/Riksa-App/server/internal/platform/response"
	"github.com/findardi/Riksa-App/server/internal/platform/validation"
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
		case errors.Is(err, service.ErrMoveDefault):
			response.Error(w, http.StatusForbidden, err.Error(), nil)
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
		log.Printf("get folders tree internal error: %v", err)
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
			log.Printf("rename folder internal error: %v", err)
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
		case errors.Is(err, service.ErrDeleteDefault):
			response.Error(w, http.StatusForbidden, err.Error(), nil)
		default:
			log.Printf("delete folder internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "delete folder success", nil)
}

func (h *ContentHandler) RequestUploadURL(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")
	fID := chi.URLParam(r, "folderID")

	res, err := h.svc.RequestUploadURL(r.Context(), wID, fID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("request upload url internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "request upload url success", res)
}

func (h *ContentHandler) CompletedUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CompleteUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = chi.URLParam(r, "workspaceID")
	req.FolderID = chi.URLParam(r, "folderID")
	req.UploadedBy = claims.ID

	res, err := h.svc.CompletedUpload(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrUploadNotFound):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("complete upload internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusCreated, "upload document success", res)
}

func (h *ContentHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")
	fID := chi.URLParam(r, "folderID")

	res, err := h.svc.ListDocuments(r.Context(), wID, fID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("list documents internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "list documents success", res)
}

func (h *ContentHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")
	dID := chi.URLParam(r, "documentID")

	res, err := h.svc.ListVersions(r.Context(), wID, dID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDocumentNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("list versions internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "list versions success", res)
}

func (h *ContentHandler) RequestUploadVersion(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")
	dID := chi.URLParam(r, "documentID")

	res, err := h.svc.RequestVersionUpload(r.Context(), wID, dID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDocumentNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("request version upload internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "request version upload url success", res)
}

func (h *ContentHandler) CompletedVersionUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req dto.CompleteVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = chi.URLParam(r, "workspaceID")
	req.DocumentID = chi.URLParam(r, "documentID")
	req.UploadedBy = claims.ID

	res, err := h.svc.CompletedVersion(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDocumentNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrUploadNotFound):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("complete version internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusCreated, "upload version success", res)
}

func (h *ContentHandler) GetDownloadURL(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")
	dID := chi.URLParam(r, "documentID")

	res, err := h.svc.GetDownloadURL(r.Context(), wID, dID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDocumentNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("get download url internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "get download url success", res)
}

func (h *ContentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	wID := chi.URLParam(r, "workspaceID")
	dID := chi.URLParam(r, "documentID")

	if err := h.svc.DeleteDocument(r.Context(), wID, dID); err != nil {
		switch {
		case errors.Is(err, service.ErrDocumentNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("delete document internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "delete document success", nil)
}

func (h *ContentHandler) MoveDocument(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	var req dto.MoveDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = chi.URLParam(r, "workspaceID")
	req.DocumentID = chi.URLParam(r, "documentID")

	if err := h.svc.MoveDocument(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, service.ErrDocumentNotFound), errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrParentCrossWorkspace):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("move document internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "move document success", nil)
}

func (h *ContentHandler) SetFolderAccess(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	var req dto.SetFolderAccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body request", nil)
		return
	}

	if errs := validation.Validate(&req); errs != nil {
		response.Error(w, http.StatusBadRequest, "validation failed", errs)
		return
	}

	req.WorkspaceID = chi.URLParam(r, "workspaceID")
	req.FolderID = chi.URLParam(r, "folderID")

	if err := h.svc.SetFolderAccess(r.Context(), req); err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrAccessTargetInvalid):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("set folder access internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "success set folder access", nil)
}

func (h *ContentHandler) RemoveFolderAccess(w http.ResponseWriter, r *http.Request) {
	WorkspaceID := chi.URLParam(r, "workspaceID")
	FolderID := chi.URLParam(r, "folderID")
	groupID := chi.URLParam(r, "groupID")

	if err := h.svc.RemoveFolderAccess(r.Context(), WorkspaceID, groupID, FolderID); err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		case errors.Is(err, service.ErrAccessTargetInvalid):
			response.Error(w, http.StatusBadRequest, err.Error(), nil)
		default:
			log.Printf("remove folder access internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "success remove folder access", nil)
}

func (h *ContentHandler) ListAccessLevel(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.ListAccessLevels(r.Context(), chi.URLParam(r, "workspaceID"))
	if err != nil {
		log.Printf("list access level internal error: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, "success get list level access", res)
}

func (h *ContentHandler) ListFolderAccess(w http.ResponseWriter, r *http.Request) {
	WorkspaceID := chi.URLParam(r, "workspaceID")
	FolderID := chi.URLParam(r, "folderID")

	res, err := h.svc.ListFolderAccess(r.Context(), WorkspaceID, FolderID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFolderNotFound):
			response.Error(w, http.StatusNotFound, err.Error(), nil)
		default:
			log.Printf("list folder access internal error: %v", err)
			response.Error(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	response.Success(w, http.StatusOK, "success get list folder access", res)
}
