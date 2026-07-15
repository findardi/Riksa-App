package dto

import "time"

type UploadURLResponse struct {
	UploadURL  string `json:"upload_url"`
	StorageKey string `json:"storage_key"`
}

type DocumentResponse struct {
	ID        string    `json:"id"`
	FolderID  string    `json:"folder_id"`
	Name      string    `json:"name"`
	VersionNo int32     `json:"version_no"`
	Mime      string    `json:"mime"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VersionResponse struct {
	ID         string    `json:"id"`
	VersionNo  int32     `json:"version_no"`
	Mime       string    `json:"mime"`
	Size       int64     `json:"size"`
	UploadedBy string    `json:"uploaded_by"`
	CreatedAt  time.Time `json:"created_at"`
}

type DownloadURLResponse struct {
	DownloadURL string `json:"download_url"`
}

type ViewMetaResponse struct {
	DocumentID  string `json:"document_id"`
	Name        string `json:"name"`
	Mime        string `json:"mime"`
	PageCount   int    `json:"page_count"`
	CanDownload bool   `json:"can_download"`
}
