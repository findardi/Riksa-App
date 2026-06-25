package dto

type CreateGroupRequest struct {
	WorkspaceID string `json:"-"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateGroupRequest struct {
	GroupID     string `json:"-"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type GroupMemberRequest struct {
	MemberID []string `json:"member_id" validate:"required"`
	GroupID  string   `json:"-"`
}
