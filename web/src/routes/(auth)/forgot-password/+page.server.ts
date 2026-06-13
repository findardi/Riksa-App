import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { forgotPassword, validateOtp, resetPassword } from '$lib/server/api';
import { t } from '$lib/i18n';

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const OTP_RE = /^\d{6}$/;

// Redirect for already-authenticated users is handled by (auth)/+layout.server.ts.
export const actions: Actions = {
	// Step 1 — send (and resend) the reset OTP. Anti-enumeration: always "sent".
	send: async ({ request }) => {
		const data = await request.formData();
		const email = (data.get('email') ?? '').toString().trim();

		const fieldErrors: Record<string, string> = {};
		if (!email) fieldErrors.email = t('err.required');
		else if (!EMAIL_RE.test(email)) fieldErrors.email = t('err.email');
		if (fieldErrors.email) {
			return fail(400, { sent: false, fieldErrors, message: undefined });
		}

		const res = await forgotPassword(email);
		if (!res.sent) {
			return fail(502, { sent: false, fieldErrors, message: res.error ?? t('err.network') });
		}

		return { sent: true, fieldErrors, message: undefined };
	},

	// Step 2 — validate the OTP. Email + code arrive as hidden fields.
	verify: async ({ request }) => {
		const data = await request.formData();
		const email = (data.get('email') ?? '').toString().trim();
		const code = (data.get('code') ?? '').toString().trim();

		const fieldErrors: Record<string, string> = {};
		if (!OTP_RE.test(code) || !EMAIL_RE.test(email)) fieldErrors.code = t('err.invalidOtp');
		if (fieldErrors.code) {
			return fail(400, { valid: false, fieldErrors, message: undefined });
		}

		const res = await validateOtp(email, code);
		if (!res.valid) {
			fieldErrors.code = t('err.invalidOtp');
			return fail(400, { valid: false, fieldErrors, message: undefined });
		}

		return { valid: true, fieldErrors, message: undefined };
	},

	// Step 3 — set the new password. Email + code stay hidden, still sent.
	reset: async ({ request }) => {
		const data = await request.formData();
		const email = (data.get('email') ?? '').toString().trim();
		const code = (data.get('code') ?? '').toString().trim();
		const password = (data.get('password') ?? '').toString();
		const confirm = (data.get('confirm') ?? '').toString();

		const fieldErrors: Record<string, string> = {};
		if (!password) fieldErrors.password = t('err.required');
		else if (password.length < 6) fieldErrors.password = t('err.min', { n: 6 });
		if (confirm !== password) fieldErrors.confirm = t('reset.mismatch');
		if (Object.keys(fieldErrors).length) {
			return fail(400, { invalidCode: false, fieldErrors, message: undefined });
		}

		const res = await resetPassword(email, code, password);
		if (!res.ok) {
			return fail(400, { invalidCode: res.invalidCode, fieldErrors, message: res.message });
		}

		redirect(303, '/login?reset=1');
	}
};
