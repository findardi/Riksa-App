import type { LayoutServerLoad } from './$types';

// Expose the authenticated user to the app shell (top bar account menu).
export const load: LayoutServerLoad = async ({ locals }) => {
	return { user: locals.user };
};
