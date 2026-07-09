import { error, redirect } from '@sveltejs/kit';
import { getDownloadUrl } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ locals, url }) => {
	if (!locals.session) error(401, t('err.invalidCredentials'));

	const workspaceId = url.searchParams.get('workspaceId');
	const documentId = url.searchParams.get('documentId');
	if (!workspaceId || !documentId) error(400, t('err.generic'));

	const res = await getDownloadUrl(locals.session, workspaceId, documentId);
	if (!res.ok) error(res.status || 500, res.message);

	redirect(302, res.data.download_url);
};
