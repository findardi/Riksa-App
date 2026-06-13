import { env } from '$env/dynamic/private';
import { t } from '$lib/i18n';
import type {
	ApiResult,
	Envelope,
	FieldError,
	LoginData,
	LoginPayload,
	RegisterData,
	RegisterPayload
} from '$lib/types';

const API_URL = env.AUTH_API_URL?.replace(/\/$/, '');

export async function registerUser(p: RegisterPayload): Promise<ApiResult<RegisterData>> {
	if (!API_URL) return stubRegister(p);
	return post<RegisterData>('/auth/register', p);
}

export async function loginUser(p: LoginPayload): Promise<ApiResult<LoginData>> {
	if (!API_URL) return stubLogin(p);
	// Backend accepts email OR username (required_without).
	const body = p.identifier.includes('@')
		? { email: p.identifier, password: p.password }
		: { username: p.identifier, password: p.password };
	return post<LoginData>('/auth/login', body);
}

// Early-warning email availability check (step 1 of register).
// 200 → available; 400 (ErrEmailUnique) → taken; 400 validation → bad format.
export async function checkEmailAvailable(
	email: string
): Promise<{ available: boolean; emailError?: string }> {
	if (!API_URL) return stubCheckEmail(email);
	const res = await post<null>('/auth/check-email', { email });
	if (res.ok) return { available: true };
	if (res.fieldErrors.email) return { available: false, emailError: res.fieldErrors.email };
	if (res.status === 400) return { available: false, emailError: t('err.emailTaken') };
	return { available: false, emailError: res.message };
}

// Request a password-reset OTP. Anti-enumeration: backend always 200 regardless
// of whether the email exists, so we surface only network/format problems.
export async function forgotPassword(email: string): Promise<{ sent: boolean; error?: string }> {
	if (!API_URL) {
		await settle();
		return { sent: true };
	}
	const res = await post<null>('/auth/forgot-password', { email });
	if (res.ok) return { sent: true };
	if (res.status === 0) return { sent: false, error: res.message }; // network only
	return { sent: true }; // hide any other backend signal (anti-enum)
}

// Step 2 — validate the reset OTP (read-only check).
export async function validateOtp(email: string, code: string): Promise<{ valid: boolean }> {
	if (!API_URL) return stubValidateOtp(code);
	const res = await post<null>('/auth/validation-otp', { email, code });
	return { valid: res.ok };
}

// Step 3 — set the new password using the verified OTP.
export async function resetPassword(
	email: string,
	code: string,
	newPassword: string
): Promise<{ ok: true } | { ok: false; invalidCode: boolean; message: string }> {
	if (!API_URL) return stubResetPassword(email, code, newPassword);
	const res = await post<null>('/auth/reset-password', { email, code, new_password: newPassword });
	if (res.ok) return { ok: true };
	// Password format is validated client-side, so a 400 here means the code is bad/expired.
	return { ok: false, invalidCode: res.status === 400, message: t('err.invalidOtp') };
}

async function post<T>(path: string, body: unknown): Promise<ApiResult<T>> {
	let res: Response;
	try {
		res = await fetch(`${API_URL}${path}`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify(body)
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

// ---- dev stub -------------------------------------------------------------

interface StubUser {
	id: string;
	email: string;
	username: string;
	password: string;
}

// Seed account so login works out of the box in stub mode.
const users: StubUser[] = [
	{ id: 'usr_demo', email: 'demo@wadi.app', username: 'demowadi', password: 'secret123' }
];
let seq = 1;

const settle = () => new Promise((r) => setTimeout(r, 450));

async function stubRegister(p: RegisterPayload): Promise<ApiResult<RegisterData>> {
	await settle();
	if (users.some((u) => u.email.toLowerCase() === p.email.toLowerCase())) {
		return {
			ok: false,
			status: 409,
			message: t('err.emailTaken'),
			fieldErrors: { email: t('err.emailTaken') }
		};
	}
	if (users.some((u) => u.username.toLowerCase() === p.username.toLowerCase())) {
		return {
			ok: false,
			status: 409,
			message: t('err.usernameTaken'),
			fieldErrors: { username: t('err.usernameTaken') }
		};
	}
	const u: StubUser = {
		id: `usr_${++seq}`,
		email: p.email,
		username: p.username,
		password: p.password
	};
	users.push(u);
	return {
		ok: true,
		message: 'success registered account',
		data: { id: u.id, username: u.username }
	};
}

async function stubCheckEmail(email: string): Promise<{ available: boolean; emailError?: string }> {
	await settle();
	if (users.some((u) => u.email.toLowerCase() === email.toLowerCase())) {
		return { available: false, emailError: t('err.emailTaken') };
	}
	return { available: true };
}

// Demo OTP for stub mode (no backend / no real email).
const STUB_OTP = '123456';

async function stubValidateOtp(code: string): Promise<{ valid: boolean }> {
	await settle();
	return { valid: code === STUB_OTP };
}

async function stubResetPassword(
	email: string,
	code: string,
	newPassword: string
): Promise<{ ok: true } | { ok: false; invalidCode: boolean; message: string }> {
	await settle();
	if (code !== STUB_OTP) return { ok: false, invalidCode: true, message: t('err.invalidOtp') };
	const u = users.find((x) => x.email.toLowerCase() === email.toLowerCase());
	if (u) u.password = newPassword;
	return { ok: true };
}

async function stubLogin(p: LoginPayload): Promise<ApiResult<LoginData>> {
	await settle();
	const idf = p.identifier.toLowerCase();
	const u = users.find((x) => x.email.toLowerCase() === idf || x.username.toLowerCase() === idf);
	if (!u || u.password !== p.password) {
		return { ok: false, status: 401, message: t('err.invalidCredentials'), fieldErrors: {} };
	}
	return {
		ok: true,
		message: 'login success',
		data: { token: `stub.access.${u.id}`, refresh_token: `stub.refresh.${u.id}` }
	};
}
