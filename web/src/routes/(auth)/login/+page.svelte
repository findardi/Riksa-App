<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import { Field, PasswordField, Button, Alert } from '$lib/components/common';
	import { t } from '$lib/i18n';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	let submitting = $state(false);

	const registered = $derived(page.url.searchParams.get('registered') === '1');
	const wasReset = $derived(page.url.searchParams.get('reset') === '1');
</script>

<svelte:head><title>{t('login.title')} · Wadi</title></svelte:head>

<section class="flex flex-col gap-6 text-center">
	<header>
		<h1 class="text-[1.625rem] font-semibold tracking-[-0.02em] text-balance">
			{t('login.title')}
		</h1>
		<p class="mt-1.5 text-[0.9375rem] text-muted">{t('login.subtitle')}</p>
	</header>

	{#if registered && !form?.message}
		<Alert variant="success">{t('login.registered')}</Alert>
	{/if}
	{#if wasReset && !form?.message}
		<Alert variant="success">{t('login.reset')}</Alert>
	{/if}
	{#if form?.message}
		<Alert variant="error">{form.message}</Alert>
	{/if}

	<form
		method="POST"
		novalidate
		class="flex flex-col gap-[1.1rem] text-left"
		use:enhance={() => {
			submitting = true;
			return async ({ update }) => {
				await update();
				submitting = false;
			};
		}}
	>
		<Field
			id="identifier"
			name="identifier"
			label={t('login.identifier')}
			autocomplete="username"
			autofocus
			value={form?.values?.identifier ?? ''}
			error={form?.fieldErrors?.identifier}
		/>
		<PasswordField
			id="password"
			name="password"
			label={t('login.password')}
			autocomplete="current-password"
			error={form?.fieldErrors?.password}
		/>
		<Button type="submit" full loading={submitting}>
			{submitting ? t('login.submitting') : t('login.submit')}
		</Button>
	</form>

	<div class="flex flex-col items-center gap-2 text-center">
		<a href="/forgot-password" class="text-sm font-medium text-primary hover:underline">
			{t('login.forgot')}
		</a>
		<p class="text-[0.9375rem] text-muted">
			{t('nav.toRegister')}
			<a href="/register" class="font-medium text-primary hover:underline"
				>{t('nav.toRegisterCta')}</a
			>
		</p>
	</div>
</section>
