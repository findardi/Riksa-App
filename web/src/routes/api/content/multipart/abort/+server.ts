import { error, json } from '@sveltejs/kit';
import { abortMultipart } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { RequestHandler } from './$types';

export const DELETE: RequestHandler = async ({ locals, request }) => {
	if (!locals.session) error(401, t('err.invalidCredentials'));

	const body = (await request.json().catch(() => null)) as {
		workspaceId?: string;
		folderId?: string;
		uploadId?: string;
		storageKey?: string;
	} | null;

	if (!body?.workspaceId || !body.folderId || !body.uploadId || !body.storageKey) {
		error(400, t('err.generic'));
	}

	const res = await abortMultipart(locals.session, body.workspaceId, body.folderId, {
		upload_id: body.uploadId,
		storage_key: body.storageKey
	});
	if (!res.ok) error(res.status || 500, res.message);

	return json({ ok: true });
};
