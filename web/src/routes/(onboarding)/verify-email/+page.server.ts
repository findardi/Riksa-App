import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { logoutUser, resendOtp, verifyEmail } from '$lib/server/api';
import { clearSession, getRefreshToken } from '$lib/server/session';
import { t } from '$lib/i18n';

const OTP_RE = /^\d{6}$/;

export const actions: Actions = {
	verify: async ({ request, locals, cookies }) => {
		if (!locals.session) redirect(303, '/login');
		const data = await request.formData();
		const code = (data.get('code') ?? '').toString().trim();

		const fieldErrors: Record<string, string> = {};
		if (!OTP_RE.test(code)) fieldErrors.code = t('err.invalidOtp');
		if (fieldErrors.code) return fail(400, { fieldErrors, message: undefined });

		const res = await verifyEmail(locals.session, code);
		if (!res.ok) return fail(400, { fieldErrors, message: res.message });

		// Status changed (pending → active), so the existing JWT carries stale claims.
		// Revoke the refresh token server-side, drop the cookies, and force a fresh
		// login to mint tokens with the new status.
		const refresh = getRefreshToken(cookies);
		if (refresh) await logoutUser(locals.session, refresh);
		clearSession(cookies);
		redirect(303, '/login?verified=1');
	},

	resend: async ({ locals }) => {
		if (!locals.session) redirect(303, '/login');
		const fieldErrors: Record<string, string> = {};
		const res = await resendOtp(locals.session);
		if (!res.sent) {
			return fail(502, { fieldErrors, message: res.error ?? t('err.network') });
		}
		return { resent: true, fieldErrors, message: undefined };
	},

	logout: async ({ locals, cookies }) => {
		const refresh = getRefreshToken(cookies);
		if (locals.session && refresh) await logoutUser(locals.session, refresh);
		clearSession(cookies);
		redirect(303, '/login');
	}
};
