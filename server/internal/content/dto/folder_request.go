package dto

type CreateFolderRequest struct {
	WorkspaceID string `json:"-"`
	CreatedBy   string `json:"-"`
	ParentID    string `json:"parent_id"`
	Name        string `json:"name" validate:"required"`
}

type MoveFolderRequest struct {
	FolderID string `json:"-"`
	ParentID string `json:"parent_id"`
}

type RenameFolderRequest struct {
	FolderID string `json:"-"`
	Name     string `json:"name" validate:"required"`
}

type SetFolderAccessRequest struct {
	WorkspaceID string `json:"-"`
	FolderID    string `json:"-"`
	GroupID     string `json:"group_id" validate:"required,uuid"`
	LevelID     string `json:"level_id" validate:"required,uuid"`
}
