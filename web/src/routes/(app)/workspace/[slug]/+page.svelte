<script lang="ts">
	import { tick } from 'svelte';
	import { applyAction, enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button, Field, TextareaField } from '$lib/components/common';
	import { t } from '$lib/i18n';
	import type { WorkspaceStatus } from '$lib/types/workspace';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const ws = $derived(data.workspace);

	const dateFmt = new Intl.DateTimeFormat('id-ID', {
		day: 'numeric',
		month: 'short',
		year: 'numeric'
	});
	const fmtDate = (iso: string) => dateFmt.format(new Date(iso));

	const statuses: WorkspaceStatus[] = ['prepare', 'active', 'archive'];
	const statusHint = $derived(t(`ws.status.hint.${ws.status}`));

	// --- Toast (transient feedback for in-place status changes) ---
	let toastMsg = $state<string | null>(null);
	let toastVariant = $state<'success' | 'error'>('success');
	let toastTimer: ReturnType<typeof setTimeout>;
	function showToast(msg: string, variant: 'success' | 'error') {
		toastMsg = msg;
		toastVariant = variant;
		clearTimeout(toastTimer);
		toastTimer = setTimeout(() => (toastMsg = null), 4000);
	}

	// --- Status change ---
	let pendingStatus = $state<WorkspaceStatus | null>(null);
	const submitStatus: SubmitFunction = ({ formData }) => {
		pendingStatus = formData.get('status') as WorkspaceStatus;
		return async ({ result }) => {
			if (result.type === 'success') {
				await invalidateAll();
				showToast(t('ws.status.updated'), 'success');
			} else if (result.type === 'redirect') {
				await applyAction(result);
			} else if (result.type === 'failure') {
				showToast((result.data?.message as string) ?? t('err.generic'), 'error');
			} else {
				showToast(t('err.generic'), 'error');
			}
			pendingStatus = null;
		};
	};

	// --- Edit name & description ---
	let editDialog = $state<HTMLDialogElement>();
	let editAlertEl = $state<HTMLDivElement>();
	let editName = $state('');
	let editDescription = $state('');
	let editSubmitting = $state(false);
	let editMessage = $state<string | null>(null);
	let editFieldErrors = $state<Record<string, string>>({});

	function openEdit() {
		editName = ws.name;
		editDescription = ws.description ?? '';
		editMessage = null;
		editFieldErrors = {};
		editDialog?.showModal();
		tick().then(() => editDialog?.querySelector<HTMLInputElement>('#edit-name')?.focus());
	}

	function focusEditError() {
		if (editFieldErrors.name) editDialog?.querySelector<HTMLInputElement>('#edit-name')?.focus();
		else if (editMessage) editAlertEl?.focus();
	}

	const submitEdit: SubmitFunction = () => {
		editSubmitting = true;
		return async ({ result }) => {
			editSubmitting = false;
			if (result.type === 'success') {
				editDialog?.close();
				await invalidateAll();
				showToast(t('ws.edit.saved'), 'success');
			} else if (result.type === 'redirect') {
				editDialog?.close();
				await applyAction(result); // reslugged — lands on the new slug
			} else if (result.type === 'failure') {
				editMessage = (result.data?.message as string | null) ?? null;
				editFieldErrors = (result.data?.fieldErrors as Record<string, string>) ?? {};
				await tick();
				focusEditError();
			} else {
				editMessage = t('err.generic');
				await tick();
				focusEditError();
			}
		};
	};

	// --- Delete (type-to-confirm) ---
	let deleteDialog = $state<HTMLDialogElement>();
	let deleteConfirm = $state('');
	let deleteSubmitting = $state(false);
	let deleteMessage = $state<string | null>(null);
	const deleteReady = $derived(deleteConfirm.trim() === ws.name);

	function openDelete() {
		deleteConfirm = '';
		deleteMessage = null;
		deleteDialog?.showModal();
	}

	const submitDelete: SubmitFunction = ({ cancel }) => {
		if (!deleteReady) return cancel();
		deleteSubmitting = true;
		return async ({ result }) => {
			deleteSubmitting = false;
			if (result.type === 'redirect') {
				deleteDialog?.close();
				await applyAction(result); // → /workspace
			} else if (result.type === 'failure') {
				deleteMessage = (result.data?.message as string) ?? t('err.generic');
			} else {
				deleteMessage = t('err.generic');
			}
		};
	};
</script>

<svelte:head><title>{ws.name} · {t('brand.name')}</title></svelte:head>

<div class="mx-auto w-full max-w-2xl px-6 py-8">
	<header class="flex items-start justify-between gap-4">
		<div class="min-w-0">
			<h1 class="truncate text-2xl font-semibold tracking-[-0.02em]">{ws.name}</h1>
			<p class="mt-1.5 flex flex-wrap items-center gap-x-2 gap-y-0.5 font-mono text-xs text-muted">
				<span>{ws.slug}</span>
				<span aria-hidden="true">·</span>
				<span>{t('ws.detail.created')} {fmtDate(ws.created_at)}</span>
				<span aria-hidden="true">·</span>
				<span>{t('ws.detail.updated')} {fmtDate(ws.updated_at)}</span>
			</p>
		</div>
		<Button variant="ghost" onclick={openEdit}>{t('ws.edit.open')}</Button>
	</header>

	{#if ws.description}
		<p class="mt-5 max-w-[65ch] text-[0.9375rem] leading-relaxed text-muted text-pretty">
			{ws.description}
		</p>
	{/if}

	<!-- Status -->
	<section class="mt-10">
		<h2 id="status-label" class="text-sm font-semibold">{t('ws.status.label')}</h2>
		<p class="mt-1 text-sm text-muted">{statusHint}</p>

		<form
			method="POST"
			action="?/updateStatus"
			use:enhance={submitStatus}
			role="group"
			aria-labelledby="status-label"
			class="mt-3 inline-flex rounded-field border border-base-content/10 p-0.5"
		>
			{#each statuses as s (s)}
				{@const current = s === ws.status}
				<button
					name="status"
					value={s}
					disabled={current || pendingStatus !== null}
					aria-pressed={current}
					class="inline-flex items-center gap-1.5 rounded-sm px-3 py-1.5 text-sm font-medium transition-colors {current
						? 'bg-primary/10 text-primary'
						: 'text-muted hover:bg-base-content/5 hover:text-base-content disabled:cursor-not-allowed disabled:text-muted/60 disabled:hover:bg-transparent disabled:hover:text-muted/60'}"
				>
					{#if pendingStatus === s}
						<span class="loading loading-spinner loading-xs"></span>
					{/if}
					{t(`ws.status.${s}`)}
				</button>
			{/each}
		</form>
	</section>

	<!-- Delete — quiet settings-style row; the red button carries the danger. -->
	<section
		class="mt-10 flex flex-col gap-3 border-t border-base-content/10 pt-6 sm:flex-row sm:items-start sm:justify-between sm:gap-6"
	>
		<div class="min-w-0">
			<h2 class="text-sm font-semibold">{t('ws.delete.title')}</h2>
			<p class="mt-1 max-w-[48ch] text-sm text-muted text-pretty">{t('ws.delete.body')}</p>
		</div>
		<div class="flex-none">
			<Button variant="danger" onclick={openDelete}>{t('ws.delete.open')}</Button>
		</div>
	</section>
</div>

<!-- Edit dialog -->
<dialog bind:this={editDialog} class="modal" aria-labelledby="edit-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="edit-title" class="text-lg font-semibold tracking-[-0.01em]">{t('ws.edit.title')}</h2>

		{#if editMessage}
			<div bind:this={editAlertEl} tabindex="-1" class="mt-4 outline-none">
				<Alert align="start">{editMessage}</Alert>
			</div>
		{/if}

		<form method="POST" action="?/update" use:enhance={submitEdit} class="mt-4 flex flex-col gap-4">
			<Field
				id="edit-name"
				name="name"
				label={t('ws.field.name')}
				bind:value={editName}
				placeholder={t('ws.field.namePlaceholder')}
				required
				maxlength={120}
				error={editFieldErrors.name}
			/>
			<TextareaField
				id="edit-description"
				name="description"
				label={t('ws.field.description')}
				bind:value={editDescription}
				placeholder={t('ws.field.descriptionPlaceholder')}
				hint={t('ws.field.descriptionHint')}
				maxlength={500}
				error={editFieldErrors.description}
			/>
			<div class="mt-2 flex justify-end gap-2">
				<Button type="button" variant="ghost" onclick={() => editDialog?.close()}>
					{t('ws.dialog.cancel')}
				</Button>
				<Button type="submit" loading={editSubmitting}>
					{editSubmitting ? t('ws.dialog.submitting') : t('ws.edit.submit')}
				</Button>
			</div>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('ws.dialog.cancel')}></button>
	</form>
</dialog>

<!-- Delete dialog -->
<dialog bind:this={deleteDialog} class="modal" aria-labelledby="delete-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="delete-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('ws.delete.title')}
		</h2>
		<p class="mt-1 text-sm text-muted text-pretty">{t('ws.delete.warning', { name: ws.name })}</p>

		{#if deleteMessage}
			<div class="mt-4"><Alert align="start">{deleteMessage}</Alert></div>
		{/if}

		<form
			method="POST"
			action="?/delete"
			use:enhance={submitDelete}
			class="mt-4 flex flex-col gap-4"
		>
			<Field
				id="delete-confirm"
				name="confirm"
				label={t('ws.delete.confirmLabel', { name: ws.name })}
				bind:value={deleteConfirm}
				placeholder={ws.name}
				autocomplete="off"
			/>
			<div class="mt-2 flex justify-end gap-2">
				<Button type="button" variant="ghost" onclick={() => deleteDialog?.close()}>
					{t('ws.dialog.cancel')}
				</Button>
				<Button type="submit" variant="danger" disabled={!deleteReady} loading={deleteSubmitting}>
					{deleteSubmitting ? t('ws.delete.submitting') : t('ws.delete.submit')}
				</Button>
			</div>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('ws.dialog.cancel')}></button>
	</form>
</dialog>

{#if toastMsg}
	<div class="pointer-events-none fixed inset-x-0 bottom-4 z-50 flex justify-center px-4">
		<div class="pointer-events-auto"><Alert variant={toastVariant}>{toastMsg}</Alert></div>
	</div>
{/if}
