import { fail, redirect } from '@sveltejs/kit';
import { addMembers, resolveWorkspaceId } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions } from './$types';

export const actions: Actions = {
	// Bulk invite. Emails arrive as repeated `email` fields; the backend decides
	// per email (invite / skip) and never reports who was already registered.
	invite: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const roleId = (form.get('roleId') ?? '').toString();
		const emails = [
			...new Set(
				form
					.getAll('email')
					.map((e) => e.toString().trim().toLowerCase())
					.filter(Boolean)
			)
		];

		if (!emails.length) return fail(400, { message: t('member.invite.empty') });
		if (!roleId) return fail(400, { fieldErrors: { role: t('err.required') } });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await addMembers(locals.session, wsId, { email: emails, role_id: roleId });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			return fail(res.status || 400, {
				fieldErrors: res.fieldErrors?.email ? { email: res.fieldErrors.email } : {},
				message: Object.keys(res.fieldErrors ?? {}).length ? null : res.message || t('err.generic')
			});
		}

		return { invited: true, results: res.data };
	}
};
