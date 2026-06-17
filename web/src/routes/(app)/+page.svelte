<script lang="ts">
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const greeting = $derived(
		data.user?.username ? t('home.welcome', { name: data.user.username }) : t('home.welcomeGeneric')
	);
</script>

<svelte:head><title>{t('brand.name')}</title></svelte:head>

<div class="mx-auto w-full max-w-3xl px-6 py-10">
	<h1 class="text-2xl font-semibold tracking-[-0.02em] text-balance">{greeting}</h1>
	<p class="mt-2 max-w-[52ch] text-[0.9375rem] leading-relaxed text-muted text-pretty">
		{t('home.subtitle')}
	</p>

	<section class="mt-8">
		<h2 class="text-sm font-medium text-muted">{t('home.quickActions')}</h2>
		<div class="mt-3 flex flex-wrap gap-3">
			<a
				href="/workspace"
				class="group flex items-center gap-3 rounded-box border border-base-content/10 bg-base-100 px-4 py-3 transition-colors hover:border-primary/40 hover:bg-primary/5"
			>
				<span
					class="grid h-9 w-9 flex-none place-items-center rounded-field bg-primary/10 text-primary"
				>
					<svg
						class="h-[18px] w-[18px]"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="1.6"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<rect x="3" y="3" width="7" height="7" rx="1.5" />
						<rect x="14" y="3" width="7" height="7" rx="1.5" />
						<rect x="3" y="14" width="7" height="7" rx="1.5" />
						<rect x="14" y="14" width="7" height="7" rx="1.5" />
					</svg>
				</span>
				<span class="flex flex-col text-left">
					<span class="text-sm font-medium">{t('home.action.workspaces')}</span>
					<span class="text-xs text-muted">{t('home.action.workspacesDesc')}</span>
				</span>
				<svg
					class="ml-2 h-4 w-4 flex-none text-muted transition-colors group-hover:text-primary"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.6"
					stroke-linecap="round"
					stroke-linejoin="round"
					aria-hidden="true"
				>
					<path d="M5 12h14M13 6l6 6-6 6" />
				</svg>
			</a>
		</div>
	</section>
</div>
