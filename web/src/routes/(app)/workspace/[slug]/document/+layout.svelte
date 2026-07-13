<script lang="ts">
	import { tick } from 'svelte';
	import { applyAction, deserialize, enhance } from '$app/forms';
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { ActionResult, SubmitFunction } from '@sveltejs/kit';
	import { UploadQueue } from '$lib/components/app';
	import { Alert, Button, Field, Toaster, showToast } from '$lib/components/common';
	import { DOCUMENT_MIME, FOLDER_MIME, filesFrom } from '$lib/dnd';
	import { t } from '$lib/i18n';
	import { findNode } from '$lib/tree';
	import type { FolderTreeNode } from '$lib/types/content';
	import type { MyAccessWorkspace } from '$lib/types/workspace';
	import { uploadQueue } from '$lib/upload/queue.svelte';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();
	const folders = $derived(data.folders);
	const workspace = $derived(data.workspace);
	const noAccess = $derived(data.noAccess);

	const ROOT = '';

	const slug = $derived(page.params.slug!);
	const actionBase = $derived(resolve('/(app)/workspace/[slug]/document', { slug }));
	const docHref = (folderId: string) =>
		resolve('/(app)/workspace/[slug]/document/[folderId]', { slug, folderId });
	const accessHref = (folderId: string) =>
		resolve('/(app)/workspace/[slug]/document/[folderId]/access', { slug, folderId });
	const activeId = $derived(page.params.folderId ?? null);

	const onAccess = $derived(page.url.pathname.endsWith('/access'));
	const rowHref = (folderId: string) => (onAccess ? accessHref(folderId) : docHref(folderId));

	const perms = $derived((page.data as { access?: MyAccessWorkspace }).access?.permissions ?? []);
	const canCreate = $derived(perms.includes('folder:create'));
	const canEdit = $derived(perms.includes('folder:edit'));
	const canDelete = $derived(perms.includes('folder:delete'));
	const canUpload = $derived(perms.includes('document:upload'));
	const canEditDoc = $derived(perms.includes('document:edit'));
	const canAssign = $derived(perms.includes('group:assign'));
	const canAct = $derived(canCreate || canEdit || canDelete || canAssign);

	const defaultFolder = $derived(folders.find((f) => f.is_default) ?? null);
	const activeFolder = $derived(activeId ? findNode(folders, activeId) : null);
	const fallbackFolder = $derived(activeFolder ?? defaultFolder);

	// Top-level folders reveal their children; deeper levels start closed, so a
	// 300-folder index opens as a readable table of contents rather than a wall.
	const DEFAULT_OPEN_DEPTH = 1;

	let expanded = $state<Record<string, boolean>>({});
	const isExpanded = (id: string, depth: number) => expanded[id] ?? depth < DEFAULT_OPEN_DEPTH;
	function toggle(id: string, depth: number) {
		expanded[id] = !isExpanded(id, depth);
	}

	let query = $state('');
	const q = $derived(query.trim().toLowerCase());

	// `null` = no search. Otherwise the ids to render: every match plus the
	// ancestors needed to reach it, so a hit never appears without its context.
	const matched = $derived.by(() => {
		if (!q) return null;
		const keep: Record<string, true> = {};
		const visit = (nodes: FolderTreeNode[], trail: string[]): void => {
			for (const n of nodes) {
				const hit = n.name.toLowerCase().includes(q) || n.number.toLowerCase().includes(q);
				if (hit) {
					keep[n.id] = true;
					for (const id of trail) keep[id] = true;
				}
				visit(n.children ?? [], [...trail, n.id]);
			}
		};
		visit(folders, []);
		return keep;
	});

	type Row = { node: FolderTreeNode; depth: number; hasChildren: boolean; open: boolean };
	function walk(nodes: FolderTreeNode[], depth: number, out: Row[]) {
		for (const n of nodes) {
			if (matched && !matched[n.id]) continue;
			const kids = (n.children ?? []).filter((k) => !matched || matched[k.id]);
			// A search result is always drilled open; otherwise honour the toggle.
			const open = matched ? true : isExpanded(n.id, depth);
			out.push({ node: n, depth, hasChildren: kids.length > 0, open });
			if (kids.length && open) walk(kids, depth + 1, out);
		}
	}
	const rows = $derived.by(() => {
		const out: Row[] = [];
		walk(folders, 0, out);
		return out;
	});

	const depthOf = $derived.by(() => {
		const map: Record<string, number> = {};
		const build = (nodes: FolderTreeNode[], depth: number) => {
			for (const n of nodes) {
				map[n.id] = depth;
				if (n.children?.length) build(n.children, depth + 1);
			}
		};
		build(folders, 0);
		return map;
	});

	const branchIds = $derived.by(() => {
		const out: string[] = [];
		const visit = (nodes: FolderTreeNode[]) => {
			for (const n of nodes) {
				if (n.children?.length) out.push(n.id);
				visit(n.children ?? []);
			}
		};
		visit(folders);
		return out;
	});

	const anyCollapsed = $derived(branchIds.some((id) => !isExpanded(id, depthOf[id] ?? 0)));

	function setAll(open: boolean) {
		const next: Record<string, boolean> = {};
		for (const id of branchIds) next[id] = open;
		expanded = next;
	}

	const totalCount = $derived.by(() => {
		let n = 0;
		const stack = [...folders];
		while (stack.length) {
			const f = stack.pop()!;
			n++;
			if (f.children?.length) stack.push(...f.children);
		}
		return n;
	});

	const parentOf = $derived.by(() => {
		const map: Record<string, string> = {};
		const build = (nodes: FolderTreeNode[], parent: string) => {
			for (const n of nodes) {
				map[n.id] = parent;
				if (n.children?.length) build(n.children, n.id);
			}
		};
		build(folders, ROOT);
		return map;
	});

	function descendantCount(node: FolderTreeNode): number {
		let n = 0;
		for (const c of node.children ?? []) n += 1 + descendantCount(c);
		return n;
	}

	function subtreeIds(node: FolderTreeNode, into: string[] = []): string[] {
		into.push(node.id);
		for (const c of node.children ?? []) subtreeIds(c, into);
		return into;
	}

	const indent = (depth: number) => `padding-left: ${depth * 1.375 + 0.25}rem`;

	// --- drag & drop -------------------------------------------------------
	// `Files` is an upload from the OS, FOLDER_MIME moves a folder, DOCUMENT_MIME
	// moves a document. `types` is the only payload readable during dragover, so
	// it is the switch.

	let draggingId = $state<string | null>(null);
	let dropTarget = $state<string | null>(null);
	let fileDragging = $state(false);
	let dragDepth = 0;

	const hasFiles = (e: DragEvent) => !!e.dataTransfer?.types.includes('Files');
	const hasFolder = (e: DragEvent) => !!e.dataTransfer?.types.includes(FOLDER_MIME);
	const hasDocument = (e: DragEvent) => !!e.dataTransfer?.types.includes(DOCUMENT_MIME);

	// A document is only ever dragged out of the folder currently on screen, so
	// the folder it came from is `activeId` — no payload lookup needed to reject
	// a drop back onto its own folder.
	const canMoveDocInto = (targetId: string) => canEditDoc && !!activeId && targetId !== activeId;

	const dropLabel = $derived.by(() => {
		if (dropTarget === ROOT) return defaultFolder?.name ?? null;
		if (dropTarget) return findNode(folders, dropTarget)?.name ?? null;
		return fallbackFolder?.name ?? null;
	});

	function resetDrag() {
		dragDepth = 0;
		fileDragging = false;
		dropTarget = null;
		draggingId = null;
	}

	function uploadTo(folderId: string, folderName: string, files: File[]) {
		if (!canUpload || !files.length) return;
		uploadQueue.enqueue(workspace.id, folderId, folderName, files);
	}

	function uploadToDefault(files: File[]) {
		if (!files.length) return;
		if (!defaultFolder) {
			showToast(t('doc.upload.noDefault'), 'error');
			return;
		}
		uploadTo(defaultFolder.id, defaultFolder.name, files);
	}

	// A loose drop belongs to whatever folder the user is looking at; only the
	// index itself (no folder open) falls back to the default folder.
	function uploadToFallback(files: File[]) {
		if (!files.length) return;
		if (!fallbackFolder) {
			showToast(t('doc.upload.noDefault'), 'error');
			return;
		}
		uploadTo(fallbackFolder.id, fallbackFolder.name, files);
	}

	function canMoveInto(targetId: string): boolean {
		if (!draggingId || !canEdit) return false;
		if (targetId === draggingId) return false;
		if (parentOf[draggingId] === targetId) return false;
		const node = findNode(folders, draggingId);
		return !!node && !subtreeIds(node).includes(targetId);
	}

	async function moveTo(folderId: string, parentId: string) {
		const body = new FormData();
		body.set('folderId', folderId);
		body.set('parentId', parentId);

		const res = await fetch(`${actionBase}?/move`, {
			method: 'POST',
			body,
			headers: { 'x-sveltekit-action': 'true' }
		});
		const result: ActionResult = deserialize(await res.text());

		if (result.type === 'success') {
			await invalidateAll();
			showToast(t('doc.moved'), 'success');
		} else if (result.type === 'failure') {
			showToast((result.data?.message as string) ?? t('err.generic'), 'error');
		} else {
			await applyAction(result);
		}
	}

	async function moveDocumentTo(documentId: string, folderId: string) {
		if (!activeId) return;

		const body = new FormData();
		body.set('documentId', documentId);
		body.set('folderId', folderId);

		const res = await fetch(`${docHref(activeId)}?/moveDocument`, {
			method: 'POST',
			body,
			headers: { 'x-sveltekit-action': 'true' }
		});
		const result: ActionResult = deserialize(await res.text());

		if (result.type === 'success') {
			await invalidateAll();
			showToast(t('doc.docs.moved'), 'success');
		} else if (result.type === 'failure') {
			showToast((result.data?.message as string) ?? t('err.generic'), 'error');
		} else {
			await applyAction(result);
		}
	}

	function rowDragStart(e: DragEvent, node: FolderTreeNode) {
		if (!canEdit || node.is_default) {
			e.preventDefault();
			return;
		}
		draggingId = node.id;
		e.dataTransfer?.setData(FOLDER_MIME, node.id);
		if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move';
	}

	// A row claims every drag of a kind it handles, even when the target is
	// illegal — it just refuses it. Letting an illegal drag fall through would
	// hand it to the rail below, silently turning "drop onto my own folder" into
	// "move to root".
	function claim(e: DragEvent) {
		e.preventDefault();
		e.stopPropagation();
	}

	function rowDragOver(e: DragEvent, node: FolderTreeNode) {
		if (hasFiles(e)) {
			if (!canUpload) return;
			claim(e);
			if (e.dataTransfer) e.dataTransfer.dropEffect = 'copy';
			dropTarget = node.id;
		} else if (hasFolder(e) && canEdit && draggingId) {
			claim(e);
			const ok = canMoveInto(node.id);
			if (e.dataTransfer) e.dataTransfer.dropEffect = ok ? 'move' : 'none';
			dropTarget = ok ? node.id : null;
		} else if (hasDocument(e) && canEditDoc && activeId) {
			claim(e);
			const ok = canMoveDocInto(node.id);
			if (e.dataTransfer) e.dataTransfer.dropEffect = ok ? 'move' : 'none';
			dropTarget = ok ? node.id : null;
		}
	}

	function rowDragLeave(e: DragEvent) {
		const next = e.relatedTarget as Node | null;
		if (next && (e.currentTarget as Element).contains(next)) return;
		if (dropTarget !== null) dropTarget = null;
	}

	function rowDrop(e: DragEvent, node: FolderTreeNode) {
		if (hasFiles(e)) {
			if (!canUpload) return;
			claim(e);
			uploadTo(node.id, node.name, filesFrom(e.dataTransfer));
			resetDrag();
		} else if (hasFolder(e) && canEdit && draggingId) {
			claim(e);
			if (canMoveInto(node.id)) void moveTo(draggingId, node.id);
			resetDrag();
		} else if (hasDocument(e) && canEditDoc && activeId) {
			claim(e);
			const documentId = e.dataTransfer!.getData(DOCUMENT_MIME);
			if (documentId && canMoveDocInto(node.id)) void moveDocumentTo(documentId, node.id);
			resetDrag();
		}
	}

	function railDragOver(e: DragEvent) {
		if (hasFiles(e)) {
			if (!canUpload || !defaultFolder) return;
			e.preventDefault();
			if (e.dataTransfer) e.dataTransfer.dropEffect = 'copy';
			dropTarget = ROOT;
		} else if (hasFolder(e) && draggingId && canEdit && parentOf[draggingId] !== ROOT) {
			e.preventDefault();
			if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
			dropTarget = ROOT;
		} else if (hasDocument(e) && defaultFolder && canMoveDocInto(defaultFolder.id)) {
			e.preventDefault();
			if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
			dropTarget = ROOT;
		}
	}

	function railDrop(e: DragEvent) {
		if (hasFiles(e)) {
			if (!canUpload || !defaultFolder) return;
			e.preventDefault();
			uploadToDefault(filesFrom(e.dataTransfer));
			resetDrag();
		} else if (hasFolder(e) && draggingId && canEdit && parentOf[draggingId] !== ROOT) {
			e.preventDefault();
			void moveTo(draggingId, ROOT);
			resetDrag();
		} else if (hasDocument(e) && defaultFolder && canMoveDocInto(defaultFolder.id)) {
			e.preventDefault();
			const documentId = e.dataTransfer!.getData(DOCUMENT_MIME);
			if (documentId) void moveDocumentTo(documentId, defaultFolder.id);
			resetDrag();
		}
	}

	function winDragEnter(e: DragEvent) {
		dragDepth++;
		if (hasFiles(e)) fileDragging = true;
	}

	function winDragLeave() {
		dragDepth = Math.max(0, dragDepth - 1);
		if (dragDepth === 0) fileDragging = false;
	}

	function winDragOver(e: DragEvent) {
		if (!hasFiles(e) || !canUpload) return;
		e.preventDefault();
		if (e.dataTransfer) e.dataTransfer.dropEffect = 'copy';
	}

	function winDrop(e: DragEvent) {
		const handled = e.defaultPrevented;
		e.preventDefault();

		if (!handled && hasFiles(e) && canUpload) uploadToFallback(filesFrom(e.dataTransfer));
		resetDrag();
	}

	// --- inline create ('' = root, id = under folder, null = idle) ---
	let creatingParent = $state<string | null>(null);
	let createError = $state<string | null>(null);
	let createSubmitting = $state(false);

	async function startCreate(parentId: string) {
		renamingId = null;
		createError = null;
		creatingParent = parentId;
		if (parentId !== ROOT) expanded[parentId] = true;
		await tick();
		document.getElementById('folder-create-input')?.focus();
	}
	function cancelCreate() {
		creatingParent = null;
		createError = null;
	}

	const submitCreate: SubmitFunction = () => {
		createSubmitting = true;
		return async ({ result }) => {
			createSubmitting = false;
			if (result.type === 'success') {
				creatingParent = null;
				createError = null;
				await invalidateAll();
				showToast(t('doc.created'), 'success');
			} else if (result.type === 'failure') {
				createError = (result.data?.message as string) ?? t('err.generic');
			} else {
				createError = t('err.generic');
			}
		};
	};

	// --- inline rename ---
	let renamingId = $state<string | null>(null);
	let renameError = $state<string | null>(null);
	let renameSubmitting = $state(false);

	async function startRename(node: FolderTreeNode) {
		creatingParent = null;
		renameError = null;
		renamingId = node.id;
		await tick();
		const el = document.getElementById('folder-rename-input') as HTMLInputElement | null;
		el?.focus();
		el?.select();
	}
	function cancelRename() {
		renamingId = null;
		renameError = null;
	}

	const submitRename: SubmitFunction = () => {
		renameSubmitting = true;
		return async ({ result }) => {
			renameSubmitting = false;
			if (result.type === 'success') {
				renamingId = null;
				renameError = null;
				await invalidateAll();
				showToast(t('doc.renamed'), 'success');
			} else if (result.type === 'failure') {
				renameError = (result.data?.message as string) ?? t('err.generic');
			} else {
				renameError = t('err.generic');
			}
		};
	};

	// --- move dialog (keyboard + touch path; drag is the shortcut) ---
	let moveDialog = $state<HTMLDialogElement>();
	let moving = $state<FolderTreeNode | null>(null);
	let moveTarget = $state(ROOT);
	let moveError = $state<string | null>(null);
	let moveSubmitting = $state(false);

	type Option = { id: string; number: string; name: string; depth: number };
	const moveOptions = $derived.by(() => {
		const out: Option[] = [];
		const skip = moving?.id;
		const build = (nodes: FolderTreeNode[], depth: number) => {
			for (const n of nodes) {
				if (n.id === skip) continue;
				out.push({ id: n.id, number: n.number, name: n.name, depth });
				if (n.children?.length) build(n.children, depth + 1);
			}
		};
		build(folders, 0);
		return out;
	});

	function openMove(node: FolderTreeNode) {
		moving = node;
		moveTarget = ROOT;
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
				showToast(t('doc.moved'), 'success');
			} else if (result.type === 'failure') {
				moveError = (result.data?.message as string) ?? t('err.generic');
			} else {
				moveError = t('err.generic');
			}
		};
	};

	// --- delete dialog ---
	let deleteDialog = $state<HTMLDialogElement>();
	let deleting = $state<FolderTreeNode | null>(null);
	let deleteError = $state<string | null>(null);
	let deleteSubmitting = $state(false);
	let deleteConfirm = $state('');
	const deletingKids = $derived(deleting ? descendantCount(deleting) : 0);
	const deleteReady = $derived(!!deleting && deleteConfirm.trim() === deleting.name);

	function openDelete(node: FolderTreeNode) {
		deleting = node;
		deleteConfirm = '';
		deleteError = null;
		deleteDialog?.showModal();
	}

	const submitDelete: SubmitFunction = ({ cancel }) => {
		if (!deleteReady) return cancel();
		deleteSubmitting = true;

		// The server cascades to descendants, so viewing any of them strands the
		// page on a folder that no longer exists.
		const removed = deleting ? subtreeIds(deleting) : [];
		const stranded = !!activeId && removed.includes(activeId);

		return async ({ result }) => {
			deleteSubmitting = false;
			if (result.type === 'success') {
				deleteDialog?.close();
				if (stranded) await goto(actionBase, { invalidateAll: true });
				else await invalidateAll();
				showToast(t('doc.deleted'), 'success');
			} else if (result.type === 'failure') {
				deleteError = (result.data?.message as string) ?? t('err.generic');
			} else {
				deleteError = t('err.generic');
			}
		};
	};
</script>

<svelte:window
	ondragenter={winDragEnter}
	ondragleave={winDragLeave}
	ondragover={winDragOver}
	ondrop={winDrop}
	ondragend={resetDrag}
/>

<svelte:head><title>{t('doc.title')} · {t('brand.name')}</title></svelte:head>

{#snippet folderIcon(open: boolean)}
	<svg
		class="h-4 w-4 flex-none text-muted"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="1.7"
		stroke-linecap="round"
		stroke-linejoin="round"
		aria-hidden="true"
	>
		{#if open}
			<path d="M3 8a2 2 0 0 1 2-2h4l2 2h7a2 2 0 0 1 2 2v1H7l-2 8" />
			<path d="M5 19h13.5a1.5 1.5 0 0 0 1.46-1.14L21.5 11H8.5a1.5 1.5 0 0 0-1.46 1.14z" />
		{:else}
			<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
		{/if}
	</svg>
{/snippet}

{#snippet createInputRow(depth: number, parentId: string)}
	<li class="flex items-start gap-1.5 py-1.5 pr-2" style={indent(depth)}>
		<span class="mt-1.5 w-6 flex-none pointer-coarse:w-9"></span>
		<span class="mt-1.5">{@render folderIcon(false)}</span>
		<div class="min-w-0 flex-1">
			<form
				method="POST"
				action="{actionBase}?/create"
				use:enhance={submitCreate}
				class="flex flex-wrap items-center gap-1.5"
			>
				<input type="hidden" name="parentId" value={parentId} />
				<input
					id="folder-create-input"
					name="name"
					autocomplete="off"
					maxlength="120"
					placeholder={t('doc.namePlaceholder')}
					aria-label={t('doc.namePlaceholder')}
					class="input input-sm w-full max-w-56"
					onkeydown={(e) => e.key === 'Escape' && cancelCreate()}
				/>
				<button type="submit" class="btn btn-primary btn-sm" disabled={createSubmitting}>
					{createSubmitting ? t('doc.adding') : t('doc.add')}
				</button>
				<button type="button" class="btn btn-ghost btn-sm" onclick={cancelCreate}>
					{t('doc.cancel')}
				</button>
			</form>
			{#if createError}<p class="mt-1 text-xs text-error">{createError}</p>{/if}
		</div>
	</li>
{/snippet}

<div class="mx-auto w-full max-w-7xl px-6 py-8">
	<header class="flex flex-wrap items-end justify-between gap-3">
		<div>
			<h1 class="text-2xl font-semibold tracking-[-0.02em]">{t('doc.title')}</h1>
			<p class="mt-1.5 text-sm text-muted">
				{t('doc.desc')}
				{#if totalCount > 0}
					<span aria-hidden="true"> · </span>
					<span class="font-mono text-xs">
						{t(totalCount === 1 ? 'doc.countOne' : 'doc.countMany', { n: totalCount })}
					</span>
				{/if}
			</p>
		</div>
		{#if canCreate}
			<Button onclick={() => startCreate(ROOT)}>
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
					<path d="M12 5v14M5 12h14" />
				</svg>
				{t('doc.newFolder')}
			</Button>
		{/if}
	</header>

	<div class="mt-6 grid gap-6 lg:grid-cols-[minmax(17rem,20rem)_1fr] lg:items-start">
		<!-- Index rail -->
		<nav
			aria-label={t('doc.index')}
			ondragover={railDragOver}
			ondragleave={rowDragLeave}
			ondrop={railDrop}
			class="flex min-h-64 flex-col rounded-box border bg-base-100 transition-colors lg:sticky lg:top-6 lg:max-h-[calc(100dvh-3rem)]
				{dropTarget === ROOT ? 'border-primary/50 bg-primary/[0.04]' : 'border-base-content/10'}"
		>
			<div class="flex-none border-b border-base-content/8">
				<div class="flex items-center justify-between gap-2 px-3 py-2">
					<h2 class="text-xs font-medium text-muted">{t('doc.index')}</h2>
					{#if branchIds.length > 0}
						<button
							type="button"
							onclick={() => setAll(anyCollapsed)}
							class="rounded-field px-1.5 py-0.5 text-xs text-muted transition-colors hover:bg-base-content/5 hover:text-base-content"
						>
							{anyCollapsed ? t('doc.expandAll') : t('doc.collapseAll')}
						</button>
					{/if}
				</div>

				{#if folders.length > 0}
					<div class="relative px-2 pb-2">
						<input
							type="search"
							bind:value={query}
							placeholder={t('doc.search.placeholder')}
							aria-label={t('doc.search.label')}
							autocomplete="off"
							class="input input-sm w-full pr-7"
						/>
						{#if query}
							<button
								type="button"
								onclick={() => (query = '')}
								aria-label={t('doc.search.clear')}
								class="absolute inset-y-0 right-3.5 my-auto grid h-6 w-6 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content"
							>
								<svg
									class="h-3.5 w-3.5"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									stroke-linejoin="round"
									aria-hidden="true"
								>
									<path d="M18 6 6 18M6 6l12 12" />
								</svg>
							</button>
						{/if}
					</div>
				{/if}
			</div>

			{#if folders.length === 0 && creatingParent === null}
				<div class="flex flex-1 flex-col items-center justify-center gap-3 px-6 py-12 text-center">
					<svg
						class="h-8 w-8 text-muted/70"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.4"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
						<path d="M12 11v5M9.5 13.5h5" />
					</svg>
					<div>
						<p class="text-sm font-medium">
							{noAccess ? t('doc.noAccess.title') : t('doc.empty.title')}
						</p>
						<p class="mx-auto mt-1 max-w-xs text-xs text-muted text-pretty">
							{noAccess ? t('doc.noAccess.body') : t('doc.empty.body')}
						</p>
					</div>
					{#if noAccess}
						<span></span>
					{:else if canCreate}
						<Button onclick={() => startCreate(ROOT)}>{t('doc.empty.cta')}</Button>
					{:else}
						<p class="text-xs text-muted">{t('doc.empty.readonly')}</p>
					{/if}
				</div>
			{:else if matched && rows.length === 0}
				<div class="flex flex-1 flex-col items-center justify-center gap-3 px-6 py-12 text-center">
					<p class="text-sm text-muted text-pretty">{t('doc.search.none', { q: query.trim() })}</p>
					<button
						type="button"
						onclick={() => (query = '')}
						class="rounded-field px-1.5 py-1 text-xs font-medium text-muted transition-colors hover:text-primary"
					>
						{t('doc.search.clear')}
					</button>
				</div>
			{:else}
				<ul class="flex-1 overflow-y-auto py-1">
					{#each rows as row (row.node.id)}
						{@const node = row.node}
						{@const open = row.open}
						{@const renaming = renamingId === node.id}
						{@const active = activeId === node.id}
						{@const targeted = dropTarget === node.id}
						<li
							draggable={canEdit && !node.is_default && !renaming}
							ondragstart={(e) => rowDragStart(e, node)}
							ondragover={(e) => rowDragOver(e, node)}
							ondragleave={rowDragLeave}
							ondrop={(e) => rowDrop(e, node)}
							class="group flex items-start gap-1.5 py-1.5 pr-1 transition-colors
								{targeted
								? 'bg-primary/8 ring-1 ring-primary/40 ring-inset'
								: active
									? 'bg-primary/6'
									: 'hover:bg-base-content/[0.045]'}
								{draggingId === node.id ? 'opacity-40' : ''}"
							style={indent(row.depth)}
						>
							{#if row.hasChildren}
								<button
									type="button"
									onclick={() => toggle(node.id, row.depth)}
									aria-expanded={open}
									aria-label={open ? t('doc.collapse') : t('doc.expand')}
									class="grid h-6 w-6 flex-none place-items-center rounded text-muted transition-colors hover:text-base-content pointer-coarse:h-9 pointer-coarse:w-9"
								>
									<svg
										class="riksa-caret h-3.5 w-3.5 {open ? 'rotate-90' : ''}"
										viewBox="0 0 24 24"
										fill="none"
										stroke="currentColor"
										stroke-width="2"
										stroke-linecap="round"
										stroke-linejoin="round"
										aria-hidden="true"
									>
										<path d="m9 6 6 6-6 6" />
									</svg>
								</button>
							{:else}
								<span class="w-6 flex-none pointer-coarse:w-9"></span>
							{/if}

							<span class="mt-1">{@render folderIcon((open && row.hasChildren) || active)}</span>

							{#if renaming}
								<div class="min-w-0 flex-1">
									<form
										method="POST"
										action="{actionBase}?/rename"
										use:enhance={submitRename}
										class="flex flex-wrap items-center gap-1.5"
									>
										<input type="hidden" name="folderId" value={node.id} />
										<input
											id="folder-rename-input"
											name="name"
											value={node.name}
											autocomplete="off"
											maxlength="120"
											aria-label={t('doc.action.rename')}
											class="input input-sm w-full max-w-56"
											onkeydown={(e) => e.key === 'Escape' && cancelRename()}
										/>
										<button
											type="submit"
											class="btn btn-primary btn-sm"
											disabled={renameSubmitting}
										>
											{renameSubmitting ? t('doc.saving') : t('doc.save')}
										</button>
										<button type="button" class="btn btn-ghost btn-sm" onclick={cancelRename}>
											{t('doc.cancel')}
										</button>
									</form>
									{#if renameError}<p class="mt-1 text-xs text-error">{renameError}</p>{/if}
								</div>
							{:else}
								<a
									href={rowHref(node.id)}
									draggable="false"
									aria-current={active ? 'page' : undefined}
									class="mt-0.5 flex min-w-0 flex-1 items-baseline gap-1.5 rounded-field no-underline"
								>
									<span class="font-mono text-xs tabular-nums text-muted">{node.number}</span>
									<span class="min-w-0 flex-1 truncate text-sm {active ? 'font-medium' : ''}">
										{node.name}
									</span>
								</a>

								{#if node.is_default}
									<span
										class="mt-0.5 flex-none rounded-selector bg-base-content/5 px-1.5 py-0.5 text-[0.6875rem] text-muted"
									>
										{t('doc.defaultBadge')}
									</span>
								{/if}

								{#if row.hasChildren && !open}
									<span
										class="mt-0.5 flex-none rounded-selector bg-base-content/5 px-1.5 py-0.5 font-mono text-[0.6875rem] text-muted"
										title={t(
											node.children.length === 1 ? 'doc.childCountOne' : 'doc.childCountMany',
											{ n: node.children.length }
										)}
									>
										{node.children.length}
									</span>
								{/if}

								{#if canAct}
									<div
										class="-mt-1 ml-1 flex flex-none items-center gap-0.5 opacity-0 transition-opacity focus-within:opacity-100 group-hover:opacity-100 pointer-coarse:-mt-2.5 pointer-coarse:gap-1 pointer-coarse:opacity-100"
									>
										{#if canCreate}
											<button
												type="button"
												onclick={() => startCreate(node.id)}
												title={t('doc.action.addSub')}
												aria-label={t('doc.action.addSubOf', { name: node.name })}
												class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content pointer-coarse:h-11 pointer-coarse:w-11"
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
													<path d="M12 8v8M8 12h8" />
												</svg>
											</button>
										{/if}
										{#if canEdit}
											<button
												type="button"
												onclick={() => startRename(node)}
												title={t('doc.action.rename')}
												aria-label={t('doc.action.renameOf', { name: node.name })}
												class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content pointer-coarse:h-11 pointer-coarse:w-11"
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
													<path d="M12 20h9" />
													<path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4z" />
												</svg>
											</button>
											{#if !node.is_default}
												<button
													type="button"
													onclick={() => openMove(node)}
													title={t('doc.action.move')}
													aria-label={t('doc.action.moveOf', { name: node.name })}
													class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content pointer-coarse:h-11 pointer-coarse:w-11"
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
										{/if}
										{#if canAssign}
											<a
												href={accessHref(node.id)}
												draggable="false"
												title={t('doc.action.access')}
												aria-label={t('doc.action.accessOf', { name: node.name })}
												class="grid h-8 w-8 place-items-center rounded-field text-muted no-underline transition-colors hover:bg-base-content/5 hover:text-base-content pointer-coarse:h-11 pointer-coarse:w-11"
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
													<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
													<circle cx="9" cy="7" r="4" />
													<path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
												</svg>
											</a>
										{/if}
										{#if canDelete && !node.is_default}
											<button
												type="button"
												onclick={() => openDelete(node)}
												title={t('doc.action.delete')}
												aria-label={t('doc.action.deleteOf', { name: node.name })}
												class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-error/10 hover:text-error pointer-coarse:h-11 pointer-coarse:w-11"
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
								{/if}
							{/if}
						</li>

						{#if creatingParent === node.id}
							{@render createInputRow(row.depth + 1, node.id)}
						{/if}
					{/each}

					{#if creatingParent === ROOT}
						{@render createInputRow(0, ROOT)}
					{/if}
				</ul>

				{#if canCreate && creatingParent !== ROOT}
					<button
						type="button"
						onclick={() => startCreate(ROOT)}
						class="m-2 mt-1 inline-flex flex-none items-center gap-1.5 self-start rounded-field px-1.5 py-1 text-xs font-medium text-muted transition-colors hover:text-primary"
					>
						<svg
							class="h-3.5 w-3.5"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							stroke-linejoin="round"
							aria-hidden="true"
						>
							<path d="M12 5v14M5 12h14" />
						</svg>
						{t('doc.newRootFolder')}
					</button>
				{/if}
			{/if}
		</nav>

		<!-- Documents pane -->
		<div class="min-w-0">
			{@render children()}
		</div>
	</div>
</div>

{#if fileDragging && canUpload}
	<div
		class="riksa-overlay pointer-events-none fixed inset-x-0 top-4 z-overlay flex justify-center px-4"
		aria-hidden="true"
	>
		<div
			class="max-w-full rounded-box border border-primary/40 bg-base-100/95 px-3.5 py-2 shadow-sm"
		>
			<p class="truncate text-xs">
				{#if dropLabel}
					{t('doc.dropAnywhere.body', { name: dropLabel })}
				{:else}
					{t('doc.upload.noDefault')}
				{/if}
			</p>
		</div>
	</div>
{/if}

<div aria-live="polite" class="sr-only">
	{#if dropTarget !== null && dropLabel}
		{t('doc.dropAnywhere.body', { name: dropLabel })}
	{/if}
</div>

<UploadQueue />

<!-- Move -->
<dialog bind:this={moveDialog} class="modal" aria-labelledby="folder-move-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="folder-move-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('doc.move.title')}
		</h2>
		{#if moving}
			<p class="mt-1 text-sm text-muted text-pretty">
				{t('doc.move.desc', { name: moving.name })}
			</p>
		{/if}

		{#if moveError}
			<div class="mt-4"><Alert align="start">{moveError}</Alert></div>
		{/if}

		<form method="POST" action="{actionBase}?/move" use:enhance={submitMove} class="mt-5">
			<input type="hidden" name="folderId" value={moving?.id ?? ''} />
			<label for="move-dest" class="mb-1.5 block text-sm font-medium">{t('doc.move.dest')}</label>
			<select
				id="move-dest"
				name="parentId"
				bind:value={moveTarget}
				class="select select-sm w-full font-mono"
			>
				<option value={ROOT}>{t('doc.move.root')}</option>
				{#each moveOptions as opt (opt.id)}
					<option value={opt.id}>{'  '.repeat(opt.depth)}{opt.number} {opt.name}</option>
				{/each}
			</select>

			<div class="mt-6 flex justify-end gap-2">
				<Button type="button" variant="ghost" onclick={() => moveDialog?.close()}>
					{t('doc.cancel')}
				</Button>
				<Button type="submit" loading={moveSubmitting}>
					{moveSubmitting ? t('doc.move.submitting') : t('doc.move.submit')}
				</Button>
			</div>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('doc.cancel')}></button>
	</form>
</dialog>

<!-- Delete -->
<dialog bind:this={deleteDialog} class="modal" aria-labelledby="folder-delete-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="folder-delete-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('doc.delete.title')}
		</h2>
		{#if deleting}
			<p class="mt-1 text-sm text-muted text-pretty">
				{t('doc.delete.warning', { name: deleting.name })}
			</p>
			<p class="mt-2 text-sm font-medium text-error text-pretty">{t('doc.delete.contents')}</p>
			{#if deletingKids > 0}
				<p class="mt-1 text-sm text-error text-pretty">
					{t(deletingKids === 1 ? 'doc.delete.cascadeOne' : 'doc.delete.cascadeMany', {
						n: deletingKids
					})}
				</p>
			{/if}
		{/if}

		{#if deleteError}
			<div class="mt-4"><Alert align="start">{deleteError}</Alert></div>
		{/if}

		<form
			method="POST"
			action="{actionBase}?/delete"
			use:enhance={submitDelete}
			class="mt-5 flex flex-col gap-4"
		>
			<input type="hidden" name="folderId" value={deleting?.id ?? ''} />
			<Field
				id="folder-delete-confirm"
				name="confirm"
				label={t('doc.delete.confirmLabel', { name: deleting?.name ?? '' })}
				bind:value={deleteConfirm}
				placeholder={deleting?.name ?? ''}
				autocomplete="off"
			/>
			<div class="mt-1 flex justify-end gap-2">
				<Button type="button" variant="ghost" onclick={() => deleteDialog?.close()}>
					{t('doc.cancel')}
				</Button>
				<Button type="submit" variant="danger" disabled={!deleteReady} loading={deleteSubmitting}>
					{deleteSubmitting ? t('doc.delete.submitting') : t('doc.delete.submit')}
				</Button>
			</div>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('doc.cancel')}></button>
	</form>
</dialog>

<Toaster />

<style>
	.riksa-caret {
		transition: transform 150ms cubic-bezier(0.22, 1, 0.36, 1);
	}
	.riksa-overlay {
		animation: riksa-fade-in 150ms ease-out;
	}
	@keyframes riksa-fade-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
	@media (prefers-reduced-motion: reduce) {
		.riksa-caret {
			transition: none;
		}
		.riksa-overlay {
			animation: none;
		}
	}
</style>
