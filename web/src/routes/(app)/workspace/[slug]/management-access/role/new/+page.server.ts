import { error, fail, redirect } from '@sveltejs/kit';
import { createRole, getPermissions, resolveWorkspaceId } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.session) redirect(303, '/login');

	const res = await getPermissions(locals.session);
	if (!res.ok) {
		if (res.status === 401) redirect(303, '/login');
		error(res.status || 502, t('role.err.loadError'));
	}

	return { catalog: res.data };
};

export const actions: Actions = {
	default: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const name = (form.get('name') ?? '').toString().trim();
		const permissions = form.getAll('permissions').map(String);

		// Mirror the backend rules (name + at least one permission) to save a round trip.
		const fieldErrors: Record<string, string> = {};
		if (!name) fieldErrors.name = t('err.required');
		if (name && permissions.length === 0) {
			return fail(400, { message: t('role.err.noPermissions'), fieldErrors });
		}
		if (!name) return fail(400, { message: null, fieldErrors });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound'), fieldErrors });

		const res = await createRole(locals.session, wsId, { name, permissions, is_system: false });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 409) {
				fieldErrors.name = t('role.err.nameTaken');
				return fail(409, { message: null, fieldErrors });
			}
			return fail(res.status || 400, { message: res.message || t('err.generic'), fieldErrors });
		}

		redirect(303, `/workspace/${params.slug}/management-access/role?flash=created`);
	}
};
