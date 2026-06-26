import type { Handle } from '@sveltejs/kit';
import { clearSession, getAccessToken, getRefreshToken, setSession } from '$lib/server/session';
import { getMe, refreshSession } from '$lib/server/api';
import { refreshSinleFlight } from '$lib/server/refresh-lock';
import { setServerLocaleSource, LOCALE_COOKIE, defaultLocale, isLocale } from '$lib/i18n';
import { getServerLocale, runWithLocale } from '$lib/i18n/server';

// Point t()'s server-side locale resolution at the request-scoped store.
setServerLocaleSource(getServerLocale);

export const handle: Handle = async ({ event, resolve }) => {
	const cookieLocale = event.cookies.get(LOCALE_COOKIE);
	const locale = isLocale(cookieLocale) ? cookieLocale : defaultLocale;

	const refresh = getRefreshToken(event.cookies);

	let token = getAccessToken(event.cookies);
	let user = token ? await getMe(token) : null;

	// Access token missing/expired but a refresh token exists → rotate transparently
	// so the session survives the 15m access TTL up to the 30d refresh TTL.
	if (!user && refresh) {
		const next = await refreshSinleFlight(refresh);
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

	return runWithLocale(locale, () =>
		resolve(event, { transformPageChunk: ({ html }) => html.replace('%lang%', locale) })
	);
};
