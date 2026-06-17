<script lang="ts">
	import { page } from '$app/state';
	import { t } from '$lib/i18n';

	const isActive = (href: string) => page.url.pathname === href;
</script>

<nav class="flex h-full flex-col gap-1 p-3" aria-label="Navigasi utama">
	<a
		href="/workspace"
		class="flex items-center gap-3 rounded-field px-3 py-2 text-[0.9375rem] font-medium transition-colors {isActive(
			'/workspace'
		)
			? 'bg-primary/10 text-primary'
			: 'text-base-content hover:bg-base-content/5'}"
		aria-current={isActive('/workspace') ? 'page' : undefined}
	>
		<svg
			class="h-4.5 w-4.5 flex-none"
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
		{t('app.nav.rooms')}
	</a>

	<!-- Not built yet: present but inert — no broken links. -->
	<button
		type="button"
		disabled
		title={t('app.nav.soon')}
		class="flex cursor-not-allowed items-center gap-3 rounded-field px-3 py-2 text-[0.9375rem] font-medium text-muted/70"
	>
		<svg
			class="h-4.5 w-4.5 flex-none"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="1.6"
			stroke-linecap="round"
			stroke-linejoin="round"
			aria-hidden="true"
		>
			<rect x="3" y="5" width="18" height="14" rx="2" />
			<path d="m3 7 9 6 9-6" />
		</svg>
		<span class="flex-1 text-left">{t('app.nav.invitations')}</span>
		<span class="text-[0.6875rem] font-normal text-muted">{t('app.nav.soon')}</span>
	</button>

	<div class="flex-1"></div>

	<div class="my-1 border-t border-base-content/10"></div>

	<button
		type="button"
		disabled
		title={t('app.nav.soon')}
		class="flex cursor-not-allowed items-center gap-3 rounded-field px-3 py-2 text-[0.9375rem] font-medium text-muted/70"
	>
		<svg
			class="h-4.5 w-4.5 flex-none"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="1.6"
			stroke-linecap="round"
			stroke-linejoin="round"
			aria-hidden="true"
		>
			<line x1="4" x2="20" y1="7" y2="7" />
			<line x1="4" x2="20" y1="12" y2="12" />
			<line x1="4" x2="20" y1="17" y2="17" />
			<circle cx="9" cy="7" r="2" fill="currentColor" stroke="none" />
			<circle cx="15" cy="12" r="2" fill="currentColor" stroke="none" />
			<circle cx="9" cy="17" r="2" fill="currentColor" stroke="none" />
		</svg>
		<span class="flex-1 text-left">{t('app.nav.settings')}</span>
		<span class="text-[0.6875rem] font-normal text-muted">{t('app.nav.soon')}</span>
	</button>
</nav>
