package dto

type CreateWorkspaceMemberRequest struct {
	WorkspaceId string `json:"-"`
	UserId      string `json:"user_id" validate:"required"`
	RoleId      string `json:"role_id" validate:"required"`
	Status      string `json:"status"`
}

type UpdateMemberRoleRequest struct {
	MemberID string `json:"-"`
	RoleId   string `json:"role_id" validate:"required"`
}

type CheckEmailRequest struct {
	Email string `json:"email" validate:"required"`
}
