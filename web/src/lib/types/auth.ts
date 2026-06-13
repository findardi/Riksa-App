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

/** Account status — `pending` until email is verified, then `active`. */
export type AccountStatus = 'pending' | 'active';

/** Current authenticated user, from GET /auth/me. */
export interface MeData {
	id: string;
	email: string;
	username: string;
	status: AccountStatus;
}
