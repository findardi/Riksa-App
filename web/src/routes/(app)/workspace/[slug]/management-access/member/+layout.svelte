<script lang="ts">
	import { page } from '$app/state';
	import { t } from '$lib/i18n';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps = $props();

	const base = $derived(`/workspace/${page.params.slug}/management-access/member`);
	const subtabs = $derived([
		{ href: base, label: t('ma.member'), count: data.members.length, exact: true },
		{ href: `${base}/invite`, label: t('ma.pending'), count: data.pendingCount, exact: false }
	]);
	const isActive = (href: string, exact: boolean) =>
		exact ? page.url.pathname === href : page.url.pathname.startsWith(href);
</script>

<!-- Segmented control: secondary nav, deliberately distinct from the underline
	 tabs above it so the two rows never read as the same control. -->
<nav
	class="inline-flex gap-1 rounded-field bg-base-content/4 p-1"
	aria-label={t('ma.member')}
>
	{#each subtabs as tab (tab.href)}
		{@const active = isActive(tab.href, tab.exact)}
		<a
			href={tab.href}
			aria-current={active ? 'page' : undefined}
			class="inline-flex items-center gap-2 rounded-selector px-3 py-1.5 text-sm font-medium transition-colors {active
				? 'bg-base-100 text-base-content shadow-sm'
				: 'text-muted hover:text-base-content'}"
		>
			{tab.label}
			<span
				class="font-mono text-xs {active ? 'text-primary' : 'text-muted'}"
				aria-hidden="true">{tab.count}</span
			>
		</a>
	{/each}
</nav>

<div class="mt-5">
	{@render children()}
</div>
