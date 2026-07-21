import { error, json } from '@sveltejs/kit';
import { multipartParts } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { RequestHandler } from './$types';

// The resume read: object storage is the source of truth for which parts landed,
// so a client that kept only the handle can still work out what is left to send.
export const GET: RequestHandler = async ({ locals, url }) => {
	if (!locals.session) error(401, t('err.invalidCredentials'));

	const workspaceId = url.searchParams.get('workspaceId');
	const folderId = url.searchParams.get('folderId');
	const uploadId = url.searchParams.get('uploadId');
	const storageKey = url.searchParams.get('storageKey');

	if (!workspaceId || !folderId || !uploadId || !storageKey) error(400, t('err.generic'));

	const res = await multipartParts(locals.session, workspaceId, folderId, uploadId, storageKey);
	if (!res.ok) error(res.status || 500, res.message);

	return json(res.data);
};
