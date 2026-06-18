<script lang="ts">
	import { untrack } from 'svelte';
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/state';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button } from '$lib/components/common';
	import { RoleForm } from '$lib/components/app';
	import { roleDisplayName } from '$lib/access/permissions';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();

	const backHref = $derived(`/workspace/${page.params.slug}/management-access/role`);
	const isSystem = $derived(data.role.is_system);

	// Seed editable state from the loaded role (initial value only; the page
	// remounts per role, so untrack keeps these out of the reactive graph).
	let name = $state(untrack(() => data.role.name));
	let selected = $state<string[]>(untrack(() => [...data.role.permissions]));

	let submitting = $state(false);
	const canSubmit = $derived(name.trim().length > 0 && selected.length > 0);

	const submit: SubmitFunction = () => {
		submitting = true;
		return async ({ result, update }) => {
			submitting = false;
			if (result.type === 'redirect') await applyAction(result);
			else await update();
		};
	};
</script>

<svelte:head>
	<title>{roleDisplayName(data.role.name)} · {t('ma.role')}</title>
</svelte:head>

<a
	href={backHref}
	class="inline-flex items-center gap-1.5 text-xs font-medium text-muted transition-colors hover:text-base-content"
>
	<svg
		class="h-3.5 w-3.5"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="1.8"
		stroke-linecap="round"
		stroke-linejoin="round"
		aria-hidden="true"><path d="m15 6-6 6 6 6" /></svg
	>
	{t('role.back')}
</a>

<header class="mt-3 flex items-center gap-2">
	<h2 class="text-lg font-semibold tracking-[-0.01em]">
		{isSystem ? t('role.view.title') : t('role.edit.title')}
	</h2>
	{#if isSystem}
		<span
			class="rounded-selector bg-base-content/10 px-1.5 py-0.5 text-[0.6875rem] font-medium text-muted"
			>{t('role.system')}</span
		>
	{/if}
</header>

{#if isSystem}
	<p class="mt-4 rounded-box bg-base-content/5 p-3 text-sm text-muted text-pretty">
		{t('role.view.systemNote')}
	</p>
	<div class="mt-6">
		<RoleForm catalog={data.catalog} {name} {selected} disabled />
	</div>
	<div class="mt-6 flex justify-end border-t border-base-content/10 pt-5">
		<a href={backHref} class="btn btn-ghost btn-sm">{t('role.back')}</a>
	</div>
{:else}
	{#if form?.message}
		<div class="mt-4"><Alert align="start">{form.message}</Alert></div>
	{/if}

	<form method="POST" use:enhance={submit} class="mt-6">
		<RoleForm catalog={data.catalog} bind:name bind:selected nameError={form?.fieldErrors?.name} />

		<div class="mt-6 flex justify-end gap-2 border-t border-base-content/10 pt-5">
			<a href={backHref} class="btn btn-ghost btn-sm">{t('role.cancel')}</a>
			<Button type="submit" loading={submitting} disabled={!canSubmit}>
				{submitting ? t('role.saving') : t('role.save')}
			</Button>
		</div>
	</form>
{/if}
