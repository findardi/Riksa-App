// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			/** Raw access token from the httpOnly session cookie, or null. */
			session: string | null;
			/** Current user resolved from GET /auth/me, or null if unauthenticated. */
			user: import('$lib/types').MeData | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
