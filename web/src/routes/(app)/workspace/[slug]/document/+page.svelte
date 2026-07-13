<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { t } from '$lib/i18n';
	import type { FolderTreeNode } from '$lib/types/content';

	const folders = $derived((page.data as { folders: FolderTreeNode[] }).folders);
	const defaultFolder = $derived(folders.find((f) => f.is_default) ?? null);
	const slug = $derived(page.params.slug!);
</script>

<div
	class="flex min-h-64 flex-col items-center justify-center gap-3 rounded-box border border-dashed border-base-content/15 px-6 py-16 text-center"
>
	<svg
		class="h-9 w-9 text-muted/70"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="1.4"
		stroke-linecap="round"
		stroke-linejoin="round"
		aria-hidden="true"
	>
		<path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z" />
		<path d="M14 3v5h5" />
	</svg>

	<div>
		<p class="text-[0.9375rem] font-medium">{t('doc.pick.title')}</p>
		<p class="mx-auto mt-1 max-w-sm text-sm text-muted text-pretty">
			{#if folders.length > 0}
				{t('doc.pick.body')}
			{:else}
				{t('doc.pick.bodyEmpty')}
			{/if}
		</p>
	</div>

	{#if defaultFolder}
		<a
			href={resolve('/(app)/workspace/[slug]/document/[folderId]', {
				slug,
				folderId: defaultFolder.id
			})}
			class="btn btn-primary btn-sm mt-1"
		>
			{t('doc.pick.openDefault', { name: defaultFolder.name })}
		</a>
	{/if}
</div>
