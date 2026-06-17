<script lang="ts">
	import { Button } from '$lib/components/common';
	import { t } from '$lib/i18n';
</script>

<svelte:head><title>{t('app.nav.rooms')} · {t('brand.name')}</title></svelte:head>

<div class="grid min-h-full place-items-center px-6 py-16">
	<section class="flex max-w-md flex-col items-center text-center">
		<span
			class="mb-5 grid h-14 w-14 place-items-center rounded-box border border-base-content/10 bg-base-100 text-muted"
		>
			<svg
				class="h-6 w-6"
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

		<h1 class="text-[1.375rem] font-semibold tracking-[-0.015em] text-balance">
			{t('app.empty.title')}
		</h1>
		<p class="mt-2 max-w-[42ch] text-[0.9375rem] leading-relaxed text-muted text-pretty">
			{t('app.empty.body')}
		</p>

		<div class="mt-6">
			<Button disabled>{t('app.empty.create')}</Button>
		</div>

		<p class="mt-8 max-w-[40ch] text-sm text-muted">{t('app.empty.guestHint')}</p>
	</section>
</div>
