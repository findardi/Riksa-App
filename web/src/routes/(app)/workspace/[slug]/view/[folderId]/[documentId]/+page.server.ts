import { error, redirect } from '@sveltejs/kit';
import { getViewMeta } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, params, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();

	const res = await getViewMeta(locals.session, workspace.id, params.documentId);
	if (!res.ok) {
		if (res.status === 401) redirect(303, '/login');
		// Guest whose group has no view access on this folder.
		if (res.status === 403) return { meta: null, forbidden: true, notViewable: false };
		// Download-only formats (spreadsheets, video, archives) have no rendition.
		if (res.status === 422) return { meta: null, forbidden: false, notViewable: true };
		if (res.status === 404) error(404, t('doc.view.err.notFound'));
		error(res.status || 500, t('doc.view.err.load'));
	}

	return { meta: res.data, forbidden: false, notViewable: false };
};
