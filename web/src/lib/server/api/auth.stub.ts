import { t } from '$lib/i18n';
import type { ApiResult, LoginData, LoginPayload, RegisterData, RegisterPayload } from '$lib/types';

// In-memory backend used when AUTH_API_URL is unset (local dev / no backend).

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

// Demo OTP for stub mode (no backend / no real email).
const STUB_OTP = '123456';

export async function stubRegister(p: RegisterPayload): Promise<ApiResult<RegisterData>> {
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

export async function stubCheckEmail(
	email: string
): Promise<{ available: boolean; emailError?: string }> {
	await settle();
	if (users.some((u) => u.email.toLowerCase() === email.toLowerCase())) {
		return { available: false, emailError: t('err.emailTaken') };
	}
	return { available: true };
}

export async function stubForgotPassword(): Promise<{ sent: boolean; error?: string }> {
	await settle();
	return { sent: true };
}

export async function stubValidateOtp(code: string): Promise<{ valid: boolean }> {
	await settle();
	return { valid: code === STUB_OTP };
}

export async function stubResetPassword(
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

export async function stubLogin(p: LoginPayload): Promise<ApiResult<LoginData>> {
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
