import { error, redirect } from '@sveltejs/kit';
import { getFoldersTree, getGroups } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { GroupWorkspaceData } from '$lib/types/workspace';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace, access } = await parent();

	const res = await getFoldersTree(locals.session, workspace.id);
	if (!res.ok) {
		if (res.status === 401) redirect(303, '/login');
		if (res.status === 403) {
			return { folders: [], groups: [], noAccess: true, accessReady: false };
		}
		error(res.status || 500, t('doc.err.load'));
	}

	const canAssign = access?.permissions?.includes('group:assign') ?? false;
	if (!canAssign) {
		return {
			folders: res.data ?? [],
			groups: [],
			noAccess: false,
			accessReady: false
		};
	}

	const groupsRes = await getGroups(locals.session, workspace.id);
	const groups: GroupWorkspaceData[] = groupsRes.ok ? (groupsRes.data ?? []) : [];

	return {
		folders: res.data ?? [],
		groups,
		noAccess: false,
		accessReady: groupsRes.ok
	};
};
