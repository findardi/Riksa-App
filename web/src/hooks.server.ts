import type { Handle } from '@sveltejs/kit';
import { getAccessToken } from '$lib/server/session';

export const handle: Handle = async ({ event, resolve }) => {
	event.locals.session = getAccessToken(event.cookies);
	return resolve(event);
};
