import type { Handle } from '@sveltejs/kit';
import { getAccessToken } from '$lib/server/session';
import { getMe } from '$lib/server/api';

export const handle: Handle = async ({ event, resolve }) => {
	const token = getAccessToken(event.cookies);
	event.locals.session = token;
	// One /me per request when a session exists — guards read locals.user.status.
	event.locals.user = token ? await getMe(token) : null;
	return resolve(event);
};
