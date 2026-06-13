import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { loginUser } from '$lib/server/api';
import { setSession } from '$lib/server/session';
import { t } from '$lib/i18n';

export const load: PageServerLoad = async ({ locals }) => {
	if (locals.session) redirect(303, '/');
};

export const actions: Actions = {
	default: async ({ request, cookies }) => {
		const data = await request.formData();
		const identifier = (data.get('identifier') ?? '').toString().trim();
		const password = (data.get('password') ?? '').toString();

		const fieldErrors: Record<string, string> = {};
		if (!identifier) fieldErrors.identifier = t('err.identifierRequired');
		if (!password) fieldErrors.password = t('err.required');
		if (Object.keys(fieldErrors).length) {
			return fail(400, { values: { identifier }, fieldErrors, message: undefined });
		}

		const res = await loginUser({ identifier, password });
		if (!res.ok) {
			return fail(res.status === 401 ? 401 : 400, {
				values: { identifier },
				fieldErrors: res.fieldErrors,
				message: res.message
			});
		}

		setSession(cookies, res.data);
		redirect(303, '/');
	}
};
