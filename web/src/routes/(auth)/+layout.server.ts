import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

// Already authenticated? Route by status: active → app, pending → verify gate.
// Unauthenticated users (locals.user null) fall through to the auth pages.
export const load: LayoutServerLoad = async ({ locals }) => {
	if (locals.user) redirect(303, locals.user.status === 'active' ? '/' : '/verify-email');
};
