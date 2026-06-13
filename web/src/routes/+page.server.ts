import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { logoutUser } from '$lib/server/api';
import { clearSession, getRefreshToken } from '$lib/server/session';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.user) redirect(303, '/login');
	// VDR gate: unverified accounts cannot reach app features.
	if (locals.user.status === 'pending') redirect(303, '/verify-email');
};

export const actions: Actions = {
	logout: async ({ locals, cookies }) => {
		const refresh = getRefreshToken(cookies);
		if (locals.session && refresh) await logoutUser(locals.session, refresh);
		clearSession(cookies);
		redirect(303, '/login');
	}
};
