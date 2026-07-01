import { error, fail, redirect } from '@sveltejs/kit';
import {
	createFolder,
	deleteFolder,
	getFoldersTree,
	moveFolder,
	renameFolder,
	resolveWorkspaceId
} from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();

	const res = await getFoldersTree(locals.session, workspace.id);
	if (!res.ok) {
		if (res.status === 401) redirect(303, '/login');
		error(res.status || 500, t('doc.err.load'));
	}

	return { folders: res.data ?? [] };
};

export const actions: Actions = {
	create: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const name = (form.get('name') ?? '').toString().trim();
		const parentId = (form.get('parentId') ?? '').toString();
		if (!name) return fail(400, { message: t('doc.err.nameRequired') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await createFolder(locals.session, wsId, { name, parent_id: parentId });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 409) return fail(409, { message: t('doc.err.nameTaken') });
			if (res.status === 404) return fail(404, { message: t('doc.err.parentNotFound') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { created: true };
	},

	rename: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const folderId = (form.get('folderId') ?? '').toString();
		const name = (form.get('name') ?? '').toString().trim();
		if (!folderId || !name) return fail(400, { message: t('doc.err.nameRequired') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await renameFolder(locals.session, wsId, folderId, { name });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 409) return fail(409, { message: t('doc.err.nameTaken') });
			if (res.status === 404) return fail(404, { message: t('doc.err.notFound') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { renamed: true };
	},

	move: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const folderId = (form.get('folderId') ?? '').toString();
		const parentId = (form.get('parentId') ?? '').toString();
		if (!folderId) return fail(400, { message: t('err.generic') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await moveFolder(locals.session, wsId, folderId, { parent_id: parentId });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 409) return fail(409, { message: t('doc.err.nameTaken') });
			if (res.status === 404) return fail(404, { message: t('doc.err.notFound') });
			if (res.status === 400) return fail(400, { message: t('doc.err.invalidMove') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { moved: true };
	},

	delete: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const folderId = (form.get('folderId') ?? '').toString();
		if (!folderId) return fail(400, { message: t('err.generic') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await deleteFolder(locals.session, wsId, folderId);
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 404) return fail(404, { message: t('doc.err.notFound') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { deleted: true };
	}
};
