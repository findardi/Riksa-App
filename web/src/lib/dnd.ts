// Drag payload types. A drag carrying `Files` is an upload from the OS; these
// two mark drags that started inside the app, and `dataTransfer.types` is the
// only part of a drag readable during `dragover` — so they are the switch.
export const FOLDER_MIME = 'application/x-riksa-folder';
export const DOCUMENT_MIME = 'application/x-riksa-document';

// `items` carries directory entries too; only real files survive. Falls back to
// `files` for browsers that leave `items` empty on drop.
export function filesFrom(dt: DataTransfer | null): File[] {
	if (!dt) return [];
	if (dt.items?.length) {
		const out: File[] = [];
		for (const item of Array.from(dt.items)) {
			if (item.kind !== 'file') continue;
			if (item.webkitGetAsEntry?.()?.isFile === false) continue;
			const f = item.getAsFile();
			if (f) out.push(f);
		}
		return out;
	}
	return Array.from(dt.files);
}

const UUID_RE = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;

export const isUuid = (v: string) => UUID_RE.test(v);

// A reorder sends `position` through a form field, so it arrives as a string.
// `null` means "no position given" — the server then appends. Anything that is
// not a non-negative integer is treated as absent rather than sent as garbage.
export function parsePosition(raw: FormDataEntryValue | null): number | null {
	if (raw === null) return null;
	const s = raw.toString().trim();
	if (!s) return null;
	const n = Number(s);
	return Number.isInteger(n) && n >= 0 ? n : null;
}
