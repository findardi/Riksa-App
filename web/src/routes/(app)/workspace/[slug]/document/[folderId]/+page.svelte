<script lang="ts">
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button, showToast } from '$lib/components/common';
	import { DOCUMENT_MIME } from '$lib/dnd';
	import { formatBytes, formatDate } from '$lib/format';
	import { t } from '$lib/i18n';
	import type { DocumentData, FolderTreeNode } from '$lib/types/content';
	import type { MyAccessWorkspace } from '$lib/types/workspace';
	import { uploadQueue } from '$lib/upload/queue.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const documents = $derived(data.documents);
	const folders = $derived(data.folders);
	const workspace = $derived(data.workspace);
	const folderId = $derived(page.params.folderId!);

	const perms = $derived((page.data as { access?: MyAccessWorkspace }).access?.permissions ?? []);
	const canUpload = $derived(perms.includes('document:upload'));
	const canDownload = $derived(perms.includes('document:download'));
	const canDelete = $derived(perms.includes('document:delete'));
	const canEditDoc = $derived(perms.includes('document:edit'));

	function findNode(nodes: FolderTreeNode[], id: string): FolderTreeNode | null {
		for (const n of nodes) {
			if (n.id === id) return n;
			const hit = findNode(n.children ?? [], id);
			if (hit) return hit;
		}
		return null;
	}
	const folder = $derived(findNode(folders, folderId));

	let fileInput = $state<HTMLInputElement>();
	let paneTargeted = $state(false);

	function filesFrom(dt: DataTransfer | null): File[] {
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

	function enqueue(files: File[]) {
		if (!canUpload || !files.length || !folder) return;
		uploadQueue.enqueue(workspace.id, folder.id, folder.name, files);
	}

	function paneDragOver(e: DragEvent) {
		if (!canUpload || !e.dataTransfer?.types.includes('Files')) return;
		e.preventDefault();
		e.dataTransfer.dropEffect = 'copy';
		paneTargeted = true;
	}

	function paneDragLeave(e: DragEvent) {
		const next = e.relatedTarget as Node | null;
		if (next && (e.currentTarget as Element).contains(next)) return;
		paneTargeted = false;
	}

	function paneDrop(e: DragEvent) {
		if (!canUpload || !e.dataTransfer?.types.includes('Files')) return;
		e.preventDefault();
		paneTargeted = false;
		enqueue(filesFrom(e.dataTransfer));
	}

	function onPick(e: Event) {
		const input = e.currentTarget as HTMLInputElement;
		enqueue(Array.from(input.files ?? []));
		input.value = '';
	}

	// --- drag out (the tree in the layout owns the drop targets) ---
	let draggingDocId = $state<string | null>(null);

	function docDragStart(e: DragEvent, doc: DocumentData) {
		if (!canEditDoc) {
			e.preventDefault();
			return;
		}
		draggingDocId = doc.id;
		e.dataTransfer?.setData(DOCUMENT_MIME, doc.id);
		if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move';
	}

	// --- move dialog (keyboard + touch path; drag is the shortcut) ---
	type Option = { id: string; number: string; name: string; depth: number };
	const moveOptions = $derived.by(() => {
		const out: Option[] = [];
		const build = (nodes: FolderTreeNode[], depth: number) => {
			for (const n of nodes) {
				if (n.id !== folderId) out.push({ id: n.id, number: n.number, name: n.name, depth });
				if (n.children?.length) build(n.children, depth + 1);
			}
		};
		build(folders, 0);
		return out;
	});

	let moveDialog = $state<HTMLDialogElement>();
	let movingDoc = $state<DocumentData | null>(null);
	let moveTarget = $state('');
	let moveError = $state<string | null>(null);
	let moveSubmitting = $state(false);

	function openMove(doc: DocumentData) {
		movingDoc = doc;
		moveTarget = moveOptions[0]?.id ?? '';
		moveError = null;
		moveDialog?.showModal();
	}

	const submitMove: SubmitFunction = () => {
		moveSubmitting = true;
		return async ({ result }) => {
			moveSubmitting = false;
			if (result.type === 'success') {
				moveDialog?.close();
				await invalidateAll();
				showToast(t('doc.docs.moved'), 'success');
			} else if (result.type === 'failure') {
				moveError = (result.data?.message as string) ?? t('err.generic');
			} else {
				moveError = t('err.generic');
			}
		};
	};

	const kindOf = (mime: string) =>
		mime.startsWith('image/')
			? 'image'
			: mime.includes('spreadsheet') || mime.includes('excel') || mime.includes('csv')
				? 'sheet'
				: 'file';

	// --- delete dialog ---
	let deleteDialog = $state<HTMLDialogElement>();
	let deleting = $state<DocumentData | null>(null);
	let deleteError = $state<string | null>(null);
	let deleteSubmitting = $state(false);

	function openDelete(doc: DocumentData) {
		deleting = doc;
		deleteError = null;
		deleteDialog?.showModal();
	}

	const submitDelete: SubmitFunction = () => {
		deleteSubmitting = true;
		return async ({ result }) => {
			deleteSubmitting = false;
			if (result.type === 'success') {
				deleteDialog?.close();
				await invalidateAll();
				showToast(t('doc.docs.deleted'), 'success');
			} else if (result.type === 'failure') {
				deleteError = (result.data?.message as string) ?? t('err.generic');
			} else {
				deleteError = t('err.generic');
			}
		};
	};
</script>

{#snippet fileIcon(kind: string)}
	<svg
		class="h-4 w-4 flex-none text-muted"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="1.6"
		stroke-linecap="round"
		stroke-linejoin="round"
		aria-hidden="true"
	>
		{#if kind === 'image'}
			<rect x="3" y="4" width="18" height="16" rx="2" />
			<circle cx="8.5" cy="9.5" r="1.5" />
			<path d="m21 16-5-5L6 20" />
		{:else if kind === 'sheet'}
			<rect x="3" y="4" width="18" height="16" rx="2" />
			<path d="M3 10h18M9 10v10M15 10v10" />
		{:else}
			<path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z" />
			<path d="M14 3v5h5" />
		{/if}
	</svg>
{/snippet}

<section
	ondragover={paneDragOver}
	ondragleave={paneDragLeave}
	ondrop={paneDrop}
	aria-label={folder?.name ?? t('doc.title')}
	class="min-h-64 rounded-box border transition-colors
		{paneTargeted ? 'border-primary/50 bg-primary/[0.04]' : 'border-base-content/10 bg-base-100'}"
>
	<header
		class="flex flex-wrap items-center justify-between gap-3 border-b border-base-content/8 px-4 py-3"
	>
		<div class="flex min-w-0 items-baseline gap-2">
			{#if folder}
				<span class="font-mono text-xs tabular-nums text-muted">{folder.number}</span>
			{/if}
			<h2 class="min-w-0 truncate text-[0.9375rem] font-semibold tracking-[-0.01em]">
				{folder?.name ?? t('doc.docs.unknownFolder')}
			</h2>
			<span class="flex-none font-mono text-xs text-muted">
				{t('doc.docs.count', { n: documents.length })}
			</span>
		</div>

		{#if canUpload}
			<div>
				<input
					bind:this={fileInput}
					onchange={onPick}
					type="file"
					multiple
					class="sr-only"
					aria-label={t('doc.docs.upload')}
				/>
				<Button onclick={() => fileInput?.click()}>
					<svg
						class="h-4 w-4"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.8"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<path d="M12 19V5M5 12l7-7 7 7" />
					</svg>
					{t('doc.docs.upload')}
				</Button>
			</div>
		{/if}
	</header>

	{#if documents.length === 0}
		<div class="flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
			<svg
				class="h-9 w-9 text-muted/70"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.4"
				stroke-linecap="round"
				stroke-linejoin="round"
				aria-hidden="true"
			>
				<path d="M12 17V9M8.5 12.5 12 9l3.5 3.5" />
				<path d="M20 16.5V8a2 2 0 0 0-2-2h-6L10 4H6a2 2 0 0 0-2 2v10.5" />
				<path d="M4 16.5A1.5 1.5 0 0 0 5.5 18h13a1.5 1.5 0 0 0 1.5-1.5" />
			</svg>
			<div>
				<p class="text-[0.9375rem] font-medium">{t('doc.docs.empty.title')}</p>
				<p class="mx-auto mt-1 max-w-sm text-sm text-muted text-pretty">
					{canUpload ? t('doc.docs.empty.body') : t('doc.docs.empty.readonly')}
				</p>
			</div>
			{#if canUpload}
				<Button onclick={() => fileInput?.click()}>{t('doc.docs.empty.cta')}</Button>
			{/if}
		</div>
	{:else}
		<ul class="divide-y divide-base-content/6">
			{#each documents as doc (doc.id)}
				<li
					draggable={canEditDoc}
					ondragstart={(e) => docDragStart(e, doc)}
					ondragend={() => (draggingDocId = null)}
					class="group flex items-center gap-2.5 px-4 py-2.5 transition-colors hover:bg-base-content/[0.025]
						{draggingDocId === doc.id ? 'opacity-40' : ''}"
				>
					{@render fileIcon(kindOf(doc.mime))}

					<span class="min-w-0 flex-1 truncate text-sm" title={doc.name}>{doc.name}</span>

					<span
						class="hidden flex-none rounded-selector bg-base-content/5 px-1.5 py-0.5 font-mono text-[0.6875rem] text-muted sm:inline"
						title={t('doc.docs.versionTitle', { n: doc.version_no })}
					>
						v{doc.version_no}
					</span>

					<span
						class="hidden w-20 flex-none text-right font-mono text-xs text-muted tabular-nums md:inline"
					>
						{formatBytes(doc.size)}
					</span>

					<span
						class="hidden w-24 flex-none text-right font-mono text-xs text-muted tabular-nums lg:inline"
					>
						{formatDate(doc.updated_at)}
					</span>

					<div
						class="flex flex-none items-center gap-0.5 opacity-0 transition-opacity focus-within:opacity-100 group-hover:opacity-100"
					>
						{#if canDownload}
							<a
								href={resolve(
									`/api/content/download?workspaceId=${workspace.id}&documentId=${doc.id}`
								)}
								draggable="false"
								title={t('doc.docs.download')}
								aria-label={t('doc.docs.downloadOf', { name: doc.name })}
								class="grid h-7 w-7 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content"
							>
								<svg
									class="h-4 w-4"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="1.8"
									stroke-linecap="round"
									stroke-linejoin="round"
									aria-hidden="true"
								>
									<path d="M12 4v11M7.5 10.5 12 15l4.5-4.5" />
									<path d="M5 19h14" />
								</svg>
							</a>
						{/if}
						{#if canEditDoc && moveOptions.length > 0}
							<button
								type="button"
								onclick={() => openMove(doc)}
								title={t('doc.docs.move')}
								aria-label={t('doc.docs.moveOf', { name: doc.name })}
								class="grid h-7 w-7 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content"
							>
								<svg
									class="h-4 w-4"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="1.8"
									stroke-linecap="round"
									stroke-linejoin="round"
									aria-hidden="true"
								>
									<path d="M5 9V5h4M19 15v4h-4M5 5l6 6M19 19l-6-6" />
								</svg>
							</button>
						{/if}
						{#if canDelete}
							<button
								type="button"
								onclick={() => openDelete(doc)}
								title={t('doc.docs.delete')}
								aria-label={t('doc.docs.deleteOf', { name: doc.name })}
								class="grid h-7 w-7 place-items-center rounded-field text-muted transition-colors hover:bg-error/10 hover:text-error"
							>
								<svg
									class="h-4 w-4"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="1.8"
									stroke-linecap="round"
									stroke-linejoin="round"
									aria-hidden="true"
								>
									<path
										d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2m2 0v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"
									/>
									<path d="M10 11v6M14 11v6" />
								</svg>
							</button>
						{/if}
					</div>
				</li>
			{/each}
		</ul>

		{#if canUpload}
			<p class="px-4 py-3 text-xs text-muted">{t('doc.docs.dropHint')}</p>
		{/if}
	{/if}
</section>

<dialog bind:this={moveDialog} class="modal" aria-labelledby="doc-move-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="doc-move-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('doc.docs.move.title')}
		</h2>
		{#if movingDoc}
			<p class="mt-1 text-sm text-muted text-pretty">
				{t('doc.docs.move.desc', { name: movingDoc.name })}
			</p>
		{/if}

		{#if moveError}
			<div class="mt-4"><Alert align="start">{moveError}</Alert></div>
		{/if}

		<form method="POST" action="?/moveDocument" use:enhance={submitMove} class="mt-5">
			<input type="hidden" name="documentId" value={movingDoc?.id ?? ''} />
			<label for="doc-move-dest" class="mb-1.5 block text-sm font-medium">
				{t('doc.docs.move.dest')}
			</label>
			<select
				id="doc-move-dest"
				name="folderId"
				bind:value={moveTarget}
				class="select select-sm w-full font-mono"
			>
				{#each moveOptions as opt (opt.id)}
					<option value={opt.id}>{'  '.repeat(opt.depth)}{opt.number} {opt.name}</option>
				{/each}
			</select>

			<div class="mt-6 flex justify-end gap-2">
				<Button type="button" variant="ghost" onclick={() => moveDialog?.close()}>
					{t('doc.cancel')}
				</Button>
				<Button type="submit" loading={moveSubmitting}>
					{moveSubmitting ? t('doc.docs.move.submitting') : t('doc.docs.move.submit')}
				</Button>
			</div>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('doc.cancel')}></button>
	</form>
</dialog>

<dialog bind:this={deleteDialog} class="modal" aria-labelledby="doc-delete-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="doc-delete-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('doc.docs.delete.title')}
		</h2>
		{#if deleting}
			<p class="mt-1 text-sm text-muted text-pretty">
				{t('doc.docs.delete.warning', { name: deleting.name })}
			</p>
			{#if deleting.version_no > 1}
				<p class="mt-2 text-sm font-medium text-error text-pretty">
					{t('doc.docs.delete.versions', { n: deleting.version_no })}
				</p>
			{/if}
		{/if}

		{#if deleteError}
			<div class="mt-4"><Alert align="start">{deleteError}</Alert></div>
		{/if}

		<form
			method="POST"
			action="?/deleteDocument"
			use:enhance={submitDelete}
			class="mt-6 flex justify-end gap-2"
		>
			<input type="hidden" name="documentId" value={deleting?.id ?? ''} />
			<Button type="button" variant="ghost" onclick={() => deleteDialog?.close()}>
				{t('doc.cancel')}
			</Button>
			<Button type="submit" variant="danger" loading={deleteSubmitting}>
				{deleteSubmitting ? t('doc.docs.delete.submitting') : t('doc.docs.delete.submit')}
			</Button>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('doc.cancel')}></button>
	</form>
</dialog>
