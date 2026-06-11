package dto

type RegisterResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
