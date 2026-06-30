package dto

type RegisterResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Status        string `json:"status"`
	EmailVerified bool   `json:"email_verified"`
}

type InvitePreviewResponse struct {
	Email         string `json:"email"`
	WorkspaceName string `json:"workspace_name"`
	RoleName      string `json:"role_name"`
}
