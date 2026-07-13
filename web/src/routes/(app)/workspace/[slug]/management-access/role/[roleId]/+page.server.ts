import { error, redirect } from '@sveltejs/kit';
import { getPermissions, getRole } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { PageServerLoad } from './$types';

// Read-only role detail: shows the role's granted permissions. Roles are fixed
// system roles (owner/admin/guest) — there is no edit surface.
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
