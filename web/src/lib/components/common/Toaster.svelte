<script lang="ts">
	import Alert from './Alert.svelte';
	import { store } from './toast.svelte';
</script>

{#if store.current}
	<div class="pointer-events-none fixed inset-x-0 bottom-4 z-toast flex justify-center px-4">
		<div class="pointer-events-auto">
			<Alert variant={store.current.variant}>{store.current.message}</Alert>
		</div>
	</div>
{/if}
