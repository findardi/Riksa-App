import { error, json } from '@sveltejs/kit';
import { completeMultipart } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { CompletedPart } from '$lib/types/content';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ locals, request }) => {
	if (!locals.session) error(401, t('err.invalidCredentials'));

	const body = (await request.json().catch(() => null)) as {
		workspaceId?: string;
		folderId?: string;
		uploadId?: string;
		name?: string;
		storageKey?: string;
		contentType?: string;
		parts?: CompletedPart[];
	} | null;

	if (
		!body?.workspaceId ||
		!body.folderId ||
		!body.uploadId ||
		!body.name ||
		!body.storageKey ||
		!body.parts?.length
	) {
		error(400, t('err.generic'));
	}

	const res = await completeMultipart(locals.session, body.workspaceId, body.folderId, {
		upload_id: body.uploadId,
		name: body.name,
		storage_key: body.storageKey,
		content_type: body.contentType || 'application/octet-stream',
		parts: body.parts
	});
	if (!res.ok) error(res.status || 500, res.message);

	return json(res.data);
};
