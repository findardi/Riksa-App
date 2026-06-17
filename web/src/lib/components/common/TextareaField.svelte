<script lang="ts">
	// Textarea counterpart of Field — same label/error/hint + a11y contract.
	type Props = {
		id: string;
		name: string;
		label: string;
		value?: string;
		error?: string;
		hint?: string;
		placeholder?: string;
		required?: boolean;
		maxlength?: number;
		rows?: number;
	};

	let {
		id,
		name,
		label,
		value = $bindable(''),
		error,
		hint,
		placeholder,
		required = false,
		maxlength,
		rows = 3
	}: Props = $props();

	const describedBy = $derived(
		[error ? `${id}-error` : null, hint && !error ? `${id}-hint` : null]
			.filter(Boolean)
			.join(' ') || undefined
	);
</script>

<div class="flex flex-col gap-1.5">
	<label class="text-sm font-medium" for={id}>{label}</label>
	<textarea
		{id}
		{name}
		{placeholder}
		{required}
		{maxlength}
		{rows}
		bind:value
		class="textarea w-full focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
		class:textarea-error={!!error}
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={describedBy}
	></textarea>
	{#if hint && !error}<p class="text-sm text-muted" id="{id}-hint">{hint}</p>{/if}
	{#if error}<p class="text-sm text-error" id="{id}-error">{error}</p>{/if}
</div>
