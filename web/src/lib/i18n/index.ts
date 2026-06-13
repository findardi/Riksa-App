import { id, type Dict } from './id';

// Single-locale seam (ID default). To add English: create en.ts with the same
// keys, then `const locales = { id, en }` and thread a locale through `t`.
const locales = { id };
export const defaultLocale: keyof typeof locales = 'id';

export type TKey = keyof Dict;

/** Translate a key, interpolating {placeholders} from `vars`. */
export function t(key: TKey, vars?: Record<string, string | number>): string {
	let str: string = locales[defaultLocale][key] ?? key;
	if (vars) {
		for (const [k, v] of Object.entries(vars)) {
			str = str.replace(`{${k}}`, String(v));
		}
	}
	return str;
}
