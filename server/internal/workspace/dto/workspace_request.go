package dto

type WorkspaceCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	OwnerID     string `json:"-"`
}

type GetWorkspace struct {
	Slug    string `json:"slug" validate:"required"`
	OwnerID string `json:"-"`
}

type WorkspaceUpdateStatusRequest struct {
	ID     string `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type WorkspaceUpdateRequest struct {
	ID          string `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type WorkspaceDeleteRequest struct {
	ID string `json:"id" validate:"required"`
}
