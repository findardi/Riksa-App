<script lang="ts">
	import { tick } from 'svelte';
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button, showToast } from '$lib/components/common';
	import { t } from '$lib/i18n';
	import type { FolderTreeNode } from '$lib/types/content';
	import type { MyAccessWorkspace } from '$lib/types/workspace';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const folders = $derived(data.folders);

	const ROOT = '';

	const perms = $derived((page.data as { access?: MyAccessWorkspace }).access?.permissions ?? []);
	const canCreate = $derived(perms.includes('folder:create'));
	const canEdit = $derived(perms.includes('folder:edit'));
	const canDelete = $derived(perms.includes('folder:delete'));
	const canAct = $derived(canCreate || canEdit || canDelete);

	// Expanded by default; toggling records a false. Reset on navigation via load.
	let expanded = $state<Record<string, boolean>>({});
	const isExpanded = (id: string) => expanded[id] ?? true;
	function toggle(id: string) {
		expanded[id] = !isExpanded(id);
	}

	type Row = { node: FolderTreeNode; depth: number; hasChildren: boolean };
	function walk(nodes: FolderTreeNode[], depth: number, out: Row[]) {
		for (const n of nodes) {
			const kids = n.children ?? [];
			out.push({ node: n, depth, hasChildren: kids.length > 0 });
			if (kids.length && isExpanded(n.id)) walk(kids, depth + 1, out);
		}
	}
	const rows = $derived.by(() => {
		const out: Row[] = [];
		walk(folders, 0, out);
		return out;
	});

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

	function descendantCount(node: FolderTreeNode): number {
		let n = 0;
		for (const c of node.children ?? []) n += 1 + descendantCount(c);
		return n;
	}

	const indent = (depth: number) => `padding-left: ${depth * 1.375 + 0.25}rem`;

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

	// --- move dialog ---
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
				if (n.id === skip) continue; // skip the moving node and its whole subtree
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
	const deletingKids = $derived(deleting ? descendantCount(deleting) : 0);

	function openDelete(node: FolderTreeNode) {
		deleting = node;
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
				showToast(t('doc.deleted'), 'success');
			} else if (result.type === 'failure') {
				deleteError = (result.data?.message as string) ?? t('err.generic');
			} else {
				deleteError = t('err.generic');
			}
		};
	};
</script>

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
		<span class="mt-1.5 w-5 flex-none"></span>
		<span class="mt-1.5">{@render folderIcon(false)}</span>
		<div class="min-w-0 flex-1">
			<form
				method="POST"
				action="?/create"
				use:enhance={submitCreate}
				class="flex flex-wrap items-center gap-1.5"
			>
				<input type="hidden" name="parentId" value={parentId} />
				<!-- svelte-ignore a11y_autofocus -->
				<input
					id="folder-create-input"
					name="name"
					autocomplete="off"
					maxlength="120"
					placeholder={t('doc.namePlaceholder')}
					aria-label={t('doc.namePlaceholder')}
					class="input input-sm w-full max-w-64"
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

<div class="mx-auto w-full max-w-4xl px-6 py-8">
	<header class="flex flex-wrap items-end justify-between gap-3">
		<div>
			<h1 class="text-2xl font-semibold tracking-[-0.02em]">{t('doc.title')}</h1>
			<p class="mt-1.5 text-sm text-muted">
				{t('doc.desc')}
				{#if totalCount > 0}
					<span aria-hidden="true"> · </span>
					<span class="font-mono text-xs">{t('doc.count', { n: totalCount })}</span>
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

	{#if folders.length === 0 && creatingParent === null}
		<div
			class="mt-8 flex flex-col items-center justify-center gap-3 border-y border-base-content/10 px-6 py-16 text-center"
		>
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
				<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" />
				<path d="M12 11v5M9.5 13.5h5" />
			</svg>
			<div>
				<p class="text-[0.9375rem] font-medium">{t('doc.empty.title')}</p>
				<p class="mx-auto mt-1 max-w-sm text-sm text-muted text-pretty">{t('doc.empty.body')}</p>
			</div>
			{#if canCreate}
				<div class="mt-1">
					<Button onclick={() => startCreate(ROOT)}>{t('doc.empty.cta')}</Button>
				</div>
			{:else}
				<p class="text-xs text-muted">{t('doc.empty.readonly')}</p>
			{/if}
		</div>
	{:else}
		<ul
			class="mt-6 divide-y divide-base-content/[0.06] border-y border-base-content/10"
			aria-label={t('doc.title')}
		>
			{#each rows as row (row.node.id)}
				{@const node = row.node}
				{@const open = isExpanded(node.id)}
				{@const renaming = renamingId === node.id}
				<li
					class="group flex items-start gap-1.5 py-1.5 pr-1 transition-colors hover:bg-base-content/[0.025]"
					style={indent(row.depth)}
				>
					{#if row.hasChildren}
						<button
							type="button"
							onclick={() => toggle(node.id)}
							aria-expanded={open}
							aria-label={open ? t('doc.collapse') : t('doc.expand')}
							class="mt-0.5 grid h-5 w-5 flex-none place-items-center rounded text-muted transition-colors hover:text-base-content"
						>
							<svg
								class="h-3.5 w-3.5 transition-transform duration-150 {open ? 'rotate-90' : ''}"
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
						<span class="mt-0.5 w-5 flex-none"></span>
					{/if}

					<span class="mt-0.5">{@render folderIcon(open && row.hasChildren)}</span>

					{#if renaming}
						<div class="min-w-0 flex-1">
							<form
								method="POST"
								action="?/rename"
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
									class="input input-sm w-full max-w-64"
									onkeydown={(e) => e.key === 'Escape' && cancelRename()}
								/>
								<button type="submit" class="btn btn-primary btn-sm" disabled={renameSubmitting}>
									{renameSubmitting ? t('doc.saving') : t('doc.save')}
								</button>
								<button type="button" class="btn btn-ghost btn-sm" onclick={cancelRename}>
									{t('doc.cancel')}
								</button>
							</form>
							{#if renameError}<p class="mt-1 text-xs text-error">{renameError}</p>{/if}
						</div>
					{:else}
						<span class="mt-[0.15rem] font-mono text-xs tabular-nums text-muted">{node.number}</span>
						<span class="mt-0 min-w-0 flex-1 truncate text-sm">{node.name}</span>

						{#if row.hasChildren && !open}
							<span
								class="mt-0.5 rounded-selector bg-base-content/5 px-1.5 py-0.5 font-mono text-[0.6875rem] text-muted"
								title={t('doc.childCount', { n: node.children.length })}
							>
								{node.children.length}
							</span>
						{/if}

						{#if canAct}
							<div
								class="ml-1 flex flex-none items-center gap-0.5 opacity-0 transition-opacity focus-within:opacity-100 group-hover:opacity-100"
							>
								{#if canCreate}
									<button
										type="button"
										onclick={() => startCreate(node.id)}
										title={t('doc.action.addSub')}
										aria-label={t('doc.action.addSubOf', { name: node.name })}
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
											<path d="M12 20h9" />
											<path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4z" />
										</svg>
									</button>
									<button
										type="button"
										onclick={() => openMove(node)}
										title={t('doc.action.move')}
										aria-label={t('doc.action.moveOf', { name: node.name })}
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
										onclick={() => openDelete(node)}
										title={t('doc.action.delete')}
										aria-label={t('doc.action.deleteOf', { name: node.name })}
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
											<path d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2m2 0v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6" />
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
				class="mt-3 inline-flex items-center gap-1.5 px-1 text-sm font-medium text-muted transition-colors hover:text-primary"
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
					<path d="M12 5v14M5 12h14" />
				</svg>
				{t('doc.newRootFolder')}
			</button>
		{/if}
	{/if}
</div>

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

		<form method="POST" action="?/move" use:enhance={submitMove} class="mt-5">
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
					<option value={opt.id}>{'  '.repeat(opt.depth)}{opt.number} {opt.name}</option>
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
			{#if deletingKids > 0}
				<p class="mt-2 text-sm font-medium text-error text-pretty">
					{t('doc.delete.cascade', { n: deletingKids })}
				</p>
			{/if}
		{/if}

		{#if deleteError}
			<div class="mt-4"><Alert align="start">{deleteError}</Alert></div>
		{/if}

		<form
			method="POST"
			action="?/delete"
			use:enhance={submitDelete}
			class="mt-6 flex justify-end gap-2"
		>
			<input type="hidden" name="folderId" value={deleting?.id ?? ''} />
			<Button type="button" variant="ghost" onclick={() => deleteDialog?.close()}>
				{t('doc.cancel')}
			</Button>
			<Button type="submit" variant="danger" loading={deleteSubmitting}>
				{deleteSubmitting ? t('doc.delete.submitting') : t('doc.delete.submit')}
			</Button>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('doc.cancel')}></button>
	</form>
</dialog>
