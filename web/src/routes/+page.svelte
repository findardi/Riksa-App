<script lang="ts">
	import { Brand, Button } from '$lib/components/common';
</script>

<svelte:head><title>Wadi</title></svelte:head>

<main class="flex min-h-dvh items-center justify-center p-5">
	<div
		class="flex w-full max-w-md flex-col items-start gap-4 rounded-box border border-base-content/10 bg-base-100 p-8"
	>
		<Brand size={32} />
		<h1 class="text-[1.375rem] font-semibold tracking-[-0.015em]">Anda sudah masuk.</h1>
		<p class="max-w-[36ch] text-[0.9375rem] leading-relaxed text-muted">
			Ruang data — daftar dokumen, izin akses, dan jejak audit — adalah build berikutnya. Untuk saat
			ini, alur autentikasi sudah aktif.
		</p>
		<form method="POST" action="?/logout" class="mt-2">
			<Button type="submit" variant="ghost">Keluar</Button>
		</form>
	</div>
</main>
