import { browser } from '$app/environment';
import { invalidateAll } from '$app/navigation';
import { defaultLocale, isLocale, LOCALE_COOKIE, type Locale } from './shared';

function readCookieLocale(): Locale {
	if (!browser) return defaultLocale;
	const match = document.cookie.match(new RegExp(`(?:^|;\\s*)${LOCALE_COOKIE}=([^;]+)`));
	return match && isLocale(match[1]) ? match[1] : defaultLocale;
}

export const localeState = $state<{ current: Locale }>({ current: readCookieLocale() });

const ONE_YEAR = 60 * 60 * 24 * 365;

export async function setLocale(next: Locale): Promise<void> {
	if (next === localeState.current) return;
	localeState.current = next;
	if (!browser) return;
	document.cookie = `${LOCALE_COOKIE}=${next}; path=/; max-age=${ONE_YEAR}; samesite=lax`;
	document.documentElement.lang = next;
	await invalidateAll();
}
