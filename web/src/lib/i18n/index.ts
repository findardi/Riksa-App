import { browser } from '$app/environment';
import { id, type Dict } from './id';
import { en } from './en';
import { localeState } from './locale.svelte';
import { defaultLocale, type Locale } from './shared';

export type { Locale } from './shared';
export { LOCALES, LOCALE_COOKIE, defaultLocale, localeLabels, isLocale } from './shared';

const locales: Record<Locale, Record<TKey, string>> = { id, en };

export type TKey = keyof Dict;

let resolveServerLocale: () => Locale = () => defaultLocale;
export function setServerLocaleSource(fn: () => Locale): void {
	resolveServerLocale = fn;
}

function activeLocale(): Locale {
	return browser ? localeState.current : resolveServerLocale();
}

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
