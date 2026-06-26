<script lang="ts">
	import { onMount } from 'svelte';
	import { enhance } from '$app/forms';
	import { invalidateAll, replaceState } from '$app/navigation';
	import { page } from '$app/state';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button } from '$lib/components/common';
	import { roleDisplayName } from '$lib/access/permissions';
	import { t } from '$lib/i18n';
	import type { WorkspaceRoleData } from '$lib/types/workspace';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const roles = $derived(data.roles);
	const customCount = $derived(roles.filter((r) => !r.is_system).length);

	const base = $derived(`/workspace/${page.params.slug}/management-access/role`);
	const permLabel = (r: WorkspaceRoleData) =>
		r.permissions.length === 0
			? t('role.permNone')
			: t('role.permCount', { n: r.permissions.length });

	// --- Toast ---
	let toastMsg = $state<string | null>(null);
	let toastVariant = $state<'success' | 'error'>('success');
	let toastTimer: ReturnType<typeof setTimeout>;
	function showToast(msg: string, variant: 'success' | 'error') {
		toastMsg = msg;
		toastVariant = variant;
		clearTimeout(toastTimer);
		toastTimer = setTimeout(() => (toastMsg = null), 4000);
	}

	// Flash from the editor redirect (?flash=created|updated) → confirm once, then strip.
	onMount(() => {
		const flash = page.url.searchParams.get('flash');
		if (flash === 'created') showToast(t('role.created'), 'success');
		else if (flash === 'updated') showToast(t('role.updated'), 'success');
		if (flash) {
			const url = new URL(page.url);
			url.searchParams.delete('flash');
			replaceState(url, page.state);
		}
	});

	// --- Delete ---
	let deleteDialog = $state<HTMLDialogElement>();
	let pending = $state<WorkspaceRoleData | null>(null);
	let deleteSubmitting = $state(false);
	let deleteMessage = $state<string | null>(null);

	function openDelete(role: WorkspaceRoleData) {
		pending = role;
		deleteMessage = null;
		deleteDialog?.showModal();
	}

	const submitDelete: SubmitFunction = () => {
		deleteSubmitting = true;
		return async ({ result }) => {
			deleteSubmitting = false;
			if (result.type === 'success') {
				deleteDialog?.close();
				await invalidateAll();
				showToast(t('role.deleted'), 'success');
			} else if (result.type === 'failure') {
				deleteMessage = (result.data?.message as string) ?? t('err.generic');
			} else {
				deleteMessage = t('err.generic');
			}
		};
	};
</script>

<svelte:head><title>{t('ma.role')} · {t('ma.title')}</title></svelte:head>

<div class="flex items-center justify-between gap-4">
	<h2 class="text-sm font-semibold">
		{t('ma.role')}
		<span class="ml-1 font-mono text-xs font-normal text-muted">{roles.length}</span>
	</h2>
	<a href="{base}/new" class="btn btn-primary btn-sm">
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
		{t('role.new')}
	</a>
</div>

<ul class="mt-4 divide-y divide-base-content/10 border-y border-base-content/10">
	{#each roles as role (role.id)}
		{@const href = `${base}/${role.id}`}
		<li class="flex items-center gap-4 py-3">
			<div class="min-w-0 flex-1">
				<div class="flex items-center gap-2">
					<span class="truncate text-[0.9375rem] font-medium">{roleDisplayName(role.name)}</span>
					{#if role.is_system}
						<span
							class="rounded-selector bg-base-content/10 px-1.5 py-0.5 text-[0.6875rem] font-medium text-muted"
							>{t('role.system')}</span
						>
					{/if}
				</div>
				<p class="mt-0.5 font-mono text-xs text-muted">{permLabel(role)}</p>
			</div>

			{#if role.is_system}
				<a
					{href}
					class="inline-flex flex-none items-center gap-1.5 rounded-field px-2.5 py-2.5 text-sm font-medium text-muted transition-colors hover:bg-base-content/5 hover:text-base-content"
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
						<path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z" />
						<circle cx="12" cy="12" r="3" />
					</svg>
					{t('role.view')}
				</a>
			{:else}
				<div class="flex flex-none items-center gap-1">
					<a
						{href}
						class="inline-flex items-center gap-1.5 rounded-field px-2.5 py-2.5 text-sm font-medium text-muted transition-colors hover:bg-base-content/5 hover:text-base-content"
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
							<path d="M16.5 3.5a2.12 2.12 0 0 1 3 3L7 19l-4 1 1-4Z" />
						</svg>
						{t('role.edit')}
					</a>
					<button
						type="button"
						onclick={() => openDelete(role)}
						class="inline-flex items-center gap-1.5 rounded-field px-2.5 py-2.5 text-sm font-medium text-muted transition-colors hover:bg-error/10 hover:text-error"
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
						{t('role.delete')}
					</button>
				</div>
			{/if}
		</li>
	{/each}
</ul>

{#if customCount === 0}
	<p class="mt-4 max-w-[60ch] text-sm text-muted text-pretty">{t('role.empty.body')}</p>
{/if}

<!-- Delete confirm -->
<dialog bind:this={deleteDialog} class="modal" aria-labelledby="role-delete-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="role-delete-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('role.delete.title')}
		</h2>
		{#if pending}
			<p class="mt-1 text-sm text-muted text-pretty">
				{t('role.delete.warning', { name: roleDisplayName(pending.name) })}
			</p>
		{/if}

		{#if deleteMessage}
			<div class="mt-4"><Alert align="start">{deleteMessage}</Alert></div>
		{/if}

		<form
			method="POST"
			action="?/delete"
			use:enhance={submitDelete}
			class="mt-5 flex justify-end gap-2"
		>
			<input type="hidden" name="roleId" value={pending?.id ?? ''} />
			<Button type="button" variant="ghost" onclick={() => deleteDialog?.close()}>
				{t('role.cancel')}
			</Button>
			<Button type="submit" variant="danger" loading={deleteSubmitting}>
				{deleteSubmitting ? t('role.delete.submitting') : t('role.delete.submit')}
			</Button>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('role.cancel')}></button>
	</form>
</dialog>

{#if toastMsg}
	<div class="pointer-events-none fixed inset-x-0 bottom-4 z-50 flex justify-center px-4">
		<div class="pointer-events-auto"><Alert variant={toastVariant}>{toastMsg}</Alert></div>
	</div>
{/if}
