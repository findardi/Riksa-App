<script lang="ts">
	import { LOCALES, localeLabels, t, type Locale } from '$lib/i18n';
	import { localeState, setLocale } from '$lib/i18n/locale.svelte';

	// 'dropdown' is a standalone globe trigger; 'inline' is a bare option list
	// meant to sit inside another menu (e.g. the topbar account dropdown).
	type Props = { variant?: 'dropdown' | 'inline' };
	let { variant = 'dropdown' }: Props = $props();

	function choose(locale: Locale) {
		setLocale(locale);
		// Close any daisyUI dropdown this lives inside.
		(document.activeElement as HTMLElement | null)?.blur();
	}
</script>

{#snippet options()}
	{#each LOCALES as locale (locale)}
		<button
			type="button"
			onclick={() => choose(locale)}
			class="flex w-full items-center justify-between gap-2 rounded-field px-3 py-2 text-left text-sm hover:bg-base-content/5 {localeState.current ===
			locale
				? 'font-medium text-base-content'
				: 'text-muted'}"
		>
			{localeLabels[locale]}
			{#if localeState.current === locale}
				<svg
					class="h-4 w-4 flex-none text-primary"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.8"
					stroke-linecap="round"
					stroke-linejoin="round"
					aria-hidden="true"
				>
					<path d="m5 13 4 4L19 7" />
				</svg>
			{/if}
		</button>
	{/each}
{/snippet}

{#if variant === 'inline'}
	<p class="px-3 pt-1 pb-0.5 text-xs text-muted">{t('app.language')}</p>
	{@render options()}
{:else}
	<div class="dropdown dropdown-end">
		<button
			tabindex="0"
			type="button"
			class="flex items-center gap-1.5 rounded-field border border-base-content/10 bg-base-100 px-2.5 py-1.5 text-sm text-base-content hover:bg-base-content/5"
			aria-label={t('app.language')}
		>
			<svg
				class="h-4 w-4 flex-none text-muted"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.6"
				stroke-linecap="round"
				stroke-linejoin="round"
				aria-hidden="true"
			>
				<circle cx="12" cy="12" r="9" />
				<path d="M3 12h18" />
				<path d="M12 3a15 15 0 0 1 0 18a15 15 0 0 1 0-18" />
			</svg>
			<span>{localeLabels[localeState.current]}</span>
		</button>
		<ul
			class="dropdown-content z-50 mt-2 w-44 rounded-box border border-base-content/10 bg-base-100 p-1.5 shadow-lg"
		>
			{@render options()}
		</ul>
	</div>
{/if}
