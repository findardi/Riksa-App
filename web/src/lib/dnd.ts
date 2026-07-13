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
