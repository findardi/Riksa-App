package dto

type CreateWorkspaceMemberRequest struct {
	WorkspaceId string `json:"workspace_id" validate:"required"`
	UserId      string `json:"user_id" validate:"required"`
	RoleId      string `json:"role_id" validate:"required"`
	Status      string `json:"status" validate:"required"`
}
