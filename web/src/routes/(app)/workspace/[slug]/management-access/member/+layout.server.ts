import { error, redirect } from '@sveltejs/kit';
import { getInvitations, getMembers, getRoles } from '$lib/server/api';
import { t } from '$lib/i18n';
import type { LayoutServerLoad } from './$types';

// Loaded once for the member section so both sub-tabs (members + invitations)
// share data and their tab counts always agree.
export const load: LayoutServerLoad = async ({ locals, parent }) => {
	if (!locals.session) redirect(303, '/login');

	const { workspace } = await parent();
	const [membersRes, rolesRes, invitesRes] = await Promise.all([
		getMembers(locals.session, workspace.id),
		getRoles(locals.session, workspace.id),
		getInvitations(locals.session, workspace.id)
	]);

	if (!membersRes.ok) {
		if (membersRes.status === 401) redirect(303, '/login');
		error(membersRes.status || 502, t('member.err.loadError'));
	}
	if (!rolesRes.ok) {
		if (rolesRes.status === 401) redirect(303, '/login');
		error(rolesRes.status || 502, t('member.err.loadError'));
	}
	if (!invitesRes.ok) {
		if (invitesRes.status === 401) redirect(303, '/login');
		error(invitesRes.status || 502, t('pending.err.loadError'));
	}

	return {
		members: membersRes.data,
		roles: rolesRes.data,
		invitations: invitesRes.data,
		ownerId: workspace.owner_id
	};
};
