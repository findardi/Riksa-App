import { browser } from '$app/environment';
import { id, type Dict } from './id';
import { en } from './en';
import { localeState } from './locale.svelte';
import { defaultLocale, type Locale } from './shared';

export type { Locale } from './shared';
export { LOCALES, LOCALE_COOKIE, defaultLocale, localeLabels, isLocale } from './shared';

const locales: Record<Locale, Record<TKey, string>> = { id, en };

export type TKey = keyof Dict;

// Server locale source, registered by hooks.server.ts. Kept behind a setter so
// index.ts stays isomorphic — it never statically imports the node-only ./server.
let resolveServerLocale: () => Locale = () => defaultLocale;
export function setServerLocaleSource(fn: () => Locale): void {
	resolveServerLocale = fn;
}

/** The active locale: reactive client store in the browser, request-scoped on the server. */
function activeLocale(): Locale {
	return browser ? localeState.current : resolveServerLocale();
}

/** Translate a key, interpolating {placeholders} from `vars`. */
export function t(key: TKey, vars?: Record<string, string | number>): string {
	const dict = locales[activeLocale()] ?? locales[defaultLocale];
	let str: string = dict[key] ?? id[key] ?? key;
	if (vars) {
		for (const [k, v] of Object.entries(vars)) {
			str = str.replace(`{${k}}`, String(v));
		}
	}
	return str;
}
