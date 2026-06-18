import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

// The index has no content of its own — land on the first tab.
export const load: PageServerLoad = ({ params }) => {
	redirect(307, `/workspace/${params.slug}/management-access/member`);
};
