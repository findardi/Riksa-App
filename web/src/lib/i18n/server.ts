// Request-scoped locale for server code. A plain module variable would leak one
// request's locale into another under concurrency, so we bind it to the request
// via AsyncLocalStorage. Wired up in hooks.server.ts. Server-only — never import
// from client code (pulls in node:async_hooks).
import { AsyncLocalStorage } from 'node:async_hooks';
import { defaultLocale, type Locale } from './shared';

const storage = new AsyncLocalStorage<Locale>();

/** Run `fn` with `locale` bound for the duration of this request. */
export function runWithLocale<T>(locale: Locale, fn: () => T): T {
	return storage.run(locale, fn);
}

/** The locale bound to the current request, or the default outside one. */
export function getServerLocale(): Locale {
	return storage.getStore() ?? defaultLocale;
}
