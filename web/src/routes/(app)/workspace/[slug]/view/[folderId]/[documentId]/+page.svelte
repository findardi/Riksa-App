<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { ViewerPage } from '$lib/components/app';
	import { showToast } from '$lib/components/common';
	import { t } from '$lib/i18n';
	import type { WorkspaceData, MyAccessWorkspace } from '$lib/types/workspace';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const meta = $derived(data.meta);
	const forbidden = $derived(data.forbidden);
	const notViewable = $derived(data.notViewable);

	const workspace = $derived((page.data as { workspace: WorkspaceData }).workspace);
	const access = $derived((page.data as { access?: MyAccessWorkspace }).access);
	const perms = $derived(access?.permissions ?? []);

	const slug = $derived(page.params.slug!);
	const folderId = $derived(page.params.folderId!);
	const documentId = $derived(page.params.documentId!);

	// The folder lives in the path, so back returns to exactly the list the
	// document was opened from.
	const backHref = $derived(
		resolve('/(app)/workspace/[slug]/document/[folderId]', { slug, folderId })
	);

	const pageCount = $derived(meta?.page_count ?? 0);
	const pages = $derived(Array.from({ length: pageCount }, (_, i) => i + 1));

	const pageSrc = (n: number) =>
		`/api/content/pages?workspaceId=${encodeURIComponent(workspace.id)}` +
		`&documentId=${encodeURIComponent(documentId)}&page=${n}`;

	// --- current page tracking (max on-screen coverage wins) ---
	// Plain Maps, deliberately non-reactive: the UI reads only `currentPage`
	// ($state); these are imperative scratch state, read in handlers, never markup.
	// eslint-disable-next-line svelte/prefer-svelte-reactivity
	const coverage = new Map<number, number>();
	let currentPage = $state(1);

	function onactive(p: number, cov: number) {
		if (cov <= 0) coverage.delete(p);
		else coverage.set(p, cov);
		if (coverage.size === 0) return;
		let best = currentPage;
		let bestCov = -1;
		for (const [pg, c] of coverage) {
			if (c > bestCov) {
				bestCov = c;
				best = pg;
			}
		}
		if (best !== currentPage) currentPage = best;
	}

	// --- element registry for jump + step navigation ---
	// eslint-disable-next-line svelte/prefer-svelte-reactivity
	const pageEls = new Map<number, HTMLElement>();
	function onregister(p: number, el: HTMLElement | null) {
		if (el) pageEls.set(p, el);
		else pageEls.delete(p);
	}

	function scrollToPage(n: number) {
		const target = Math.min(Math.max(1, n), pageCount);
		const el = pageEls.get(target);
		if (!el) return;
		const reduce = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
		el.scrollIntoView({ behavior: reduce ? 'auto' : 'smooth', block: 'start' });
	}

	// --- jump-to-page input (display follows scroll unless the user is typing) ---
	let jumpEl = $state<HTMLInputElement>();
	let editing = $state(false);

	$effect(() => {
		if (!editing && jumpEl) jumpEl.value = String(currentPage);
	});

	function commitJump() {
		const n = Number.parseInt(jumpEl?.value ?? '', 10);
		editing = false;
		if (Number.isFinite(n)) scrollToPage(n);
		if (jumpEl) jumpEl.value = String(currentPage);
	}

	function onWindowKey(e: KeyboardEvent) {
		if (e.key === 'Escape' && !editing) goto(backHref);
	}

	// --- download (view-and-download access) ---
	let downloading = $state(false);
	async function download() {
		downloading = true;
		try {
			const q = new URLSearchParams({ workspaceId: workspace.id, documentId });
			const res = await fetch(`/api/content/download?${q}`);
			if (res.status === 403) {
				showToast(t('doc.docs.err.forbiddenDownload'), 'error');
				return;
			}
			if (!res.ok) {
				showToast(t('err.generic'), 'error');
				return;
			}
			const { download_url } = (await res.json()) as { download_url: string };
			window.location.href = download_url;
		} catch {
			showToast(t('err.network'), 'error');
		} finally {
			downloading = false;
		}
	}

	const canDownloadOnly = $derived(perms.includes('document:download'));
</script>

<svelte:head>
	<title>{meta?.name ?? t('doc.view.tab')} · {t('brand.name')}</title>
</svelte:head>

<svelte:window onkeydown={onWindowKey} />

<div class="flex h-full min-h-0 flex-col bg-base-200">
	<!-- Reader chrome -->
	<header
		class="flex flex-none items-center gap-2 border-b border-base-content/10 bg-base-100 px-3 py-2 sm:gap-3 sm:px-4"
	>
		<a
			href={backHref}
			class="flex flex-none items-center gap-1.5 rounded-field px-2 py-1.5 text-sm text-muted no-underline transition-colors hover:bg-base-content/5 hover:text-base-content"
		>
			<svg
				class="h-4 w-4"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="1.8"
				stroke-linecap="round"
				stroke-linejoin="round"
				aria-hidden="true"
			>
				<path d="M19 12H5M12 19l-7-7 7-7" />
			</svg>
			<span class="hidden sm:inline">{t('doc.view.back')}</span>
			<span class="sr-only sm:hidden">{t('doc.view.back')}</span>
		</a>

		<span class="h-5 w-px flex-none bg-base-content/10" aria-hidden="true"></span>

		<h1 class="min-w-0 flex-1 truncate text-sm font-medium" title={meta?.name}>
			{meta?.name ?? t('doc.view.tab')}
		</h1>

		{#if meta && pageCount > 0}
			<!-- Page stepper -->
			<div class="flex flex-none items-center gap-0.5">
				<button
					type="button"
					onclick={() => scrollToPage(currentPage - 1)}
					disabled={currentPage <= 1}
					aria-label={t('doc.view.prev')}
					title={t('doc.view.prev')}
					class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content disabled:pointer-events-none disabled:opacity-40"
				>
					<svg
						class="h-4 w-4"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<path d="m18 15-6-6-6 6" />
					</svg>
				</button>

				<div class="flex items-center gap-1 font-mono text-xs text-muted tabular-nums">
					<input
						id="viewer-jump"
						bind:this={jumpEl}
						type="text"
						inputmode="numeric"
						onfocus={() => {
							editing = true;
							jumpEl?.select();
						}}
						onblur={commitJump}
						onkeydown={(e) => {
							if (e.key === 'Enter') {
								e.preventDefault();
								commitJump();
								jumpEl?.blur();
							} else if (e.key === 'Escape') {
								e.stopPropagation();
								editing = false;
								jumpEl?.blur();
							}
						}}
						aria-label={t('doc.view.jumpLabel')}
						class="w-9 rounded-field border border-base-content/15 bg-base-100 px-1 py-0.5 text-center text-xs tabular-nums focus:border-primary focus:outline-none"
					/>
					<span aria-hidden="true">/</span>
					<span aria-label={t('doc.view.pageOf', { n: currentPage, total: pageCount })}>
						{pageCount}
					</span>
				</div>

				<button
					type="button"
					onclick={() => scrollToPage(currentPage + 1)}
					disabled={currentPage >= pageCount}
					aria-label={t('doc.view.next')}
					title={t('doc.view.next')}
					class="grid h-8 w-8 place-items-center rounded-field text-muted transition-colors hover:bg-base-content/5 hover:text-base-content disabled:pointer-events-none disabled:opacity-40"
				>
					<svg
						class="h-4 w-4"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
						aria-hidden="true"
					>
						<path d="m6 9 6 6 6-6" />
					</svg>
				</button>
			</div>

			<span class="hidden h-5 w-px flex-none bg-base-content/10 sm:block" aria-hidden="true"></span>

			<!-- Protection signal — trust is shown, not claimed -->
			<span
				class="hidden flex-none items-center gap-1.5 text-xs text-muted sm:flex"
				title={t('doc.view.protected')}
			>
				<svg
					class="h-4 w-4 text-primary"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="1.7"
					stroke-linecap="round"
					stroke-linejoin="round"
					aria-hidden="true"
				>
					<path d="M12 3 4 6v5c0 4.5 3 8 8 10 5-2 8-5.5 8-10V6z" />
					<path d="m9 12 2 2 4-4" />
				</svg>
				<span class="hidden lg:inline">{t('doc.view.watermarked')}</span>
			</span>

			{#if meta.can_download_original}
				<button
					type="button"
					onclick={download}
					disabled={downloading}
					class="btn btn-ghost btn-sm flex-none gap-1.5"
				>
					{#if downloading}
						<span class="loading loading-spinner loading-xs"></span>
					{:else}
						<svg
							class="h-4 w-4"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="1.8"
							stroke-linecap="round"
							stroke-linejoin="round"
							aria-hidden="true"
						>
							<path d="M12 4v11M7.5 10.5 12 15l4.5-4.5" />
							<path d="M5 19h14" />
						</svg>
					{/if}
					<span class="hidden sm:inline">{t('doc.docs.download')}</span>
				</button>
			{/if}
		{/if}
	</header>

	{#if forbidden}
		<div class="flex flex-1 items-center justify-center overflow-y-auto px-6 py-16">
			<div class="flex max-w-sm flex-col items-center gap-3 text-center">
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
					<rect x="4" y="10.5" width="16" height="10" rx="2" />
					<path d="M8 10.5V7a4 4 0 0 1 8 0v3.5" />
				</svg>
				<div>
					<p class="text-[0.9375rem] font-medium">{t('doc.view.forbidden.title')}</p>
					<p class="mt-1 text-sm text-muted text-pretty">{t('doc.view.forbidden.body')}</p>
				</div>
			</div>
		</div>
	{:else if notViewable}
		<div class="flex flex-1 items-center justify-center overflow-y-auto px-6 py-16">
			<div class="flex max-w-sm flex-col items-center gap-4 text-center">
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
					<path d="M12 12v5M9.5 14.5 12 17l2.5-2.5" />
				</svg>
				<div>
					<p class="text-[0.9375rem] font-medium">{t('doc.view.downloadOnly.title')}</p>
					<p class="mt-1 text-sm text-muted text-pretty">{t('doc.view.downloadOnly.body')}</p>
				</div>
				{#if canDownloadOnly}
					<button
						type="button"
						onclick={download}
						disabled={downloading}
						class="btn btn-primary btn-sm gap-1.5"
					>
						{#if downloading}
							<span class="loading loading-spinner loading-xs"></span>
						{:else}
							<svg
								class="h-4 w-4"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="1.8"
								stroke-linecap="round"
								stroke-linejoin="round"
								aria-hidden="true"
							>
								<path d="M12 4v11M7.5 10.5 12 15l4.5-4.5" />
								<path d="M5 19h14" />
							</svg>
						{/if}
						{t('doc.docs.download')}
					</button>
				{:else}
					<p class="text-xs text-muted text-pretty">{t('doc.view.downloadOnly.noPerm')}</p>
				{/if}
			</div>
		</div>
	{:else if meta && pageCount > 0}
		<div class="min-h-0 flex-1 overflow-y-auto" aria-label={meta.name}>
			<div class="mx-auto flex max-w-[820px] flex-col gap-4 px-3 py-6 sm:px-4">
				{#each pages as n (n)}
					<ViewerPage pageNumber={n} total={pageCount} src={pageSrc(n)} {onactive} {onregister} />
				{/each}
			</div>
		</div>
	{:else}
		<div class="flex flex-1 items-center justify-center overflow-y-auto px-6 py-16">
			<p class="max-w-sm text-center text-sm text-muted text-pretty">{t('doc.view.emptyPages')}</p>
		</div>
	{/if}
</div>
