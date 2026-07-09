package dto

import "time"

type FolderResponse struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspace_id"`
	ParentID    string    `json:"parent_id"`
	Name        string    `json:"name"`
	Position    int32     `json:"position"`
	IsDefault   bool      `json:"is_default"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FolderTreeNode struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Number    string           `json:"number"`
	Position  int32            `json:"position"`
	IsDefault bool             `json:"is_default"`
	Children  []FolderTreeNode `json:"children"`
}
