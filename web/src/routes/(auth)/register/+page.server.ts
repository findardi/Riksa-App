import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { checkEmailAvailable, registerUser } from '$lib/server/api';
import { t } from '$lib/i18n';

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

// Redirect for already-authenticated users is handled by (auth)/+layout.server.ts.
export const actions: Actions = {
	// Step 1 — early-warning email availability.
	check: async ({ request }) => {
		const data = await request.formData();
		const email = (data.get('email') ?? '').toString().trim();

		const fieldErrors: Record<string, string> = {};
		if (!email) fieldErrors.email = t('err.required');
		else if (!EMAIL_RE.test(email)) fieldErrors.email = t('err.email');
		if (fieldErrors.email) {
			return fail(400, { step: 'check', available: false, fieldErrors, message: undefined });
		}

		const res = await checkEmailAvailable(email);
		if (!res.available) {
			fieldErrors.email = res.emailError ?? t('err.emailTaken');
			return fail(400, { step: 'check', available: false, fieldErrors, message: undefined });
		}

		const empty: Record<string, string> = {};
		return { step: 'check', available: true, fieldErrors: empty, message: undefined };
	},

	// Step 2 — full account creation.
	register: async ({ request }) => {
		const data = await request.formData();
		const email = (data.get('email') ?? '').toString().trim();
		const username = (data.get('username') ?? '').toString().trim();
		const password = (data.get('password') ?? '').toString();

		const fieldErrors: Record<string, string> = {};
		if (!email || !EMAIL_RE.test(email)) fieldErrors.email = t('err.email');
		if (!username) fieldErrors.username = t('err.required');
		else if (username.length < 6) fieldErrors.username = t('err.min', { n: 6 });
		if (!password) fieldErrors.password = t('err.required');
		else if (password.length < 6) fieldErrors.password = t('err.min', { n: 6 });

		if (Object.keys(fieldErrors).length) {
			return fail(400, { step: 'register', available: true, fieldErrors, message: undefined });
		}

		const res = await registerUser({ email, username, password });
		if (!res.ok) {
			return fail(res.status, {
				step: 'register',
				available: true,
				fieldErrors: res.fieldErrors,
				message: Object.keys(res.fieldErrors).length ? undefined : res.message
			});
		}

		redirect(303, '/login?registered=1');
	}
};
