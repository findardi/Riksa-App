import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

// Gate for authenticated-but-pending accounts. No session → login;
// already active → app. Pending users stay here to verify.
export const load: LayoutServerLoad = async ({ locals }) => {
	if (!locals.user) redirect(303, '/login');
	if (locals.user.status === 'active') redirect(303, '/');
	return { email: locals.user.email };
};
