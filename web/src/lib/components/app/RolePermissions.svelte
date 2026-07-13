<script lang="ts">
	import { t } from '$lib/i18n';
	import { groupPermissions } from '$lib/access/permissions';

	type Props = {
		/** Full permission catalog (from getPermissions). */
		catalog: string[];
		/** Permissions granted to this role. */
		granted: string[];
	};

	let { catalog, granted }: Props = $props();

	const groups = $derived(groupPermissions(catalog));
	const grantedSet = $derived(new Set(granted));
</script>

<fieldset class="flex flex-col gap-1.5">
	<div class="flex items-baseline justify-between gap-3">
		<legend class="text-sm font-medium">{t('role.field.permissions')}</legend>
		<span class="font-mono text-xs text-muted">{t('role.permCount', { n: granted.length })}</span>
	</div>

	<div class="mt-2 rounded-box border border-base-content/10">
		{#each groups as g, i (g.resource)}
			<div class="p-4 {i > 0 ? 'border-t border-base-content/10' : ''}">
				<span class="text-sm font-semibold">{g.label}</span>

				<div class="mt-3 grid grid-cols-2 gap-x-4 gap-y-2.5 sm:grid-cols-3">
					{#each g.items as item (item.value)}
						{@const on = grantedSet.has(item.value)}
						<span class="inline-flex items-center gap-2 text-sm {on ? '' : 'text-muted/60'}">
							{#if on}
								<svg
									class="h-4 w-4 flex-none text-success"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2.2"
									stroke-linecap="round"
									stroke-linejoin="round"
									aria-hidden="true"><path d="M20 6 9 17l-5-5" /></svg
								>
							{:else}
								<svg
									class="h-4 w-4 flex-none text-base-content/25"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
									stroke-linecap="round"
									aria-hidden="true"><path d="M5 12h14" /></svg
								>
							{/if}
							{item.label}
						</span>
					{/each}
				</div>
			</div>
		{/each}
	</div>
</fieldset>
