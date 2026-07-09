import { error, json } from '@sveltejs/kit';
import { completeUpload } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ locals, request }) => {
	if (!locals.session) error(401, t('err.invalidCredentials'));

	const body = (await request.json().catch(() => null)) as {
		workspaceId?: string;
		folderId?: string;
		name?: string;
		storageKey?: string;
	} | null;

	if (!body?.workspaceId || !body.folderId || !body.name || !body.storageKey) {
		error(400, t('err.generic'));
	}

	const res = await completeUpload(locals.session, body.workspaceId, body.folderId, {
		name: body.name,
		storage_key: body.storageKey
	});
	if (!res.ok) error(res.status || 500, res.message);

	return json(res.data);
};
