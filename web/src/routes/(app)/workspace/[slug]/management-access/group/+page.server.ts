import { error, fail, redirect } from '@sveltejs/kit';
import {
	createGroup,
	deleteGroup,
	getGroups,
	resolveWorkspaceId,
	updateGroup
} from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();
	const res = await getGroups(locals.session, workspace.id);
	if (!res.ok) {
		if (res.status === 401) redirect(303, '/login');
		error(res.status || 502, t('group.err.loadError'));
	}

	return { groups: res.data ?? [] };
};

export const actions: Actions = {
	create: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const name = (form.get('name') ?? '').toString().trim();
		const description = (form.get('description') ?? '').toString().trim();
		if (!name) return fail(400, { message: t('group.err.nameRequired') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await createGroup(locals.session, wsId, { name, description });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 409) return fail(409, { message: t('group.err.nameTaken') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { created: true };
	},

	update: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const groupId = (form.get('groupId') ?? '').toString();
		const name = (form.get('name') ?? '').toString().trim();
		const description = (form.get('description') ?? '').toString().trim();
		if (!groupId || !name) return fail(400, { message: t('group.err.nameRequired') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await updateGroup(locals.session, wsId, groupId, { name, description });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 404) return fail(404, { message: t('group.err.notFound') });
			if (res.status === 409) return fail(409, { message: t('group.err.nameTaken') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { updated: true };
	},

	delete: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const groupId = (form.get('groupId') ?? '').toString();
		if (!groupId) return fail(400, { message: t('err.generic') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await deleteGroup(locals.session, wsId, groupId);
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 404) return fail(404, { message: t('group.err.notFound') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { deleted: true };
	}
};
