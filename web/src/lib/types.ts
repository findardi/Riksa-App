// API contract — mirrors the Go backend (internal/auth/dto + platform/response).

export interface FieldError {
	field: string;
	message: string;
}

export interface Envelope<T = unknown> {
	success: boolean;
	message: string;
	data?: T;
	errors?: FieldError[] | null;
	meta?: unknown;
}

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

/** Normalized result every server-side API call returns. */
export type ApiResult<T> =
	| { ok: true; message: string; data: T }
	| { ok: false; status: number; message: string; fieldErrors: Record<string, string> };
