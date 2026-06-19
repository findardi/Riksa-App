<script lang="ts">
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button } from '$lib/components/common';
	import { roleDisplayName } from '$lib/access/permissions';
	import { t } from '$lib/i18n';
	import type { MemberStatus, WorkspaceMemberData } from '$lib/types/workspace';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const members = $derived(data.members);
	// Don't offer "owner" as an assignable role — one owner per room.
	const roleOptions = $derived(data.roles.filter((r) => r.name !== 'owner'));

	const isOwner = (m: WorkspaceMemberData) => m.user_id === data.ownerId;
	const initial = (m: WorkspaceMemberData) => (m.username || m.email || '?').charAt(0).toUpperCase();

	const statusMeta = (s: MemberStatus) => {
		if (s === 'active') return { label: t('member.status.active'), dot: 'bg-success' };
		if (s === 'invited') return { label: t('member.status.invited'), dot: 'bg-warning' };
		if (s === 'suspended') return { label: t('member.status.suspended'), dot: 'bg-error' };
		return { label: s, dot: 'bg-base-content/30' };
	};

	// Per-member select value: a writable derived seeded from server truth. The
	// select reassigns it on change; any reload re-derives, so a failed change
	// visibly reverts to the stored role.
	let choice = $derived(Object.fromEntries(data.members.map((m) => [m.id, m.role_id])));

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

	// --- Change role (inline select, auto-submit) ---
	let savingId = $state<string | null>(null);
	const submitRole = (memberId: string): SubmitFunction => {
		return () => {
			savingId = memberId;
			return async ({ result }) => {
				savingId = null;
				await invalidateAll(); // success applies the new role; failure reverts the select
				if (result.type === 'success') showToast(t('member.roleChanged'), 'success');
				else if (result.type === 'failure')
					showToast((result.data?.message as string) ?? t('err.generic'), 'error');
				else showToast(t('err.generic'), 'error');
			};
		};
	};

	// --- Remove ---
	let removeDialog = $state<HTMLDialogElement>();
	let pending = $state<WorkspaceMemberData | null>(null);
	let removeSubmitting = $state(false);
	let removeMessage = $state<string | null>(null);

	function openRemove(m: WorkspaceMemberData) {
		pending = m;
		removeMessage = null;
		removeDialog?.showModal();
	}

	const submitRemove: SubmitFunction = () => {
		removeSubmitting = true;
		return async ({ result }) => {
			removeSubmitting = false;
			if (result.type === 'success') {
				removeDialog?.close();
				await invalidateAll();
				showToast(t('member.removed'), 'success');
			} else if (result.type === 'failure') {
				removeMessage = (result.data?.message as string) ?? t('err.generic');
			} else {
				removeMessage = t('err.generic');
			}
		};
	};
</script>

<svelte:head><title>{t('ma.member')} · {t('ma.title')}</title></svelte:head>

<div class="flex items-center justify-between gap-4">
	<h2 class="text-sm font-semibold">
		{t('ma.member')}
		<span class="ml-1 font-mono text-xs font-normal text-muted">{members.length}</span>
	</h2>
	<!-- Invite flow not built yet: present but inert. -->
	<button type="button" disabled class="btn btn-primary btn-sm" title={t('app.nav.soon')}>
		{t('member.invite')}
		<span class="text-[0.6875rem] font-normal opacity-80">· {t('app.nav.soon')}</span>
	</button>
</div>

<ul class="mt-4 divide-y divide-base-content/10 border-y border-base-content/10">
	{#each members as m (m.id)}
		{@const status = statusMeta(m.status)}
		{@const owner = isOwner(m)}
		{@const groups = m.group_names ?? []}
		<li class="flex items-center gap-3 py-3">
			<span
				class="grid h-9 w-9 flex-none place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary"
				aria-hidden="true">{initial(m)}</span
			>

			<div class="min-w-0 flex-1">
				<div class="flex items-center gap-2">
					<span class="truncate text-[0.9375rem] font-medium">{m.username || m.email}</span>
					{#if owner}
						<span
							class="rounded-selector bg-base-content/10 px-1.5 py-0.5 text-[0.6875rem] font-medium text-muted"
							>{roleDisplayName('owner')}</span
						>
					{/if}
				</div>
				<p class="mt-0.5 flex flex-wrap items-center gap-x-2 gap-y-0.5 text-xs text-muted">
					<span class="truncate font-mono">{m.email}</span>
					<span aria-hidden="true">·</span>
					<span class="inline-flex items-center gap-1.5">
						<span class="h-1.5 w-1.5 rounded-full {status.dot}"></span>
						{status.label}
					</span>
				</p>
				{#if groups.length}
					<div class="mt-1.5 flex flex-wrap gap-1">
						{#each groups as g (g)}
							<span class="rounded-selector bg-base-content/5 px-1.5 py-0.5 text-[0.6875rem] text-muted"
								>{g}</span
							>
						{/each}
					</div>
				{/if}
			</div>

			{#if owner}
				<span
					class="inline-flex w-36 flex-none items-center justify-end gap-1.5 text-sm text-muted"
					title={t('member.role.locked')}
				>
					{roleDisplayName(m.role_name)}
					<svg
						class="h-3.5 w-3.5 flex-none"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.8"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<rect x="5" y="11" width="14" height="9" rx="2" />
						<path d="M8 11V7a4 4 0 0 1 8 0v4" />
					</svg>
				</span>
			{:else}
				<form method="POST" action="?/updateRole" use:enhance={submitRole(m.id)} class="contents">
					<input type="hidden" name="memberId" value={m.id} />
					<select
						name="roleId"
						bind:value={choice[m.id]}
						onchange={(e) => e.currentTarget.form?.requestSubmit()}
						disabled={savingId === m.id}
						aria-label={t('member.changeRole', { name: m.username || m.email })}
						class="select select-sm w-36 flex-none"
					>
						{#each roleOptions as r (r.id)}
							<option value={r.id}>{roleDisplayName(r.name)}</option>
						{/each}
					</select>
				</form>
				<button
					type="button"
					onclick={() => openRemove(m)}
					aria-label={t('member.remove')}
					title={t('member.remove')}
					class="flex-none rounded-field p-2 text-muted transition-colors hover:bg-error/10 hover:text-error"
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
		</li>
	{/each}
</ul>

<!-- Remove confirm -->
<dialog bind:this={removeDialog} class="modal" aria-labelledby="member-remove-title">
	<div class="modal-box w-full max-w-md rounded-box border border-base-content/10 bg-base-100 p-6">
		<h2 id="member-remove-title" class="text-lg font-semibold tracking-[-0.01em]">
			{t('member.remove.title')}
		</h2>
		{#if pending}
			<p class="mt-1 text-sm text-muted text-pretty">
				{t('member.remove.warning', { name: pending.username || pending.email })}
			</p>
		{/if}

		{#if removeMessage}
			<div class="mt-4"><Alert align="start">{removeMessage}</Alert></div>
		{/if}

		<form
			method="POST"
			action="?/delete"
			use:enhance={submitRemove}
			class="mt-5 flex justify-end gap-2"
		>
			<input type="hidden" name="memberId" value={pending?.id ?? ''} />
			<Button type="button" variant="ghost" onclick={() => removeDialog?.close()}>
				{t('member.cancel')}
			</Button>
			<Button type="submit" variant="danger" loading={removeSubmitting}>
				{removeSubmitting ? t('member.remove.submitting') : t('member.remove.submit')}
			</Button>
		</form>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button aria-label={t('member.cancel')}></button>
	</form>
</dialog>

{#if toastMsg}
	<div class="pointer-events-none fixed inset-x-0 bottom-4 z-50 flex justify-center px-4">
		<div class="pointer-events-auto"><Alert variant={toastVariant}>{toastMsg}</Alert></div>
	</div>
{/if}
