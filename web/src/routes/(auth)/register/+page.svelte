<script lang="ts">
	import { enhance } from '$app/forms';
	import Field from '$lib/components/Field.svelte';
	import PasswordField from '$lib/components/PasswordField.svelte';
	import Button from '$lib/components/Button.svelte';
	import Alert from '$lib/components/Alert.svelte';
	import { t } from '$lib/i18n';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();

	let email = $state('');
	let username = $state('');
	let unlocked = $state(false);
	let checking = $state(false);
	let submitting = $state(false);
</script>

<svelte:head><title>{t('register.title')} · Wadi</title></svelte:head>

<section class="flex flex-col gap-6 text-center">
	<header>
		<h1 class="text-[1.625rem] font-semibold tracking-[-0.02em] text-balance">
			{t('register.title')}
		</h1>
		<p class="mt-1.5 text-[0.9375rem] text-muted">{t('register.subtitle')}</p>
	</header>

	{#if form?.message}
		<Alert variant="error">{form.message}</Alert>
	{/if}

	{#if !unlocked}
		<!-- Step 1 — email only -->
		<form
			method="POST"
			action="?/check"
			novalidate
			class="flex flex-col gap-[1.1rem]"
			use:enhance={() => {
				checking = true;
				return async ({ result, update }) => {
					if (result.type === 'success' && result.data?.available === true) unlocked = true;
					await update({ reset: false });
					checking = false;
				};
			}}
		>
			<Field
				id="email"
				name="email"
				type="email"
				label={t('register.email')}
				autocomplete="email"
				inputmode="email"
				bind:value={email}
				error={form?.fieldErrors?.email}
			/>
			<Button type="submit" full loading={checking}>
				{checking ? t('register.checking') : t('register.emailContinue')}
			</Button>
		</form>
	{:else}
		<!-- Step 2 — locked email + username + password -->
		<div class="flex flex-col gap-1.5">
			<div
				class="flex items-center justify-between gap-3 rounded-(--radius-field) border border-base-content/10 bg-base-100 px-3 py-2.5"
			>
				<span class="truncate font-mono text-sm">{email}</span>
				<button
					type="button"
					class="flex-none cursor-pointer text-sm font-medium text-primary hover:underline"
					onclick={() => (unlocked = false)}
				>
					{t('register.changeEmail')}
				</button>
			</div>
			{#if form?.fieldErrors?.email}
				<p class="text-sm text-error">{form.fieldErrors.email}</p>
			{:else}
				<p class="text-sm text-muted">{t('register.emailOk')}</p>
			{/if}
		</div>

		<form
			method="POST"
			action="?/register"
			novalidate
			class="flex flex-col gap-[1.1rem]"
			use:enhance={() => {
				submitting = true;
				return async ({ update }) => {
					await update({ reset: false });
					submitting = false;
				};
			}}
		>
			<input type="hidden" name="email" value={email} />
			<Field
				id="username"
				name="username"
				label={t('register.username')}
				autocomplete="username"
				hint={t('register.usernameHint')}
				bind:value={username}
				error={form?.fieldErrors?.username}
			/>
			<PasswordField
				id="password"
				name="password"
				label={t('register.password')}
				autocomplete="new-password"
				hint={t('register.passwordHint')}
				error={form?.fieldErrors?.password}
			/>
			<Button type="submit" full loading={submitting}>
				{submitting ? t('register.submitting') : t('register.submit')}
			</Button>
		</form>
	{/if}

	<div class="flex flex-col items-center gap-2 text-center">
		<p class="text-[0.9375rem] text-muted">
			{t('nav.toLogin')}
			<a href="/login" class="font-medium text-primary hover:underline">{t('nav.toLoginCta')}</a>
		</p>
	</div>
</section>
