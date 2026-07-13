import { invalidateAll } from '$app/navigation';
import { t } from '$lib/i18n';

export type UploadStatus = 'pending' | 'uploading' | 'done' | 'error' | 'canceled';

export interface UploadItem {
	id: string;
	name: string;
	size: number;
	progress: number;
	status: UploadStatus;
	message: string | null;
	workspaceId: string;
	folderId: string;
	folderName: string;
}

const MAX_CONCURRENT = 3;

let items = $state<UploadItem[]>([]);
let panelOpen = $state(true);

// Plain Maps, not SvelteMap: their values are File and XMLHttpRequest handles.
// A reactive proxy around either breaks `xhr.send(file)`, and nothing renders
// from them — the reactive view of an upload lives in `items`.
/* eslint-disable svelte/prefer-svelte-reactivity */
const pendingFiles = new Map<string, File>();
const requests = new Map<string, XMLHttpRequest>();
/* eslint-enable svelte/prefer-svelte-reactivity */
let running = 0;

function find(id: string): UploadItem | undefined {
	return items.find((i) => i.id === id);
}

async function messageOf(res: Response): Promise<string> {
	const body = (await res.json().catch(() => null)) as { message?: string } | null;
	return body?.message || t('err.generic');
}

function put(id: string, url: string, file: File): Promise<void> {
	return new Promise((resolve, reject) => {
		const xhr = new XMLHttpRequest();
		requests.set(id, xhr);

		xhr.open('PUT', url, true);
		xhr.setRequestHeader('Content-Type', file.type || 'application/octet-stream');

		xhr.upload.onprogress = (e) => {
			if (!e.lengthComputable) return;
			const item = find(id);
			if (item) item.progress = Math.round((e.loaded / e.total) * 100);
		};
		xhr.onload = () =>
			xhr.status >= 200 && xhr.status < 300
				? resolve()
				: reject(new Error(t('doc.upload.err.storage')));
		xhr.onerror = () => reject(new Error(t('err.network')));
		xhr.onabort = () => reject(new Error(t('doc.upload.status.canceled')));

		xhr.send(file);
	});
}

async function run(id: string): Promise<void> {
	const item = find(id);
	const file = pendingFiles.get(id);

	if (!item || !file) {
		running--;
		return;
	}

	item.status = 'uploading';
	item.progress = 0;
	item.message = null;

	try {
		const presign = await fetch('/api/content/upload-url', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ workspaceId: item.workspaceId, folderId: item.folderId })
		});
		if (!presign.ok) throw new Error(await messageOf(presign));
		const { upload_url, storage_key } = (await presign.json()) as {
			upload_url: string;
			storage_key: string;
		};

		await put(id, upload_url, file);

		const complete = await fetch('/api/content/documents', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({
				workspaceId: item.workspaceId,
				folderId: item.folderId,
				name: item.name,
				storageKey: storage_key
			})
		});
		if (!complete.ok) throw new Error(await messageOf(complete));

		item.status = 'done';
		item.progress = 100;
		pendingFiles.delete(id);
	} catch (e) {
		// Re-read: `cancel()` may have flipped the status while we were awaiting.
		const latest = find(id);
		if (latest && latest.status !== 'canceled') {
			latest.status = 'error';
			latest.message = e instanceof Error ? e.message : t('err.generic');
		}
	} finally {
		requests.delete(id);
		running--;
		pump();
		settle();
	}
}

function pump(): void {
	while (running < MAX_CONCURRENT) {
		const next = items.find((i) => i.status === 'pending');
		if (!next) return;
		running++;
		void run(next.id);
	}
}

function settle(): void {
	const busy = items.some((i) => i.status === 'pending' || i.status === 'uploading');
	if (busy) return;
	if (items.some((i) => i.status === 'done')) void invalidateAll();
}

function drop(id: string): void {
	pendingFiles.delete(id);
	requests.delete(id);
	items = items.filter((i) => i.id !== id);
}

export const uploadQueue = {
	get items(): UploadItem[] {
		return items;
	},
	get busy(): number {
		return items.filter((i) => i.status === 'pending' || i.status === 'uploading').length;
	},
	get failed(): number {
		return items.filter((i) => i.status === 'error').length;
	},
	get open(): boolean {
		return panelOpen;
	},
	set open(v: boolean) {
		panelOpen = v;
	},

	enqueue(workspaceId: string, folderId: string, folderName: string, files: File[]): void {
		if (!files.length) return;

		for (const file of files) {
			const id = crypto.randomUUID();
			pendingFiles.set(id, file);
			items.push({
				id,
				name: file.name,
				size: file.size,
				progress: 0,
				status: 'pending',
				message: null,
				workspaceId,
				folderId,
				folderName
			});
		}

		panelOpen = true;
		pump();
	},

	cancel(id: string): void {
		const item = find(id);
		if (!item) return;

		const wasUploading = item.status === 'uploading';
		item.status = 'canceled';
		item.message = null;
		if (wasUploading) requests.get(id)?.abort();
		else settle();
	},

	retry(id: string): void {
		const item = find(id);
		if (!item || !pendingFiles.has(id)) return;

		item.status = 'pending';
		item.progress = 0;
		item.message = null;
		pump();
	},

	remove(id: string): void {
		const item = find(id);
		if (item?.status === 'uploading') requests.get(id)?.abort();
		drop(id);
		settle();
	},

	clearFinished(): void {
		for (const item of items) {
			if (item.status !== 'pending' && item.status !== 'uploading') pendingFiles.delete(item.id);
		}
		items = items.filter((i) => i.status === 'pending' || i.status === 'uploading');
	}
};
