<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import Field from '$lib/components/Field.svelte';
	import PasswordField from '$lib/components/PasswordField.svelte';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';
	import { t } from '$lib/i18n';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	let submitting = $state(false);

	let failedAttempts = $state(0);
	const FORGOT_THRESHOLD = 3;

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
		class="flex flex-col gap-[1.1rem]"
		use:enhance={() => {
			submitting = true;
			return async ({ result, update }) => {
				if (result.type === 'failure' && result.status === 401) failedAttempts += 1;
				else if (result.type === 'redirect' || result.type === 'success') failedAttempts = 0;
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
		<p class="text-[0.9375rem] text-muted">
			{t('nav.toRegister')}
			<a href="/register" class="font-medium text-primary hover:underline"
				>{t('nav.toRegisterCta')}</a
			>
		</p>

		{#if failedAttempts >= FORGOT_THRESHOLD}
			<p class="text-sm text-muted">
				{t('login.forgot')}
				<a href="/forgot-password" class="font-medium text-primary hover:underline"
					>{t('click.here')}</a
				>
			</p>
		{/if}
	</div>
</section>
