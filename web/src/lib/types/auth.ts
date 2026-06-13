// Auth feature contract — mirrors internal/auth/dto on the Go backend.

export interface RegisterPayload {
	email: string;
	username: string;
	password: string;
}

export interface RegisterData {
	id: string;
	username: string;
}

export interface LoginPayload {
	/** email or username — backend accepts either (required_without). */
	identifier: string;
	password: string;
}

export interface LoginData {
	token: string;
	refresh_token: string;
}
