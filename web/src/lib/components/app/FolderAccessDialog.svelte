<script lang="ts">
	import { tick } from 'svelte';
	import { enhance } from '$app/forms';
	import { resolve } from '$app/paths';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button, showToast } from '$lib/components/common';
	import { t, type TKey } from '$lib/i18n';
	import type {
		AccessLevelData,
		FolderAccessPanel,
		FolderTreeNode,
		InheritedFolderAccess
	} from '$lib/types/content';
	import type { GroupWorkspaceData } from '$lib/types/workspace';

	type Props = {
		folder: FolderTreeNode | null;
		groups: GroupWorkspaceData[];
		levels: AccessLevelData[];
		ready: boolean;
		workspaceId: string;
		actionBase: string;
		slug: string;
		open?: boolean;
	};

	let {
		folder,
		groups,
		levels,
		ready,
		workspaceId,
		actionBase,
		slug,
		open = $bindable(false)
	}: Props = $props();

	const LEVEL_KEYS: Record<string, TKey> = {
		view: 'level.view',
		download: 'level.download',
		watermark: 'level.watermark'
	};

	function levelLabel(name: string): string {
		const key = LEVEL_KEYS[name.toLowerCase()];
		return key ? t(key) : name;
	}

	let dialog = $state<HTMLDialogElement>();
	let panel = $state<FolderAccessPanel | null>(null);
	let loading = $state(false);
	let loadError = $state<string | null>(null);
	let formError = $state<string | null>(null);

	let adding = $state(false);
	let addGroupId = $state('');
	let addLevelId = $state('');
	let addSubmitting = $state(false);
	let savingGroup = $state<string | null>(null);
	let revokingGroup = $state<string | null>(null);

	const direct = $derived(panel?.direct ?? []);
	const inherited = $derived(panel?.inherited ?? []);
	const directIds = $derived(new Set(direct.map((r) => r.group_id)));
	const addable = $derived(groups.filter((g) => !directIds.has(g.id)));
	const defaultLevelId = $derived(levels.find((l) => l.name === 'view')?.id ?? levels[0]?.id ?? '');
	const selectedLevel = $derived(levels.find((l) => l.id === addLevelId) ?? null);

	const selectedCaps = $derived.by(() => {
		const l = selectedLevel;
		if (!l) return [] as string[];
		const caps: string[] = [];
		if (l.can_view) caps.push(t('level.cap.view'));
		if (l.can_download) caps.push(t('level.cap.download'));
		if (l.can_watermark) caps.push(t('level.cap.watermark'));
		return caps;
	});

	function endpoint(folderId: string): string {
		const q = new URLSearchParams({ workspaceId, folderId });
		return `/api/content/folder-access?${q}`;
	}

	async function refresh(): Promise<FolderAccessPanel | null> {
		if (!folder) return null;
		try {
			const res = await fetch(endpoint(folder.id));
			if (!res.ok) return null;
			panel = (await res.json()) as FolderAccessPanel;
			return panel;
		} catch {
			return null;
		}
	}

	async function load(folderId: string) {
		loading = true;
		loadError = null;
		panel = null;
		try {
			const res = await fetch(endpoint(folderId));
			if (!res.ok) throw new Error(String(res.status));
			panel = (await res.json()) as FolderAccessPanel;
		} catch {
			loadError = t('facc.err.load');
		} finally {
			loading = false;
		}
	}

	$effect(() => {
		if (!open || !folder) return;
		adding = false;
		addGroupId = '';
		addLevelId = '';
		formError = null;
		if (!dialog?.open) dialog?.showModal();
		void load(folder.id);
	});

	function focusIn(selector: string) {
		tick().then(() => dialog?.querySelector<HTMLElement>(selector)?.focus());
	}

	function startAdd() {
		adding = true;
		addGroupId = addable[0]?.id ?? '';
		addLevelId = defaultLevelId;
		formError = null;
		focusIn('#facc-add-group');
	}

	function startOverride(row: InheritedFolderAccess) {
		adding = true;
		addGroupId = row.group_id;
		addLevelId = row.level_id;
		formError = null;
		focusIn('#facc-add-level');
	}

	function cancelAdd() {
		adding = false;
		formError = null;
	}

	const submitLevel =
		(groupId: string, groupName: string): SubmitFunction =>
		() => {
			savingGroup = groupId;
			formError = null;
			return async ({ result }) => {
				savingGroup = null;
				if (result.type === 'success') {
					await refresh();
					showToast(t('facc.saved', { group: groupName }), 'success');
				} else if (result.type === 'failure') {
					formError = (result.data?.message as string) ?? t('err.generic');
					await refresh();
				} else {
					formError = t('err.generic');
				}
			};
		};

	const submitRevoke =
		(groupId: string, groupName: string): SubmitFunction =>
		() => {
			revokingGroup = groupId;
			formError = null;
			return async ({ result }) => {
				revokingGroup = null;
				if (result.type === 'success') {
					const next = await refresh();
					const back = next?.inherited.find((r) => r.group_id === groupId);
					showToast(
						back
							? t('facc.revokedInherits', { group: groupName, name: back.source_folder_name })
							: t('facc.revoked', { group: groupName }),
						'success'
					);
				} else if (result.type === 'failure') {
					formError = (result.data?.message as string) ?? t('err.generic');
				} else {
					formError = t('err.generic');
				}
			};
		};

	const submitAdd: SubmitFunction = ({ cancel }) => {
		if (!addGroupId || !addLevelId) {
			formError = t('facc.err.pick');
			cancel();
			return;
		}
		const groupName = groups.find((g) => g.id === addGroupId)?.name ?? '';
		addSubmitting = true;
		formError = null;
		return async ({ result }) => {
			addSubmitting = false;
			if (result.type === 'success') {
				await refresh();
				adding = false;
				addGroupId = '';
				showToast(t('facc.saved', { group: groupName }), 'success');
			} else if (result.type === 'failure') {
				formError = (result.data?.message as string) ?? t('err.generic');
			} else {
				formError = t('err.generic');
			}
		};
	};
</script>

<dialog
	bind:this={dialog}
	onclose={() => (open = false)}
	class="modal"
	aria-labelledby="facc-title"
>
	<div class="modal-box w-full max-w-lg rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="facc-title" class="text-lg font-semibold tracking-[-0.01em]">{t('facc.title')}</h2>
		{#if folder}
			<p class="mt-1 flex items-baseline gap-2 text-sm text-muted">
				{#if folder.number}
					<span class="font-mono text-xs">{folder.number}</span>
				{/if}
				<span class="min-w-0 truncate">{folder.name}</span>
			</p>
		{/if}

		{#if loading}
			<div class="mt-5 space-y-2" aria-busy="true">
				{#each [0, 1] as i (i)}
					<div class="h-11 animate-pulse rounded-field bg-base-content/5"></div>
				{/each}
			</div>
		{:else if loadError}
			<div class="mt-5"><Alert align="start">{loadError}</Alert></div>
		{:else if !ready}
			<div class="mt-5"><Alert align="start">{t('facc.err.load')}</Alert></div>
		{:else if groups.length === 0}
			<div class="mt-5 rounded-box border border-base-content/10 p-5 text-center">
				<p class="text-sm font-semibold">{t('facc.noGroups.title')}</p>
				<p class="mx-auto mt-1 max-w-[46ch] text-sm text-muted text-pretty">
					{t('facc.noGroups.body')}
				</p>
				<a
					href={resolve('/(app)/workspace/[slug]/management-access/group', { slug })}
					class="mt-3 inline-block text-sm font-medium text-primary hover:underline"
				>
					{t('facc.noGroups.cta')}
				</a>
			</div>
		{:else}
			{#if formError}
				<div class="mt-4"><Alert align="start">{formError}</Alert></div>
			{/if}

			<section class="mt-5">
				<div class="flex items-baseline justify-between gap-2">
					<h3 class="text-sm font-semibold">{t('facc.direct')}</h3>
					{#if direct.length}
						<span class="font-mono text-xs text-muted">
							{t('facc.directCount', { n: direct.length })}
						</span>
					{/if}
				</div>

				{#if direct.length}
					<ul class="mt-1 divide-y divide-base-content/8">
						{#each direct as row (row.group_id)}
							<li class="flex items-center gap-2 py-2">
								<span class="min-w-0 flex-1 truncate text-sm font-medium">{row.group_name}</span>

								<form
									method="POST"
									action="{actionBase}?/setAccess"
									use:enhance={submitLevel(row.group_id, row.group_name)}
									class="flex-none"
								>
									<input type="hidden" name="folderId" value={folder?.id ?? ''} />
									<input type="hidden" name="groupId" value={row.group_id} />
									<select
										name="levelId"
										value={row.level_id}
										disabled={savingGroup === row.group_id}
										onchange={(e) => e.currentTarget.form?.requestSubmit()}
										aria-label={t('facc.levelOf', { group: row.group_name })}
										class="select select-sm w-44"
									>
										{#each levels as l (l.id)}
											<option value={l.id}>{levelLabel(l.name)}</option>
										{/each}
									</select>
								</form>

								<form
									method="POST"
									action="{actionBase}?/removeAccess"
									use:enhance={submitRevoke(row.group_id, row.group_name)}
									class="flex-none"
								>
									<input type="hidden" name="folderId" value={folder?.id ?? ''} />
									<input type="hidden" name="groupId" value={row.group_id} />
									<button
										type="submit"
										disabled={revokingGroup === row.group_id}
										aria-label={t('facc.revokeOf', { group: row.group_name })}
										title={t('facc.revoke')}
										class="grid h-9 w-9 place-items-center rounded-field text-muted transition-colors hover:bg-error/10 hover:text-error disabled:pointer-events-none disabled:opacity-50"
									>
										{#if revokingGroup === row.group_id}
											<span class="loading loading-spinner loading-xs"></span>
										{:else}
											<svg
												class="h-4 w-4"
												viewBox="0 0 24 24"
												fill="none"
												stroke="currentColor"
												stroke-width="1.8"
												stroke-linecap="round"
												aria-hidden="true"
											>
												<path d="M18 6 6 18M6 6l12 12" />
											</svg>
										{/if}
									</button>
								</form>
							</li>
						{/each}
					</ul>
				{:else if !inherited.length}
					<div class="mt-2 rounded-box border border-base-content/10 p-5 text-center">
						<p class="text-sm font-semibold">{t('facc.empty.title')}</p>
						<p class="mx-auto mt-1 max-w-[48ch] text-sm text-muted text-pretty">
							{t('facc.empty.body')}
						</p>
					</div>
				{/if}

				{#if adding}
					<form
						method="POST"
						action="{actionBase}?/setAccess"
						use:enhance={submitAdd}
						class="mt-3 rounded-box border border-base-content/10 p-3"
					>
						<input type="hidden" name="folderId" value={folder?.id ?? ''} />
						<div class="flex flex-wrap items-end gap-2">
							<div class="min-w-[9rem] flex-1">
								<label class="text-xs font-medium" for="facc-add-group">
									{t('facc.add.group')}
								</label>
								<select
									id="facc-add-group"
									name="groupId"
									bind:value={addGroupId}
									class="select select-sm mt-1 w-full"
								>
									{#each addable as g (g.id)}
										<option value={g.id}>{g.name}</option>
									{/each}
								</select>
							</div>
							<div class="min-w-[9rem] flex-1">
								<label class="text-xs font-medium" for="facc-add-level">
									{t('facc.add.level')}
								</label>
								<select
									id="facc-add-level"
									name="levelId"
									bind:value={addLevelId}
									class="select select-sm mt-1 w-full"
								>
									{#each levels as l (l.id)}
										<option value={l.id}>{levelLabel(l.name)}</option>
									{/each}
								</select>
							</div>
						</div>

						{#if selectedCaps.length}
							<p class="mt-2 text-xs text-muted">{selectedCaps.join(' · ')}</p>
						{/if}

						<div class="mt-3 flex justify-end gap-2">
							<Button type="button" variant="ghost" onclick={cancelAdd}>
								{t('facc.add.cancel')}
							</Button>
							<Button type="submit" loading={addSubmitting} disabled={!addGroupId || !addLevelId}>
								{addSubmitting ? t('facc.add.submitting') : t('facc.add.submit')}
							</Button>
						</div>
					</form>
				{:else if addable.length}
					<button
						type="button"
						onclick={startAdd}
						class="mt-3 inline-flex items-center gap-1.5 rounded-field px-2 py-1.5 text-sm font-medium text-primary transition-colors hover:bg-primary/8"
					>
						<svg
							class="h-4 w-4"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							aria-hidden="true"
						>
							<path d="M12 5v14M5 12h14" />
						</svg>
						{t('facc.add')}
					</button>
				{:else}
					<p class="mt-3 text-xs text-muted">{t('facc.allGranted')}</p>
				{/if}
			</section>

			{#if inherited.length}
				<section class="mt-6">
					<h3 class="text-sm font-semibold">{t('facc.inherited')}</h3>
					<ul class="mt-1 divide-y divide-base-content/8">
						{#each inherited as row (row.group_id)}
							<li class="flex items-center gap-2 py-2">
								<div class="min-w-0 flex-1">
									<p class="truncate text-sm font-medium">{row.group_name}</p>
									<p class="mt-0.5 truncate text-xs text-muted">
										{t('facc.inheritedFrom', { name: row.source_folder_name })}
									</p>
								</div>
								<span class="flex-none text-sm text-muted">{levelLabel(row.level_name)}</span>
								<button
									type="button"
									onclick={() => startOverride(row)}
									aria-label={t('facc.overrideOf', { group: row.group_name })}
									class="flex-none rounded-field px-2 py-1.5 text-sm font-medium text-primary transition-colors hover:bg-primary/8"
								>
									{t('facc.override')}
								</button>
							</li>
						{/each}
					</ul>
				</section>
			{/if}

			<p
				class="mt-6 flex items-start gap-2 rounded-box border border-base-content/10 bg-base-content/3 p-3 text-xs text-muted text-pretty"
			>
				<svg
					class="mt-px h-4 w-4 flex-none"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.8"
					stroke-linecap="round"
					stroke-linejoin="round"
					aria-hidden="true"
				>
					<circle cx="12" cy="12" r="9" />
					<path d="M12 16v-5M12 8h.01" />
				</svg>
				<span>{t('facc.flow')}</span>
			</p>
		{/if}

		<div class="mt-6 flex justify-end">
			<Button type="button" variant="ghost" onclick={() => dialog?.close()}>
				{t('facc.close')}
			</Button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('facc.close')}></button>
	</form>
</dialog>
