import { error, redirect } from '@sveltejs/kit';
import { getFoldersTree } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();

	const res = await getFoldersTree(locals.session, workspace.id);
	if (!res.ok) {
		if (res.status === 401) redirect(303, '/login');
		error(res.status || 500, t('doc.err.load'));
	}

	return { folders: res.data ?? [] };
};
