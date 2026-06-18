import { error, fail, redirect } from '@sveltejs/kit';
import { getPermissions, getRole, resolveWorkspaceId, updateRole } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, params, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();
	const [roleRes, permRes] = await Promise.all([
		getRole(locals.session, workspace.id, params.roleId),
		getPermissions(locals.session)
	]);

	if (!roleRes.ok) {
		if (roleRes.status === 401) redirect(303, '/login');
		if (roleRes.status === 404) error(404, t('role.err.notFound'));
		error(roleRes.status || 502, t('role.err.loadError'));
	}
	if (!permRes.ok) {
		if (permRes.status === 401) redirect(303, '/login');
		error(permRes.status || 502, t('role.err.loadError'));
	}

	return { role: roleRes.data, catalog: permRes.data };
};

export const actions: Actions = {
	default: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const name = (form.get('name') ?? '').toString().trim();
		const permissions = form.getAll('permissions').map(String);

		const fieldErrors: Record<string, string> = {};
		if (!name) fieldErrors.name = t('err.required');
		if (name && permissions.length === 0) {
			return fail(400, { message: t('role.err.noPermissions'), fieldErrors });
		}
		if (!name) return fail(400, { message: null, fieldErrors });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound'), fieldErrors });

		// System roles are read-only — the backend doesn't enforce this yet, so guard here.
		const existing = await getRole(locals.session, wsId, params.roleId);
		if (existing.ok && existing.data.is_system) {
			return fail(403, { message: t('role.view.systemNote'), fieldErrors });
		}

		const res = await updateRole(locals.session, wsId, params.roleId, { name, permissions });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 404) return fail(404, { message: t('role.err.notFound'), fieldErrors });
			if (res.status === 409) {
				fieldErrors.name = t('role.err.nameTaken');
				return fail(409, { message: null, fieldErrors });
			}
			return fail(res.status || 400, { message: res.message || t('err.generic'), fieldErrors });
		}

		redirect(303, `/workspace/${params.slug}/management-access/role?flash=updated`);
	}
};
