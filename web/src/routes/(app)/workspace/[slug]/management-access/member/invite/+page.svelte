<script lang="ts">
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
</script>

<svelte:head><title>{t('ma.pending')} · {t('ma.title')}</title></svelte:head>

<div class="flex justify-end">
	<button type="button" onclick={() => (inviteOpen = true)} class="btn btn-primary btn-sm">
		{t('member.invite')}
	</button>
</div>

{#if invitations.length}
	<ul class="mt-4 divide-y divide-base-content/10 border-y border-base-content/10">
		{#each invitations as inv (inv.id)}
			<li class="flex items-center gap-3 py-3">
				<span
					class="grid h-9 w-9 flex-none place-items-center rounded-full bg-warning/10 text-sm font-semibold text-warning"
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
							<span class="h-1.5 w-1.5 rounded-full bg-warning"></span>
							{t('pending.status.pending')}
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

				<!-- Resend / revoke: affordance shown; backend handlers not wired yet. -->
				<div class="flex flex-none items-center gap-1">
					<span
						class="mr-0.5 rounded-selector bg-base-content/5 px-1.5 py-0.5 text-[0.6875rem] text-muted"
						>{t('app.nav.soon')}</span
					>
					<button
						type="button"
						disabled
						aria-label={t('pending.resend')}
						title={t('pending.resend')}
						class="cursor-not-allowed rounded-field p-2 text-muted opacity-40"
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
							<path d="M21 2v6h-6" />
							<path d="M3 12a9 9 0 0 1 15-6.7L21 8" />
							<path d="M3 22v-6h6" />
							<path d="M21 12a9 9 0 0 1-15 6.7L3 16" />
						</svg>
					</button>
					<button
						type="button"
						disabled
						aria-label={t('pending.revoke')}
						title={t('pending.revoke')}
						class="cursor-not-allowed rounded-field p-2 text-muted opacity-40"
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

<InviteDialog bind:open={inviteOpen} roles={data.roles} action="?/invite" />
