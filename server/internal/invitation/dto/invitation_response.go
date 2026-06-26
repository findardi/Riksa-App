package dto

import "time"

type GetMyInvitationsRow struct {
	ID            string    `json:"id"`
	WorkspaceName string    `json:"workspace_name"`
	RoleName      string    `json:"role_name"`
	InvitedBy     string    `json:"invited_by"`
	ExpiresAt     time.Time `json:"expires_at"`
	Status        string    `json:"status"`
}
