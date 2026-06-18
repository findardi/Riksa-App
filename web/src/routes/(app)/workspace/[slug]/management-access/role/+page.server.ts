import { error, fail, redirect } from '@sveltejs/kit';
import { deleteRole, getRoles, resolveWorkspaceId } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();
	const res = await getRoles(locals.session, workspace.id);
	if (!res.ok) {
		if (res.status === 401) redirect(303, '/login');
		error(res.status || 502, t('role.err.loadError'));
	}

	return { roles: res.data };
};

export const actions: Actions = {
	delete: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const roleId = (form.get('roleId') ?? '').toString();
		if (!roleId) return fail(400, { message: t('err.generic') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await deleteRole(locals.session, wsId, roleId);
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			// The only 400 on delete is the FK guard: role still assigned to members.
			if (res.status === 400) return fail(400, { message: t('role.err.inUse') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { deleted: true };
	}
};
