import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { clearSession } from '$lib/server/session';

export const load: PageServerLoad = async ({ locals }) => {
	if (!locals.session) redirect(303, '/login');
};

export const actions: Actions = {
	// Local sign-out (clears the cookie). Backend POST /auth/logout wiring is a
	// follow-up pass — out of this craft's login+register scope.
	logout: async ({ cookies }) => {
		clearSession(cookies);
		redirect(303, '/login');
	}
};
