import { error, json } from '@sveltejs/kit';
import { multipartPartUrls } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ locals, request }) => {
	if (!locals.session) error(401, t('err.invalidCredentials'));

	const body = (await request.json().catch(() => null)) as {
		workspaceId?: string;
		folderId?: string;
		uploadId?: string;
		storageKey?: string;
		partNumbers?: number[];
	} | null;

	if (
		!body?.workspaceId ||
		!body.folderId ||
		!body.uploadId ||
		!body.storageKey ||
		!body.partNumbers?.length
	) {
		error(400, t('err.generic'));
	}

	const res = await multipartPartUrls(locals.session, body.workspaceId, body.folderId, {
		upload_id: body.uploadId,
		storage_key: body.storageKey,
		part_numbers: body.partNumbers
	});
	if (!res.ok) error(res.status || 500, res.message);

	return json(res.data);
};
