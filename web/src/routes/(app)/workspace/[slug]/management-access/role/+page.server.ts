import { error, redirect } from '@sveltejs/kit';
import { getRoles } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { PageServerLoad } from './$types';

// Roles are fixed system roles (owner/admin/guest); this page is read-only.
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
