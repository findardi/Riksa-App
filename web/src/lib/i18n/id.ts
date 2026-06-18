// Bahasa Indonesia — default locale. Add en.ts and register it in index.ts to expand.
export const id = {
	'brand.name': 'Wadi',
	'brand.tagline': 'Ruang data aman untuk dokumen deal Anda.',
	'brand.reassure.title': 'Dibangun untuk dokumen rahasia',
	'brand.reassure.body':
		'Akses terkontrol dan teraudit. Setiap berkas, izin, dan aktivitas tercatat — Anda selalu tahu siapa mengakses apa.',

	'nav.toLogin': 'Sudah punya akun?',
	'nav.toLoginCta': 'Masuk',
	'nav.toRegister': 'Belum punya akun?',
	'nav.toRegisterCta': 'Buat akun',

	'login.title': 'Masuk ke Wadi',
	'login.subtitle': 'Lanjutkan ke ruang data Anda.',
	'login.identifier': 'Email atau username',
	'login.password': 'Kata sandi',
	'login.submit': 'Masuk',
	'login.submitting': 'Memproses…',
	'login.registered': 'Akun berhasil dibuat. Silakan masuk.',
	'login.reset': 'Kata sandi berhasil diatur ulang. Silakan masuk.',
	'login.verified': 'Email berhasil diverifikasi. Silakan masuk kembali.',
	'login.forgot': 'Lupa kata sandi?',
	'login.or': 'atau',
	'login.ssoError': 'Gagal masuk dengan Google atau GitHub. Coba lagi.',
	'login.ssoConflict': 'Email ini sudah terdaftar dengan kata sandi. Masuk memakai kata sandi.',

	'sso.google': 'Lanjutkan dengan Google',
	'sso.github': 'Lanjutkan dengan GitHub',
	'sso.redirecting': 'Mengalihkan…',

	'register.title': 'Buat akun Wadi',
	'register.subtitle': 'Mulai kelola ruang data Anda dalam hitungan menit.',
	'register.email': 'Email',
	'register.username': 'Username',
	'register.password': 'Kata sandi',
	'register.passwordHint': 'Minimal 6 karakter.',
	'register.usernameHint': 'Minimal 6 karakter.',
	'register.submit': 'Buat akun',
	'register.submitting': 'Membuat akun…',
	'register.emailContinue': 'Lanjutkan',
	'register.checking': 'Memeriksa…',
	'register.changeEmail': 'Ubah',
	'register.emailOk': 'Email tersedia — lengkapi data akun Anda.',

	'forgot.title': 'Lupa kata sandi',
	'forgot.subtitle': 'Masukkan email akun Anda. Kami kirim kode untuk mengatur ulang kata sandi.',
	'forgot.email': 'Email',
	'forgot.send': 'Kirim kode',
	'forgot.sending': 'Mengirim…',
	'forgot.sent': 'Kode OTP telah dikirim ke email Anda.',
	'forgot.otpTitle': 'Masukkan kode 6 digit',
	'forgot.otpSubtitle': 'Kami mengirim kode ke {email}.',
	'forgot.expiresIn': 'Kode kedaluwarsa dalam {time}',
	'forgot.expired': 'Kode kedaluwarsa.',
	'forgot.resend': 'Kirim ulang',
	'forgot.resent': 'Kode dikirim ulang.',
	'forgot.changeEmail': 'Ubah email',
	'forgot.verify': 'Verifikasi kode',
	'forgot.verifying': 'Memverifikasi…',
	'forgot.back': 'Kembali ke halaman masuk',

	'reset.subtitle': 'Buat kata sandi baru untuk {email}.',
	'reset.newPassword': 'Kata sandi baru',
	'reset.confirmPassword': 'Konfirmasi kata sandi',
	'reset.passwordHint': 'Minimal 6 karakter.',
	'reset.submit': 'Simpan kata sandi baru',
	'reset.submitting': 'Menyimpan…',
	'reset.mismatch': 'Konfirmasi kata sandi tidak cocok.',

	'verify.title': 'Verifikasi email Anda',
	'verify.subtitle': 'Masukkan kode 6 digit yang kami kirim ke {email}.',
	'verify.otpTitle': 'Kode verifikasi',
	'verify.sent': 'Kode verifikasi telah dikirim ke email Anda.',
	'verify.submit': 'Verifikasi',
	'verify.verifying': 'Memverifikasi…',
	'verify.noCode': 'Tidak menerima kode?',
	'verify.resend': 'Kirim ulang',
	'verify.resending': 'Mengirim…',
	'verify.resendIn': 'Kirim ulang ({s}d)',
	'verify.resent': 'Kode dikirim ulang.',
	'verify.logout': 'Keluar',

	// App shell (post-login, level-0 / user-scoped)
	'app.nav.rooms': 'Ruang data',
	'app.nav.invitations': 'Undangan',
	'app.nav.settings': 'Pengaturan',
	'app.nav.soon': 'Segera',
	'app.search.placeholder': 'Cari…',
	'app.menu.open': 'Buka menu navigasi',
	'app.account.signedInAs': 'Masuk sebagai',
	'app.account.logout': 'Keluar',

	// Workspaces (ruang data)
	'ws.title': 'Ruang data',
	'ws.create': 'Buat ruang data',
	'ws.count': '{n} ruang data',
	'ws.limitReached': 'Anda sudah mencapai batas 3 ruang data sebagai pemilik.',
	'ws.created': 'Ruang data "{name}" dibuat.',
	'ws.loadError': 'Gagal memuat ruang data. Coba muat ulang halaman.',
	'ws.empty.title': 'Belum ada ruang data',
	'ws.empty.body':
		'Buat ruang data pertama untuk mulai membagikan dokumen dengan akses terkontrol dan teraudit.',
	'ws.dialog.title': 'Buat ruang data',
	'ws.dialog.subtitle': 'Beri nama ruang data Anda.',
	'ws.dialog.reassure':
		'Privat secara default — akses terkontrol dan teraudit. Anda mengundang anggota setelahnya.',
	'ws.field.name': 'Nama',
	'ws.field.namePlaceholder': 'mis. Project Falcon',
	'ws.field.description': 'Deskripsi',
	'ws.field.descriptionHint': 'Opsional — ringkas tujuan ruang data ini.',
	'ws.field.descriptionPlaceholder': 'Opsional',
	'ws.dialog.cancel': 'Batal',
	'ws.dialog.submit': 'Buat',
	'ws.dialog.submitting': 'Membuat…',
	'ws.err.nameTaken': 'Nama ruang data sudah dipakai.',
	'ws.err.nameInvalid': 'Nama harus mengandung huruf atau angka.',
	'ws.err.limit': 'Maksimal 3 ruang data per akun.',
	'ws.err.invalidStatus': 'Status tidak valid.',

	// Workspace detail (/workspace/[slug])
	'ws.detail.back': 'Ruang data',
	'ws.detail.notFound': 'Ruang data tidak ditemukan.',
	'ws.detail.forbidden': 'Anda tidak punya akses ke ruang data ini.',
	'ws.detail.created': 'Dibuat',
	'ws.detail.updated': 'Diperbarui',
	'ws.section.overview': 'Ikhtisar',
	'ws.section.documents': 'Dokumen',
	'ws.section.activity': 'Aktivitas',
	'ws.section.access': 'Manajemen Akses',

	// Workspace status (lifecycle)
	'ws.status.label': 'Status',
	'ws.status.prepare': 'Persiapan',
	'ws.status.active': 'Aktif',
	'ws.status.archive': 'Arsip',
	'ws.status.hint.prepare': 'Ruang masih disiapkan — belum dibagikan ke pihak luar.',
	'ws.status.hint.active': 'Ruang aktif — pihak dengan akses dapat membukanya.',
	'ws.status.hint.archive': 'Ruang diarsipkan — hanya-baca, disimpan untuk audit.',
	'ws.status.updated': 'Status ruang data diperbarui.',

	// Edit room
	'ws.edit.open': 'Edit',
	'ws.edit.title': 'Edit ruang data',
	'ws.edit.submit': 'Simpan',
	'ws.edit.saved': 'Perubahan disimpan.',

	// Delete room
	'ws.delete.open': 'Hapus',
	'ws.delete.title': 'Hapus ruang data',
	'ws.delete.body':
		'Menghapus ruang data bersifat permanen — seluruh dokumen dan jejaknya hilang dan tidak dapat dipulihkan.',
	'ws.delete.warning': 'Tindakan ini permanen. Ruang "{name}" dan seluruh isinya akan dihapus.',
	'ws.delete.confirmLabel': 'Ketik {name} untuk konfirmasi',
	'ws.delete.submit': 'Hapus permanen',
	'ws.delete.submitting': 'Menghapus…',

	// Error page (in-app)
	'error.title': 'Ada yang tidak beres',
	'error.home': 'Kembali ke beranda',

	// Home (post-login landing)
	'home.welcome': 'Selamat datang, {name}',
	'home.welcomeGeneric': 'Selamat datang di Wadi',
	'home.subtitle': 'Buka ruang data untuk melanjutkan, atau buat yang baru.',
	'home.quickActions': 'Aksi cepat',
	'home.action.workspaces': 'Ruang data',
	'home.action.workspacesDesc': 'Kelola dan buat ruang data Anda',

	'password.show': 'Tampilkan kata sandi',
	'password.hide': 'Sembunyikan kata sandi',

	'otp.group': 'Kode OTP',
	'otp.digit': 'Digit ke-{n}',

	// Field-level validation (client + mapped from backend)
	'err.required': 'Wajib diisi',
	'err.email': 'Format email tidak valid',
	'err.min': 'Minimal {n} karakter',
	'err.max': 'Maksimal {n} karakter',
	'err.identifierRequired': 'Masukkan email atau username',

	// Form-level
	'err.invalidCredentials': 'Email/username atau kata sandi salah.',
	'err.emailTaken': 'Email sudah terdaftar.',
	'err.usernameTaken': 'Username sudah dipakai.',
	'err.network': 'Tidak dapat terhubung ke server. Coba lagi.',
	'err.generic': 'Terjadi kesalahan. Coba lagi sebentar.',
	'err.invalidOtp': 'Kode OTP salah atau kedaluwarsa.',

	// management access
	'ma.title': 'Manajemen akses',
	'ma.desc': 'Kelola anggota, peran, dan grup ruang ini.',
	'ma.member': 'Anggota',
	'ma.role': 'Peran',
	'ma.group': 'Grup'
} as const;

export type Dict = typeof id;
