package dto

type CreateWorkspaceRoleRequest struct {
	WorkspaceID string   `json:"-"`
	Permission  []string `json:"permissions" validate:"required"`
	Name        string   `json:"name" validate:"required"`
}

type UpdateWorkspaceRoleRequest struct {
	WorkspaceID string   `json:"-"`
	RoleID      string   `json:"-"`
	Name        string   `json:"name"`
	Permission  []string `json:"permissions"`
}
