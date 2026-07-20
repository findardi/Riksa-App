package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/findardi/Riksa-App/server/internal/content/dto"
	contentdb "github.com/findardi/Riksa-App/server/internal/content/repository/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *ContentService) RequestUploadURL(ctx context.Context, workspaceID, folderID string) (dto.UploadURLResponse, error) {
	var fID pgtype.UUID
	if err := fID.Scan(folderID); err != nil {
		return dto.UploadURLResponse{}, fmt.Errorf("folder id parse: %w", err)
	}

	folder, err := s.repo.GetFolderByID(ctx, fID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.UploadURLResponse{}, ErrFolderNotFound
	}

	if err != nil {
		return dto.UploadURLResponse{}, fmt.Errorf("get folder: %w", err)
	}

	if uuidString(folder.WorkspaceID) != workspaceID {
		return dto.UploadURLResponse{}, ErrFolderNotFound
	}

	key := storageKey(workspaceID, folderID)
	url, err := s.store.PresignedPut(ctx, key, uploadURLTTL)
	if err != nil {
		return dto.UploadURLResponse{}, fmt.Errorf("presign put: %w", err)
	}

	return dto.UploadURLResponse{
		UploadURL:  url,
		StorageKey: key,
	}, nil
}

func (s *ContentService) CompletedUpload(ctx context.Context, req dto.CompleteUploadRequest) (dto.DocumentResponse, error) {
	var wID, fID, uID pgtype.UUID
	if err := wID.Scan(req.WorkspaceID); err != nil {
		return dto.DocumentResponse{}, fmt.Errorf("workspace id parse: %w", err)
	}

	if err := fID.Scan(req.FolderID); err != nil {
		return dto.DocumentResponse{}, fmt.Errorf("folder id parse: %w", err)
	}

	if err := uID.Scan(req.UploadedBy); err != nil {
		return dto.DocumentResponse{}, fmt.Errorf("user id parse: %w", err)
	}

	folder, err := s.repo.GetFolderByID(ctx, fID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.DocumentResponse{}, ErrFolderNotFound
	}

	if err != nil {
		return dto.DocumentResponse{}, fmt.Errorf("get folder: %w", err)
	}

	if folder.WorkspaceID != wID {
		return dto.DocumentResponse{}, ErrFolderNotFound
	}

	size, mime, err := s.store.Stat(ctx, req.StorageKey)
	if err != nil {
		return dto.DocumentResponse{}, ErrUploadNotFound
	}

	var doc contentdb.Document
	var ver contentdb.DocumentVersion

	err = s.repo.ExecTx(ctx, func(q *contentdb.Queries) error {
		maxPos, err := q.GetMaxPosition(ctx, fID)
		if err != nil {
			return err
		}

		doc, err = q.CreateDocument(ctx, contentdb.CreateDocumentParams{
			WorkspaceID: wID,
			FolderID:    fID,
			Name:        req.Name,
			Position:    maxPos + 1,
			UploadedBy:  uID,
		})

		if err != nil {
			return err
		}

		ver, err = q.CreateDocumentVersion(ctx, contentdb.CreateDocumentVersionParams{
			DocumentID: doc.ID,
			VersionNo:  1,
			Mime:       mime,
			Size:       size,
			StorageKey: req.StorageKey,
			UploadedBy: uID,
		})

		if err != nil {
			return err
		}

		return q.SetCurrentVersion(ctx, contentdb.SetCurrentVersionParams{
			ID:               doc.ID,
			CurrentVersionID: ver.ID,
		})
	})

	if err != nil {
		_ = s.store.Delete(ctx, req.StorageKey)
		return dto.DocumentResponse{}, fmt.Errorf("delete document: %w", err)
	}

	return dto.DocumentResponse{
		ID:        uuidString(doc.ID),
		FolderID:  uuidString(doc.FolderID),
		Name:      doc.Name,
		VersionNo: ver.VersionNo,
		Mime:      ver.Mime,
		Size:      ver.Size,
		CreatedAt: doc.CreatedAt.Time,
		UpdatedAt: doc.UpdatedAt.Time,
	}, nil
}

func (s *ContentService) RequestVersionUpload(ctx context.Context, workspaceID, documentID string) (dto.UploadURLResponse, error) {
	var dID pgtype.UUID
	if err := dID.Scan(documentID); err != nil {
		return dto.UploadURLResponse{}, fmt.Errorf("document id parse: %w", err)
	}

	doc, err := s.repo.GetDocumentByID(ctx, dID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.UploadURLResponse{}, ErrDocumentNotFound
	}

	if err != nil {
		return dto.UploadURLResponse{}, fmt.Errorf("get document: %w", err)
	}

	if uuidString(doc.WorkspaceID) != workspaceID {
		return dto.UploadURLResponse{}, ErrDocumentNotFound
	}

	key := storageKey(uuidString(doc.WorkspaceID), uuidString(doc.FolderID))
	url, err := s.store.PresignedPut(ctx, key, uploadURLTTL)
	if err != nil {
		return dto.UploadURLResponse{}, fmt.Errorf("presign put: %w", err)
	}

	return dto.UploadURLResponse{
		UploadURL:  url,
		StorageKey: key,
	}, nil
}

func (s *ContentService) CompletedVersion(ctx context.Context, req dto.CompleteVersionRequest) (dto.DocumentResponse, error) {
	var dID, uID pgtype.UUID
	if err := dID.Scan(req.DocumentID); err != nil {
		return dto.DocumentResponse{}, fmt.Errorf("document id parse: %w", err)
	}

	if err := uID.Scan(req.UploadedBy); err != nil {
		return dto.DocumentResponse{}, fmt.Errorf("user id parse: %w", err)
	}

	doc, err := s.repo.GetDocumentByID(ctx, dID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.DocumentResponse{}, ErrDocumentNotFound
	}

	if err != nil {
		return dto.DocumentResponse{}, fmt.Errorf("get document: %w", err)
	}

	if uuidString(doc.WorkspaceID) != req.WorkspaceID {
		return dto.DocumentResponse{}, ErrDocumentNotFound
	}

	size, mime, err := s.store.Stat(ctx, req.StorageKey)
	if err != nil {
		return dto.DocumentResponse{}, ErrUploadNotFound
	}

	var ver contentdb.DocumentVersion
	err = s.repo.ExecTx(ctx, func(q *contentdb.Queries) error {
		next, err := q.GetNextVersionNo(ctx, dID)
		if err != nil {
			return err
		}

		ver, err = q.CreateDocumentVersion(ctx, contentdb.CreateDocumentVersionParams{
			DocumentID: dID,
			VersionNo:  next,
			Mime:       mime,
			Size:       size,
			StorageKey: req.StorageKey,
			UploadedBy: uID,
		})

		if err != nil {
			return err
		}

		return q.SetCurrentVersion(ctx, contentdb.SetCurrentVersionParams{
			ID:               dID,
			CurrentVersionID: ver.ID,
		})
	})

	if err != nil {
		_ = s.store.Delete(ctx, req.StorageKey)
		return dto.DocumentResponse{}, fmt.Errorf("create version: %w", err)
	}

	return dto.DocumentResponse{
		ID:        uuidString(doc.ID),
		FolderID:  uuidString(doc.FolderID),
		Name:      doc.Name,
		VersionNo: ver.VersionNo,
		Mime:      ver.Mime,
		Size:      ver.Size,
		CreatedAt: doc.CreatedAt.Time,
		UpdatedAt: ver.CreatedAt.Time,
	}, nil
}

func (s *ContentService) ListDocuments(ctx context.Context, workspaceID, folderID string, actor Actor) ([]dto.DocumentResponse, error) {
	var fID pgtype.UUID
	if err := fID.Scan(folderID); err != nil {
		return []dto.DocumentResponse{}, fmt.Errorf("folder id parse: %w", err)
	}

	folder, err := s.repo.GetFolderByID(ctx, fID)
	if errors.Is(err, pgx.ErrNoRows) {
		return []dto.DocumentResponse{}, ErrFolderNotFound
	}

	if err != nil {
		return []dto.DocumentResponse{}, fmt.Errorf("get folder: %w", err)
	}

	if uuidString(folder.WorkspaceID) != workspaceID {
		return []dto.DocumentResponse{}, ErrFolderNotFound
	}

	if err := s.requireFolderView(ctx, workspaceID, folderID, actor); err != nil {
		return []dto.DocumentResponse{}, err
	}

	rows, err := s.repo.ListDocumentsByFolder(ctx, fID)
	if err != nil {
		return []dto.DocumentResponse{}, fmt.Errorf("list documents: %w", err)
	}

	docs := make([]dto.DocumentResponse, 0, len(rows))
	for _, r := range rows {
		docs = append(docs, dto.DocumentResponse{
			ID:        uuidString(r.ID),
			FolderID:  uuidString(r.FolderID),
			Name:      r.Name,
			VersionNo: r.VersionNo,
			Mime:      r.Mime,
			Size:      r.Size,
			CreatedAt: r.CreatedAt.Time,
			UpdatedAt: r.UpdatedAt.Time,
		})
	}

	return docs, nil
}

func (s *ContentService) ListVersions(ctx context.Context, workspaceID, documentID string, actor Actor) ([]dto.VersionResponse, error) {
	var dID pgtype.UUID
	if err := dID.Scan(documentID); err != nil {
		return []dto.VersionResponse{}, nil
	}

	doc, err := s.repo.GetDocumentByID(ctx, dID)
	if errors.Is(err, pgx.ErrNoRows) {
		return []dto.VersionResponse{}, ErrDocumentNotFound
	}

	if err != nil {
		return []dto.VersionResponse{}, fmt.Errorf("get document: %w", err)
	}

	if uuidString(doc.WorkspaceID) != workspaceID {
		return []dto.VersionResponse{}, ErrDocumentNotFound
	}

	if err := s.requireFolderView(ctx, workspaceID, uuidString(doc.FolderID), actor); err != nil {
		return []dto.VersionResponse{}, err
	}

	rows, err := s.repo.ListVersionByDocument(ctx, dID)
	if err != nil {
		return []dto.VersionResponse{}, fmt.Errorf("list versions: %w", err)
	}

	vers := make([]dto.VersionResponse, 0, len(rows))
	for _, r := range rows {
		vers = append(vers, dto.VersionResponse{
			ID:         uuidString(r.ID),
			VersionNo:  r.VersionNo,
			Mime:       r.Mime,
			Size:       r.Size,
			UploadedBy: uuidString(r.UploadedBy),
			CreatedAt:  r.CreatedAt.Time,
		})
	}

	return vers, nil
}

func (s *ContentService) GetDownloadURL(ctx context.Context, workspaceID, documentID string, actor Actor) (dto.DownloadURLResponse, error) {
	var dID pgtype.UUID
	if err := dID.Scan(documentID); err != nil {
		return dto.DownloadURLResponse{}, fmt.Errorf("document id parse: %w", err)
	}

	doc, err := s.repo.GetDocumentByID(ctx, dID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.DownloadURLResponse{}, ErrDocumentNotFound
	}

	if err != nil {
		return dto.DownloadURLResponse{}, fmt.Errorf("get document: %w", err)
	}

	if uuidString(doc.WorkspaceID) != workspaceID {
		return dto.DownloadURLResponse{}, ErrDocumentNotFound
	}

	if err := s.requireFolderDownloadOriginal(ctx, workspaceID, uuidString(doc.FolderID), actor); err != nil {
		return dto.DownloadURLResponse{}, err
	}

	current, err := s.repo.GetCurrentVersion(ctx, dID)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.DownloadURLResponse{}, ErrDocumentNotFound
	}

	if err != nil {
		return dto.DownloadURLResponse{}, fmt.Errorf("get current version: %w", err)
	}

	url, err := s.store.PresignedGet(ctx, current.StorageKey, doc.Name, downloadURLTTL)
	if err != nil {
		return dto.DownloadURLResponse{}, fmt.Errorf("presign get: %w", err)
	}

	return dto.DownloadURLResponse{
		DownloadURL: url,
	}, nil
}

func (s *ContentService) DeleteDocument(ctx context.Context, workspaceID, documentID string) error {
	var dID pgtype.UUID
	if err := dID.Scan(documentID); err != nil {
		return fmt.Errorf("document id parse: %w", err)
	}

	doc, err := s.repo.GetDocumentByID(ctx, dID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrDocumentNotFound
	}

	if err != nil {
		return fmt.Errorf("get document: %w", err)
	}

	if uuidString(doc.WorkspaceID) != workspaceID {
		return ErrDocumentNotFound
	}

	versions, err := s.repo.ListVersionByDocument(ctx, dID)
	if err != nil {
		return fmt.Errorf("list version: %w", err)
	}

	if err := s.repo.DeleteDocument(ctx, dID); err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	for _, v := range versions {
		_ = s.store.Delete(ctx, v.StorageKey)
	}

	return nil
}

func (s *ContentService) MoveDocument(ctx context.Context, req dto.MoveDocumentRequest) error {
	var dID, fID pgtype.UUID
	if err := dID.Scan(req.DocumentID); err != nil {
		return fmt.Errorf("document id parse: %w", err)
	}
	if err := fID.Scan(req.FolderID); err != nil {
		return fmt.Errorf("folder id parse: %w", err)
	}

	return s.repo.ExecTx(ctx, func(q *contentdb.Queries) error {
		doc, err := q.GetDocumentByID(ctx, dID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDocumentNotFound
		}
		if err != nil {
			return fmt.Errorf("get document: %w", err)
		}
		if uuidString(doc.WorkspaceID) != req.WorkspaceID {
			return ErrDocumentNotFound
		}

		if err := q.LockWorkspaceStructure(ctx, doc.WorkspaceID); err != nil {
			return fmt.Errorf("lock workspace structure: %w", err)
		}

		folder, err := q.GetFolderByID(ctx, fID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrFolderNotFound
		}
		if err != nil {
			return fmt.Errorf("get target folder: %w", err)
		}
		if uuidString(folder.WorkspaceID) != req.WorkspaceID {
			return ErrParentCrossWorkspace
		}

		oldFolder := doc.FolderID

		maxPos, err := q.GetMaxPosition(ctx, fID)
		if err != nil {
			return fmt.Errorf("check max position: %w", err)
		}

		pos := maxPos + 1
		if req.Position != nil {
			pos = clampPosition(int32(*req.Position), maxPos+1)
		}

		if err := q.MoveDocument(ctx, contentdb.MoveDocumentParams{
			ID:       dID,
			FolderID: fID,
			Position: pos,
		}); err != nil {
			return fmt.Errorf("move document: %w", err)
		}

		if err := q.ReindexDocumentSiblings(ctx, contentdb.ReindexDocumentSiblingsParams{
			FolderID: fID,
			MovedID:  dID,
		}); err != nil {
			return fmt.Errorf("reindex target siblings: %w", err)
		}

		if oldFolder != fID {
			if err := q.ReindexDocumentSiblings(ctx, contentdb.ReindexDocumentSiblingsParams{
				FolderID: oldFolder,
				MovedID:  dID,
			}); err != nil {
				return fmt.Errorf("reindex source siblings: %w", err)
			}
		}

		return nil
	})
}
