<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		type?: 'button' | 'submit';
		variant?: 'primary' | 'ghost';
		loading?: boolean;
		disabled?: boolean;
		full?: boolean;
		onclick?: () => void;
		children: Snippet;
	};

	let {
		type = 'button',
		variant = 'primary',
		loading = false,
		disabled = false,
		full = false,
		onclick,
		children
	}: Props = $props();

	const variantClass = $derived(variant === 'primary' ? 'btn-primary' : 'btn-ghost');
</script>

<button
	{type}
	{onclick}
	class="btn {variantClass}"
	class:btn-block={full}
	disabled={disabled || loading}
	aria-busy={loading}
>
	{#if loading}<span class="loading loading-spinner loading-sm"></span>{/if}
	{@render children()}
</button>
