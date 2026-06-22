<script lang="ts">
	import { enhance } from '$app/forms';
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button } from '$lib/components/common';
	import { InviteDialog } from '$lib/components/app';
	import { roleDisplayName } from '$lib/access/permissions';
	import { t } from '$lib/i18n';
	import type { InvitationData } from '$lib/types/workspace';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const invitations = $derived(data.invitations);

	let inviteOpen = $state(false);

	const initial = (inv: InvitationData) => (inv.email || '?').charAt(0).toUpperCase();

	const dateFmt = new Intl.DateTimeFormat('id-ID', {
		day: '2-digit',
		month: 'short',
		year: 'numeric'
	});
	const fmtDate = (iso: string) => dateFmt.format(new Date(iso));

	// Status → label + color (full class strings keep Tailwind JIT happy).
	const statusMeta = (s: string) => {
		switch (s) {
			case 'pending':
				return { label: t('inv.status.pending'), dot: 'bg-warning', avatar: 'bg-warning/10 text-warning' };
			case 'accepted':
				return { label: t('inv.status.accepted'), dot: 'bg-success', avatar: 'bg-success/10 text-success' };
			case 'expired':
				return {
					label: t('inv.status.expired'),
					dot: 'bg-base-content/40',
					avatar: 'bg-base-content/10 text-base-content/60'
				};
			case 'revoked':
				return { label: t('inv.status.revoked'), dot: 'bg-error', avatar: 'bg-error/10 text-error' };
			case 'rejected':
				return { label: t('inv.status.rejected'), dot: 'bg-error', avatar: 'bg-error/10 text-error' };
			default:
				return { label: s, dot: 'bg-base-content/30', avatar: 'bg-base-content/10 text-muted' };
		}
	};
	// Resend/revoke only make sense while an invite is still live.
	const canManage = (s: string) => s === 'pending' || s === 'expired';

	const filters = ['', 'pending', 'accepted', 'expired', 'revoked', 'rejected'];
	function onFilter(e: Event & { currentTarget: HTMLSelectElement }) {
		const v = e.currentTarget.value;
		const u = new URL(page.url);
		if (v) u.searchParams.set('status', v);
		else u.searchParams.delete('status');
		goto(u, { keepFocus: true, noScroll: true });
	}
	function clearFilter() {
		const u = new URL(page.url);
		u.searchParams.delete('status');
		goto(u, { keepFocus: true, noScroll: true });
	}

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

	// --- Resend (direct, per row) ---
	let resendingId = $state<string | null>(null);
	const submitResend = (id: string): SubmitFunction => {
		return () => {
			resendingId = id;
			return async ({ result }) => {
				resendingId = null;
				if (result.type === 'success') {
					await invalidateAll();
					showToast(t('pending.resend.done'), 'success');
				} else if (result.type === 'failure') {
					showToast((result.data?.message as string) ?? t('err.generic'), 'error');
				} else {
					showToast(t('err.generic'), 'error');
				}
			};
		};
	};

	// --- Revoke (confirm, destructive) ---
	let revokeDialog = $state<HTMLDialogElement>();
	let revokeTarget = $state<InvitationData | null>(null);
	let revokeSubmitting = $state(false);
	let revokeMessage = $state<string | null>(null);

	function openRevoke(inv: InvitationData) {
		revokeTarget = inv;
		revokeMessage = null;
		revokeDialog?.showModal();
	}

	const submitRevoke: SubmitFunction = () => {
		revokeSubmitting = true;
		return async ({ result }) => {
			revokeSubmitting = false;
			if (result.type === 'success') {
				revokeDialog?.close();
				await invalidateAll();
				showToast(t('pending.revoke.done'), 'success');
			} else if (result.type === 'failure') {
				revokeMessage = (result.data?.message as string) ?? t('err.generic');
			} else {
				revokeMessage = t('err.generic');
			}
		};
	};
</script>

<svelte:head><title>{t('ma.pending')} · {t('ma.title')}</title></svelte:head>

<div class="flex items-center justify-between gap-3">
	<label class="sr-only" for="inv-filter">{t('inv.filter.label')}</label>
	<select
		id="inv-filter"
		value={data.status}
		onchange={onFilter}
		class="select select-sm w-auto"
		aria-label={t('inv.filter.label')}
	>
		{#each filters as f (f)}
			<option value={f}>{f ? statusMeta(f).label : t('inv.filter.all')}</option>
		{/each}
	</select>
	<button type="button" onclick={() => (inviteOpen = true)} class="btn btn-primary btn-sm">
		{t('member.invite')}
	</button>
</div>

{#if invitations.length}
	<ul class="mt-4 divide-y divide-base-content/10 border-y border-base-content/10">
		{#each invitations as inv (inv.id)}
			{@const meta = statusMeta(inv.status)}
			{@const manageable = canManage(inv.status)}
			{@const resending = resendingId === inv.id}
			<li class="flex items-center gap-3 py-3">
				<span
					class="grid h-9 w-9 flex-none place-items-center rounded-full text-sm font-semibold {meta.avatar}"
					aria-hidden="true">{initial(inv)}</span
				>

				<div class="min-w-0 flex-1">
					<div class="flex items-center gap-2">
						<span class="truncate font-mono text-[0.9375rem] font-medium">{inv.email}</span>
						<span
							class="rounded-selector bg-base-content/10 px-1.5 py-0.5 text-[0.6875rem] font-medium text-muted"
							>{roleDisplayName(inv.role_name)}</span
						>
					</div>
					<p class="mt-0.5 flex flex-wrap items-center gap-x-2 gap-y-0.5 text-xs text-muted">
						<span class="inline-flex items-center gap-1.5">
							<span class="h-1.5 w-1.5 rounded-full {meta.dot}"></span>
							{meta.label}
						</span>
						<span aria-hidden="true">·</span>
						<span>
							{t('pending.expires')}
							<span class="font-mono">{fmtDate(inv.expires_at)}</span>
						</span>
						{#if inv.invited_by_username}
							<span aria-hidden="true">·</span>
							<span>{t('pending.invitedBy', { name: inv.invited_by_username })}</span>
						{/if}
					</p>
				</div>

				<div class="flex flex-none items-center gap-1">
					<form method="POST" action="?/resend" use:enhance={submitResend(inv.id)} class="contents">
						<input type="hidden" name="invitationId" value={inv.id} />
						<button
							type="submit"
							disabled={!manageable || resending}
							aria-label={t('pending.resend')}
							title={t('pending.resend')}
							class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-primary/10 hover:text-primary disabled:pointer-events-none disabled:opacity-30"
						>
							{#if resending}
								<span class="loading loading-spinner loading-xs"></span>
							{:else}
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
									<path d="M21 2v6h-6" />
									<path d="M3 12a9 9 0 0 1 15-6.7L21 8" />
									<path d="M3 22v-6h6" />
									<path d="M21 12a9 9 0 0 1-15 6.7L3 16" />
								</svg>
							{/if}
						</button>
					</form>
					<button
						type="button"
						onclick={() => openRevoke(inv)}
						disabled={!manageable}
						aria-label={t('pending.revoke')}
						title={t('pending.revoke')}
						class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-error/10 hover:text-error disabled:pointer-events-none disabled:opacity-30"
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
							<circle cx="12" cy="12" r="9" />
							<path d="m15 9-6 6M9 9l6 6" />
						</svg>
					</button>
				</div>
			</li>
		{/each}
	</ul>
{:else if data.status}
	<!-- Empty because of an active filter, not because there are no invites. -->
	<div
		class="mt-4 flex flex-col items-center gap-3 rounded-box border border-dashed border-base-content/15 px-6 py-12 text-center"
	>
		<p class="text-sm text-muted">{t('inv.empty.filtered')}</p>
		<button type="button" onclick={clearFilter} class="text-sm font-medium text-primary hover:underline">
			{t('inv.empty.reset')}
		</button>
	</div>
{:else}
	<div
		class="mt-4 grid place-items-center rounded-box border border-dashed border-base-content/15 px-6 py-14 text-center"
	>
		<span
			class="grid h-11 w-11 place-items-center rounded-full bg-base-content/5 text-muted"
			aria-hidden="true"
		>
			<svg
				class="h-5 w-5"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.6"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<rect x="3" y="5" width="18" height="14" rx="2" />
				<path d="m3 7 9 6 9-6" />
			</svg>
		</span>
		<h3 class="mt-3 text-sm font-semibold">{t('pending.empty.title')}</h3>
		<p class="mt-1 max-w-sm text-sm text-muted text-pretty">{t('pending.empty.desc')}</p>
		<button
			type="button"
			onclick={() => (inviteOpen = true)}
			class="btn btn-primary btn-sm mt-4">{t('member.invite')}</button
		>
	</div>
{/if}

<!-- Revoke confirm -->
<dialog bind:this={revokeDialog} class="modal" aria-labelledby="invite-revoke-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="invite-revoke-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('pending.revoke.title')}
		</h2>
		{#if revokeTarget}
			<p class="mt-1 text-sm text-muted text-pretty">
				{t('pending.revoke.warning', { email: revokeTarget.email })}
			</p>
		{/if}

		{#if revokeMessage}
			<div class="mt-4"><Alert align="start">{revokeMessage}</Alert></div>
		{/if}

		<form
			method="POST"
			action="?/revoke"
			use:enhance={submitRevoke}
			class="mt-5 flex justify-end gap-2"
		>
			<input type="hidden" name="invitationId" value={revokeTarget?.id ?? ''} />
			<Button type="button" variant="ghost" onclick={() => revokeDialog?.close()}>
				{t('member.cancel')}
			</Button>
			<Button type="submit" variant="danger" loading={revokeSubmitting}>
				{revokeSubmitting ? t('pending.revoke.submitting') : t('pending.revoke.submit')}
			</Button>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('member.cancel')}></button>
	</form>
</dialog>

<InviteDialog bind:open={inviteOpen} roles={data.roles} action="?/invite" />

{#if toastMsg}
	<div class="pointer-events-none fixed inset-x-0 bottom-4 z-50 flex justify-center px-4">
		<div class="pointer-events-auto"><Alert variant={toastVariant}>{toastMsg}</Alert></div>
	</div>
{/if}
