<script lang="ts">
	import { page } from '$app/state';
	import { RolePermissions } from '$lib/components/app';
	import { roleDescription, roleDisplayName } from '$lib/access/permissions';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const backHref = $derived(`/workspace/${page.params.slug}/management-access/role`);
	const desc = $derived(roleDescription(data.role.name));
</script>

<svelte:head>
	<title>{roleDisplayName(data.role.name)} · {t('ma.role')}</title>
</svelte:head>

<a
	href={backHref}
	class="inline-flex items-center gap-1.5 text-xs font-medium text-muted transition-colors hover:text-base-content"
>
	<svg
		class="h-3.5 w-3.5"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="1.8"
		stroke-linecap="round"
		stroke-linejoin="round"
		aria-hidden="true"><path d="m15 6-6 6 6 6" /></svg
	>
	{t('role.back')}
</a>

<header class="mt-3 flex items-center gap-2">
	<h2 class="text-lg font-semibold tracking-[-0.01em]">{roleDisplayName(data.role.name)}</h2>
	{#if data.role.is_system}
		<span
			class="rounded-selector bg-base-content/10 px-1.5 py-0.5 text-[0.6875rem] font-medium text-muted"
			>{t('role.system')}</span
		>
	{/if}
</header>

{#if desc}
	<p class="mt-1 max-w-[60ch] text-sm text-muted text-pretty">{desc}</p>
{/if}

<p class="mt-4 rounded-box bg-base-content/5 p-3 text-sm text-muted text-pretty">
	{t('role.view.systemNote')}
</p>

<div class="mt-6">
	<RolePermissions catalog={data.catalog} granted={data.role.permissions} />
</div>

<div class="mt-6 flex justify-end border-t border-base-content/10 pt-5">
	<a href={backHref} class="btn btn-ghost btn-sm">{t('role.back')}</a>
</div>
