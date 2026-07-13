import { fail, redirect } from '@sveltejs/kit';
import {
	createFolder,
	deleteFolder,
	moveFolder,
	removeFolderAccess,
	renameFolder,
	resolveWorkspaceId,
	setFolderAccess
} from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions } from './$types';

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
			if (res.status === 403) return fail(403, { message: t('doc.err.defaultLocked') });
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
			if (res.status === 403) return fail(403, { message: t('doc.err.defaultLocked') });
			if (res.status === 404) return fail(404, { message: t('doc.err.notFound') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { deleted: true };
	},

	setAccess: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const folderId = (form.get('folderId') ?? '').toString();
		const groupId = (form.get('groupId') ?? '').toString();
		const levelId = (form.get('levelId') ?? '').toString();
		if (!folderId) return fail(400, { message: t('err.generic') });
		if (!groupId || !levelId) return fail(400, { message: t('facc.err.pick') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await setFolderAccess(locals.session, wsId, folderId, {
			group_id: groupId,
			level_id: levelId
		});
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 404) return fail(404, { message: t('facc.err.notFound') });
			if (res.status === 400) return fail(400, { message: t('facc.err.invalid') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { accessSet: true };
	},

	removeAccess: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const folderId = (form.get('folderId') ?? '').toString();
		const groupId = (form.get('groupId') ?? '').toString();
		if (!folderId || !groupId) return fail(400, { message: t('err.generic') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await removeFolderAccess(locals.session, wsId, folderId, groupId);
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 404) return fail(404, { message: t('facc.err.notFound') });
			if (res.status === 400) return fail(400, { message: t('facc.err.invalid') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { accessRemoved: true };
	}
};
