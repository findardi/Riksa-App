// Shared toast state — a single transient notification, app-wide. Pages call
// `showToast(...)` and render one `<Toaster />`; no per-page duplication.
type Variant = 'success' | 'error';

// Reactive holder; `.current` is the live notification (null when none).
export const store = $state<{ current: { message: string; variant: Variant } | null }>({
	current: null
});

let timer: ReturnType<typeof setTimeout>;

export function showToast(message: string, variant: Variant = 'success') {
	store.current = { message, variant };
	clearTimeout(timer);
	timer = setTimeout(() => (store.current = null), 4000);
}
