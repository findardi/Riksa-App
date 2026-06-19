import { error, fail, redirect } from '@sveltejs/kit';
import {
	checkUser,
	deleteMember,
	getMembers,
	getRoles,
	resolveWorkspaceId,
	updateMemberRole
} from '$lib/server/api';
import { t } from '$lib/i18n';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();
	const [membersRes, rolesRes] = await Promise.all([
		getMembers(locals.session, workspace.id),
		getRoles(locals.session, workspace.id)
	]);

	if (!membersRes.ok) {
		if (membersRes.status === 401) redirect(303, '/login');
		error(membersRes.status || 502, t('member.err.loadError'));
	}
	if (!rolesRes.ok) {
		if (rolesRes.status === 401) redirect(303, '/login');
		error(rolesRes.status || 502, t('member.err.loadError'));
	}

	return { members: membersRes.data, roles: rolesRes.data, ownerId: workspace.owner_id };
};

export const actions: Actions = {
	// Email lookup for the add/invite flow — returns whether the user exists.
	check: async ({ locals, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const email = (form.get('email') ?? '').toString().trim();
		if (!email) return fail(400, { fieldErrors: { email: t('err.required') } });

		const res = await checkUser(locals.session, email);
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			return fail(res.status || 400, {
				fieldErrors: (res.fieldErrors?.email ? { email: res.fieldErrors.email } : {}) as Record<
					string,
					string
				>,
				message: Object.keys(res.fieldErrors ?? {}).length ? null : res.message || t('err.generic')
			});
		}

		return { checked: true, email, exists: res.data };
	},

	updateRole: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const memberId = (form.get('memberId') ?? '').toString();
		const roleId = (form.get('roleId') ?? '').toString();
		if (!memberId || !roleId) return fail(400, { message: t('err.generic') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await updateMemberRole(locals.session, wsId, memberId, { role_id: roleId });
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			if (res.status === 404) return fail(404, { message: t('member.err.notFound') });
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { roleUpdated: true };
	},

	delete: async ({ locals, params, request }) => {
		if (!locals.session) redirect(303, '/login');

		const form = await request.formData();
		const memberId = (form.get('memberId') ?? '').toString();
		if (!memberId) return fail(400, { message: t('err.generic') });

		const wsId = await resolveWorkspaceId(locals.session, params.slug);
		if (!wsId) return fail(404, { message: t('ws.detail.notFound') });

		const res = await deleteMember(locals.session, wsId, memberId);
		if (!res.ok) {
			if (res.status === 401) redirect(303, '/login');
			return fail(res.status || 400, { message: res.message || t('err.generic') });
		}

		return { deleted: true };
	}
};
