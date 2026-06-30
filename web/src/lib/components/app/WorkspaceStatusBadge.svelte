<script lang="ts">
	import { t } from '$lib/i18n';
	import type { WorkspaceStatus } from '$lib/types/workspace';

	type Props = { status: WorkspaceStatus; class?: string };
	let { status, class: cls = '' }: Props = $props();

	// Color carries the meaning: active lights up (success), prepare is neutral
	// and present, archive is a dimmed hollow ring (dormant, read-only).
	const config: Record<WorkspaceStatus, { dot: string; text: string }> = {
		prepare: { dot: 'bg-base-content/40', text: 'text-base-content' },
		active: { dot: 'bg-success', text: 'text-base-content' },
		archive: { dot: 'border border-base-content/40', text: 'text-muted' }
	};
	// Network data isn't type-guaranteed — degrade gracefully instead of crashing
	// the whole app shell on one unexpected value.
	const fallback = { dot: 'border border-base-content/40', text: 'text-muted' };
	const c = $derived(config[status] ?? fallback);
	const label = $derived(
		status === 'prepare'
			? t('ws.status.prepare')
			: status === 'active'
				? t('ws.status.active')
				: status === 'archive'
					? t('ws.status.archive')
					: String(status ?? '')
	);
	// What the status means — surfaced as a tooltip + accessible name so the
	// dimmed/hollow dot isn't the only thing carrying it.
	const hint = $derived(
		status === 'prepare'
			? t('ws.status.hint.prepare')
			: status === 'active'
				? t('ws.status.hint.active')
				: status === 'archive'
					? t('ws.status.hint.archive')
					: ''
	);
</script>

<span
	class="inline-flex items-center gap-1.5 {cls}"
	title={hint || undefined}
	aria-label={hint ? `${label} — ${hint}` : undefined}
>
	<span class="h-1.5 w-1.5 flex-none rounded-full {c.dot}" aria-hidden="true"></span>
	<span class="text-xs font-medium {c.text}">{label}</span>
</span>
