<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		variant?: 'error' | 'success';
		children: Snippet;
	};
	let { variant = 'error', children }: Props = $props();

	const variantClass = $derived(variant === 'error' ? 'alert-error' : 'alert-success');
</script>

<div role="alert" class="alert alert-soft {variantClass} justify-center text-center text-sm">
	{@render children()}
</div>
