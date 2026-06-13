<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		variant?: 'error' | 'success';
		children: Snippet;
	};
	let { variant = 'error', children }: Props = $props();

	const variantClass = $derived(variant === 'error' ? 'alert-error' : 'alert-success');
</script>

<div
	role="alert"
	class="wadi-alert-in alert alert-soft {variantClass} justify-center text-center text-sm"
>
	{@render children()}
</div>

<style>
	/* State feedback enters with a subtle settle; instant for reduced motion. */
	.wadi-alert-in {
		animation: wadi-alert-in 180ms ease-out;
	}
	@keyframes wadi-alert-in {
		from {
			opacity: 0;
			transform: translateY(-2px);
		}
		to {
			opacity: 1;
			transform: none;
		}
	}
	@media (prefers-reduced-motion: reduce) {
		.wadi-alert-in {
			animation: none;
		}
	}
</style>
