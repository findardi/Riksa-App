<script lang="ts">
	import { page } from '$app/state';
	import { roleDescription, roleDisplayName } from '$lib/access/permissions';
	import { t } from '$lib/i18n';
	import type { WorkspaceRoleData } from '$lib/types/workspace';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const roles = $derived(data.roles);
	const base = $derived(`/workspace/${page.params.slug}/management-access/role`);

	const permLabel = (r: WorkspaceRoleData) =>
		r.permissions.length === 0
			? t('role.permNone')
			: t('role.permCount', { n: r.permissions.length });
</script>

<svelte:head><title>{t('ma.role')} · {t('ma.title')}</title></svelte:head>

<div class="flex items-center justify-between gap-4">
	<h2 class="text-sm font-semibold">
		{t('ma.role')}
		<span class="ml-1 font-mono text-xs font-normal text-muted">{roles.length}</span>
	</h2>
</div>

<p class="mt-1 max-w-[60ch] text-sm text-muted text-pretty">{t('role.fixedNote')}</p>

<ul class="mt-4 divide-y divide-base-content/10 border-y border-base-content/10">
	{#each roles as role (role.id)}
		{@const desc = roleDescription(role.name)}
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
				{#if desc}
					<p class="mt-0.5 max-w-[60ch] text-sm text-muted text-pretty">{desc}</p>
				{/if}
				<p class="mt-0.5 font-mono text-xs text-muted">{permLabel(role)}</p>
			</div>

			<a
				href={`${base}/${role.id}`}
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
		</li>
	{/each}
</ul>
