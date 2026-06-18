<script lang="ts">
	import { Field } from '$lib/components/common';
	import { t } from '$lib/i18n';
	import { groupPermissions, type PermissionGroup } from '$lib/access/permissions';

	type Props = {
		catalog: string[];
		name?: string;
		selected?: string[];
		/** Read-only mode for system roles: no name field, checkboxes locked. */
		disabled?: boolean;
		nameError?: string;
	};

	let {
		catalog,
		name = $bindable(''),
		selected = $bindable<string[]>([]),
		disabled = false,
		nameError
	}: Props = $props();

	const groups = $derived(groupPermissions(catalog));
	const selectedSet = $derived(new Set(selected));

	const groupValues = (g: PermissionGroup) => g.items.map((i) => i.value);
	const allChecked = (g: PermissionGroup) => groupValues(g).every((v) => selectedSet.has(v));
	const someChecked = (g: PermissionGroup) => groupValues(g).some((v) => selectedSet.has(v));

	function toggleGroup(g: PermissionGroup, on: boolean) {
		const values = groupValues(g);
		if (on) selected = [...selected, ...values.filter((v) => !selected.includes(v))];
		else selected = selected.filter((v) => !values.includes(v));
	}
</script>

<div class="flex flex-col gap-6">
	{#if disabled}
		<div class="flex flex-col gap-1.5">
			<span class="text-sm font-medium">{t('role.field.name')}</span>
			<p class="text-[0.9375rem]">{name}</p>
		</div>
	{:else}
		<Field
			id="role-name"
			name="name"
			label={t('role.field.name')}
			bind:value={name}
			placeholder={t('role.field.namePlaceholder')}
			error={nameError}
			required
			maxlength={60}
			autofocus
		/>
	{/if}

	<fieldset class="flex flex-col gap-1.5">
		<div class="flex items-baseline justify-between gap-3">
			<legend class="text-sm font-medium">{t('role.field.permissions')}</legend>
			<span class="font-mono text-xs text-muted">{t('role.selected', { n: selected.length })}</span>
		</div>
		{#if !disabled}
			<p class="text-sm text-muted">{t('role.field.permissionsHint')}</p>
		{/if}

		<div class="mt-2 rounded-box border border-base-content/10">
			{#each groups as g, i (g.resource)}
				<div class="p-4 {i > 0 ? 'border-t border-base-content/10' : ''}">
					<div class="flex items-center justify-between gap-3">
						<span class="text-sm font-semibold">{g.label}</span>
						<label
							class="inline-flex cursor-pointer items-center gap-2 text-xs text-muted hover:text-base-content"
						>
							<input
								type="checkbox"
								class="checkbox checkbox-xs"
								checked={allChecked(g)}
								indeterminate={someChecked(g) && !allChecked(g)}
								{disabled}
								onchange={(e) => toggleGroup(g, e.currentTarget.checked)}
							/>
							{t('role.selectAll')}
						</label>
					</div>

					<div class="mt-3 grid grid-cols-2 gap-x-4 gap-y-2.5 sm:grid-cols-3">
						{#each g.items as item (item.value)}
							<label class="inline-flex items-center gap-2 text-sm {disabled ? '' : 'cursor-pointer'}">
								<input
									type="checkbox"
									name="permissions"
									value={item.value}
									bind:group={selected}
									{disabled}
									class="checkbox checkbox-sm"
								/>
								{item.label}
							</label>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	</fieldset>
</div>
