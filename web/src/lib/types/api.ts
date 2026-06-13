// Shared API contract — mirrors platform/response on the Go backend.

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

/** Normalized result every server-side API call returns. */
export type ApiResult<T> =
	| { ok: true; message: string; data: T }
	| { ok: false; status: number; message: string; fieldErrors: Record<string, string> };
