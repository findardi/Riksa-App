import { error, json } from '@sveltejs/kit';
import { bulkCreateFolders } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { BulkFolderNode } from '$lib/types/content';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ locals, request }) => {
	if (!locals.session) error(401, t('err.invalidCredentials'));

	const body = (await request.json().catch(() => null)) as {
		workspaceId?: string;
		parentId?: string;
		folders?: BulkFolderNode[];
	} | null;

	// `parentId` may legitimately be '' — that means the root of the room.
	if (!body?.workspaceId || body.parentId === undefined || !body.folders?.length) {
		error(400, t('err.generic'));
	}

	const res = await bulkCreateFolders(locals.session, body.workspaceId, {
		parent_id: body.parentId,
		folders: body.folders
	});
	if (!res.ok) error(res.status || 500, res.message);

	return json(res.data);
};
