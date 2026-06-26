// Client-side active locale. A module-level rune is safe in the browser (one JS
// context per user) and lets `t()` re-render reactively when the locale flips.
// The server never reads this — it uses the request-scoped seam in ./server.
import { browser } from '$app/environment';
import { invalidateAll } from '$app/navigation';
import { defaultLocale, isLocale, LOCALE_COOKIE, type Locale } from './shared';

function readCookieLocale(): Locale {
	if (!browser) return defaultLocale;
	const match = document.cookie.match(new RegExp(`(?:^|;\\s*)${LOCALE_COOKIE}=([^;]+)`));
	return match && isLocale(match[1]) ? match[1] : defaultLocale;
}

// Initialised from the cookie so client hydration matches what the server rendered.
export const localeState = $state<{ current: Locale }>({ current: readCookieLocale() });

const ONE_YEAR = 60 * 60 * 24 * 365;

/**
 * Switch the active locale: update the reactive store (instant client re-render),
 * persist a cookie so the server picks it up, then re-run load functions so
 * server-rendered strings (validation/API errors) refresh too.
 */
export async function setLocale(next: Locale): Promise<void> {
	if (next === localeState.current) return;
	localeState.current = next;
	if (!browser) return;
	document.cookie = `${LOCALE_COOKIE}=${next}; path=/; max-age=${ONE_YEAR}; samesite=lax`;
	document.documentElement.lang = next;
	await invalidateAll();
}
