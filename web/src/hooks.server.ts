import type { Handle } from '@sveltejs/kit';
import { clearSession, getAccessToken, getRefreshToken, setSession } from '$lib/server/session';
import { getMe, refreshSession } from '$lib/server/api';

export const handle: Handle = async ({ event, resolve }) => {
	const refresh = getRefreshToken(event.cookies);

	let token = getAccessToken(event.cookies);
	let user = token ? await getMe(token) : null;

	// Access token missing/expired but a refresh token exists → rotate transparently
	// so the session survives the 15m access TTL up to the 30d refresh TTL.
	if (!user && refresh) {
		const next = await refreshSession(refresh);
		if (next) {
			setSession(event.cookies, next);
			token = next.token;
			user = await getMe(token);
		} else {
			clearSession(event.cookies); // refresh invalid/expired → full sign-out
			token = null;
		}
	}

	event.locals.session = token;
	event.locals.user = user;
	return resolve(event);
};
