import { getMyInvitation } from '$lib/server/api';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	let invitationCount = 0;
	if (locals.session) {
		const res = await getMyInvitation(locals.session);
		if (res.ok) invitationCount = res.data.length;
	}
	return { user: locals.user, invitationCount };
};
