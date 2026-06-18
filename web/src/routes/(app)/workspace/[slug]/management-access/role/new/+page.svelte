<script lang="ts">
	import { applyAction, enhance } from '$app/forms';
	import { page } from '$app/state';
	import type { SubmitFunction } from '@sveltejs/kit';
	import { Alert, Button } from '$lib/components/common';
	import { RoleForm } from '$lib/components/app';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data, form }: PageProps = $props();

	const backHref = $derived(`/workspace/${page.params.slug}/management-access/role`);

	let name = $state('');
	let selected = $state<string[]>([]);
	let submitting = $state(false);
	const canSubmit = $derived(name.trim().length > 0 && selected.length > 0);

	const submit: SubmitFunction = () => {
		submitting = true;
		return async ({ result, update }) => {
			submitting = false;
			if (result.type === 'redirect') await applyAction(result);
			else await update(); // failure → surface errors, keep what was typed
		};
	};
</script>

<svelte:head><title>{t('role.create.title')} · {t('ma.title')}</title></svelte:head>

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

<header class="mt-3">
	<h2 class="text-lg font-semibold tracking-[-0.01em]">{t('role.create.title')}</h2>
	<p class="mt-1 text-sm text-muted">{t('role.create.desc')}</p>
</header>

{#if form?.message}
	<div class="mt-4"><Alert align="start">{form.message}</Alert></div>
{/if}

<form method="POST" use:enhance={submit} class="mt-6">
	<RoleForm catalog={data.catalog} bind:name bind:selected nameError={form?.fieldErrors?.name} />

	<div class="mt-6 flex justify-end gap-2 border-t border-base-content/10 pt-5">
		<a href={backHref} class="btn btn-ghost btn-sm">{t('role.cancel')}</a>
		<Button type="submit" loading={submitting} disabled={!canSubmit}>
			{submitting ? t('role.create.submitting') : t('role.create.submit')}
		</Button>
	</div>
</form>
