# 02: Bangun CLI dan vertical slice TUI

**What to build:** Pengguna dapat menjalankan `reamd <path>` untuk membuka satu file di terminal interaktif, membaca dan menggulir tampilan dasar, keluar kembali ke shell, atau menerima teks biasa saat output dialihkan.

**Blocked by:** Tidak ada; Ticket 01 selesai.

**Status:** ready-for-agent

**Acuan:** Fase 2-3 di `docs/development-workflow.md`; `AC-CLI` dan bagian dasar `AC-TUI`/`AC-THEME`.

- [ ] Buat fondasi Go sesuai keputusan stack yang disetujui dan jaga pemisahan CLI, pemuatan file, model dokumen, layout, dan TUI.
- [ ] Dukung satu path file, `--help`/`-h`, `--version`/`-v`, dan `--theme` dengan validasi sesuai `docs/command-spec.md`.
- [ ] Saat tanpa argumen, tampilkan bantuan dan contoh dalam Bahasa Inggris lalu keluar dengan kode 0.
- [ ] Validasi jumlah argumen, opsi tema, file tidak ditemukan, file tidak dapat dibaca, input non-UTF-8/biner, pesan stderr, dan kode keluar 1 atau 2 sesuai jenis error.
- [ ] Di terminal interaktif, tampilkan konten dasar dalam TUI yang dapat digulir dengan keyboard; `q`, `Esc`, dan `?` berfungsi. **Dipindahkan dari checklist Ticket 01:** uji keyboard adalah acceptance aplikasi; API/simulasi spike bukan bukti interaktif.
- [ ] Tangani ukuran minimum 40×10 dan resize; pulihkan terminal setelah keluar normal, error, panic, dan interupsi. **Dipindahkan dari checklist Ticket 01:** resize runtime dan pemulihan terminal nyata milik vertical slice, bukan gerbang pemilihan stack.
- [ ] Saat output dialihkan atau dipipe, jangan ambil alih layar TUI dan jangan keluarkan ANSI secara default.
- [ ] Verifikasi path dengan spasi/Unicode serta pastikan file tidak berubah setelah dibuka.
- [ ] Jalankan smoke test vertical slice pada platform yang tersedia dan catat platform/terminal serta batasannya; matriks seluruh OS/terminal tetap milik Ticket 10. **Dipindahkan dari checklist Ticket 01:** smoke aplikasi bukan bukti spike.

**Verifikasi:** acceptance criteria di atas lulus pada platform yang tersedia; cantumkan command yang benar-benar dijalankan dan keterbatasan platform yang belum diuji.
