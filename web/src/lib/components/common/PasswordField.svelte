<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { t } from '$lib/i18n';

	type Props = {
		id: string;
		name: string;
		label: string;
		value?: string;
		error?: string;
		hint?: string;
		autocomplete?: HTMLInputAttributes['autocomplete'];
		required?: boolean;
		autofocus?: boolean;
	};

	let {
		id,
		name,
		label,
		value = $bindable(''),
		error,
		hint,
		autocomplete = 'current-password',
		required = false,
		autofocus = false
	}: Props = $props();

	let show = $state(false);
	let inputEl = $state<HTMLInputElement>();
	$effect(() => {
		if (autofocus) inputEl?.focus();
	});

	const describedBy = $derived(
		[error ? `${id}-error` : null, hint && !error ? `${id}-hint` : null]
			.filter(Boolean)
			.join(' ') || undefined
	);
</script>

<div class="flex flex-col gap-1.5">
	<label class="text-sm font-medium" for={id}>{label}</label>
	<label
		class="input w-full focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-primary"
		class:input-error={!!error}
	>
		<input
			bind:this={inputEl}
			{id}
			{name}
			{required}
			{autocomplete}
			bind:value
			type={show ? 'text' : 'password'}
			class="grow text-left focus:outline-none"
			aria-invalid={error ? 'true' : undefined}
			aria-describedby={describedBy}
		/>
		<button
			type="button"
			class="-mr-1 cursor-pointer text-muted hover:text-base-content"
			onclick={() => (show = !show)}
			aria-label={show ? t('password.hide') : t('password.show')}
			aria-pressed={show}
		>
			{#if show}
				<svg class="size-[1.15rem]" viewBox="0 0 24 24" fill="none" aria-hidden="true">
					<path
						d="M3 3l18 18M10.6 10.7a2 2 0 002.7 2.8M9.4 5.2A9.5 9.5 0 0112 5c5 0 9 4.5 9 7a12 12 0 01-2.4 3.3M6.5 6.6A12.4 12.4 0 003 12c0 2.5 4 7 9 7 1.3 0 2.5-.3 3.6-.8"
						stroke="currentColor"
						stroke-width="1.6"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			{:else}
				<svg class="size-[1.15rem]" viewBox="0 0 24 24" fill="none" aria-hidden="true">
					<path
						d="M3 12s3.5-7 9-7 9 7 9 7-3.5 7-9 7-9-7-9-7z"
						stroke="currentColor"
						stroke-width="1.6"
						stroke-linejoin="round"
					/>
					<circle cx="12" cy="12" r="2.5" stroke="currentColor" stroke-width="1.6" />
				</svg>
			{/if}
		</button>
	</label>
	{#if hint && !error}<p class="text-sm text-muted" id="{id}-hint">{hint}</p>{/if}
	{#if error}<p class="text-sm text-error" id="{id}-error">{error}</p>{/if}
</div>
