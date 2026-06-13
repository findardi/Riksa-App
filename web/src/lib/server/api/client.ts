import { env } from '$env/dynamic/private';
import { t } from '$lib/i18n';
import type { ApiResult, Envelope, FieldError } from '$lib/types';

// Shared HTTP layer for every feature's API module.
export const API_URL = env.AUTH_API_URL?.replace(/\/$/, '');

// `token` attaches `Authorization: Bearer` for JWT-protected endpoints.
async function request<T>(
	method: 'GET' | 'POST',
	path: string,
	body: unknown,
	token?: string
): Promise<ApiResult<T>> {
	const headers: Record<string, string> = { 'content-type': 'application/json' };
	if (token) headers.authorization = `Bearer ${token}`;

	let res: Response;
	try {
		res = await fetch(`${API_URL}${path}`, {
			method,
			headers,
			body: body === undefined ? undefined : JSON.stringify(body)
		});
	} catch {
		return { ok: false, status: 0, message: t('err.network'), fieldErrors: {} };
	}

	let env: Envelope<T>;
	try {
		env = (await res.json()) as Envelope<T>;
	} catch {
		return { ok: false, status: res.status, message: t('err.generic'), fieldErrors: {} };
	}

	if (res.ok && env.success) {
		return { ok: true, message: env.message, data: env.data as T };
	}

	return {
		ok: false,
		status: res.status,
		message: translateMessage(res.status, env.message),
		fieldErrors: translateFieldErrors(env.errors)
	};
}

export const post = <T>(path: string, body: unknown, token?: string) =>
	request<T>('POST', path, body, token);

export const get = <T>(path: string, token?: string) => request<T>('GET', path, undefined, token);

function translateFieldErrors(errs?: FieldError[] | null): Record<string, string> {
	const out: Record<string, string> = {};
	if (!errs) return out;
	for (const e of errs) out[e.field] = translateFieldMessage(e.message);
	return out;
}

function translateFieldMessage(m: string): string {
	if (m === 'required') return t('err.required');
	if (m === 'invalid email format') return t('err.email');
	let match = m.match(/^minimal (\d+) characters$/);
	if (match) return t('err.min', { n: match[1] });
	match = m.match(/^maximal (\d+) characters$/);
	if (match) return t('err.max', { n: match[1] });
	if (m.startsWith('must fill if')) return t('err.identifierRequired');
	return m;
}

function translateMessage(status: number, raw: string): string {
	if (status === 401) return t('err.invalidCredentials');
	if (status === 409) {
		const m = raw.toLowerCase();
		if (m.includes('email')) return t('err.emailTaken');
		if (m.includes('username')) return t('err.usernameTaken');
	}
	if (status >= 500 || status === 0) return t('err.generic');
	return raw || t('err.generic');
}
