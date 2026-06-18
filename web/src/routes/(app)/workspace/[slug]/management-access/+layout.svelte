<script lang="ts">
	import { page } from '$app/state';
	import { t } from '$lib/i18n';
	import type { LayoutProps } from './$types';

	let { children }: LayoutProps = $props();

	const base = $derived(`/workspace/${page.params.slug}/management-access`);
	const tabs = $derived([
		{ href: `${base}/member`, label: t('ma.member') },
		{ href: `${base}/role`, label: t('ma.role') },
		{ href: `${base}/group`, label: t('ma.group') }
	]);
	const isActive = (href: string) => page.url.pathname === href;
</script>

<div class="mx-auto w-full max-w-4xl px-6 py-8">
	<header>
		<h1 class="text-2xl font-semibold tracking-[-0.02em]">{t('ma.title')}</h1>
		<p class="mt-1.5 text-sm text-muted">{t('ma.desc')}</p>
	</header>

	<!-- Underline tabs: flat, route-based; active item carries the primary voice. -->
	<nav class="mt-6 flex gap-6 border-b border-base-content/10" aria-label={t('ma.title')}>
		{#each tabs as tab (tab.href)}
			<a
				href={tab.href}
				aria-current={isActive(tab.href) ? 'page' : undefined}
				class="-mb-px border-b-2 px-0.5 pb-2.5 text-sm font-medium transition-colors {isActive(
					tab.href
				)
					? 'border-primary text-primary'
					: 'border-transparent text-muted hover:text-base-content'}"
			>
				{tab.label}
			</a>
		{/each}
	</nav>

	<section class="mt-6">
		{@render children()}
	</section>
</div>
