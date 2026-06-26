import { AsyncLocalStorage } from 'node:async_hooks';
import { defaultLocale, type Locale } from './shared';

const storage = new AsyncLocalStorage<Locale>();

export function runWithLocale<T>(locale: Locale, fn: () => T): T {
	return storage.run(locale, fn);
}

export function getServerLocale(): Locale {
	return storage.getStore() ?? defaultLocale;
}
