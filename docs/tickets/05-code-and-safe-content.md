# 05: Render code dan amankan konten Markdown

**What to build:** Code block dan Mermaid tetap terbaca, sementara raw HTML dan karakter kontrol dokumen tidak pernah dijalankan sebagai instruksi terminal.

**Blocked by:** 02, Bangun CLI dan vertical slice TUI.

**Status:** blocked

**Acuan:** Fase 4 di `docs/development-workflow.md`; `AC-MD`, `AC-RENDER`, `AC-SEC`; `FX-005`, `FX-007`, `FX-008`, dan `FX-009`.

- [ ] Render fenced code dan inline code sebagai konten kode; wrap baris code panjang dan beri penanda kelanjutan.
- [ ] Tampilkan blok Mermaid dengan penanda diagram dan source code, tanpa merender diagram.
- [ ] Tampilkan raw HTML sebagai teks aman; jangan mengeksekusi tag atau isinya.
- [ ] Netralisasi ESC, BEL, ANSI/OSC, dan karakter kontrol agar tidak dapat membersihkan layar, mengubah judul terminal, atau mengubah clipboard.
- [ ] Tambahkan fixture test yang berisi byte kontrol aktual, bukan hanya representasi teks `\\x1b`.
- [ ] Pastikan tidak ada network request ketika dokumen atau gambar dirender.

**Verifikasi:** jalankan test sanitasi dan periksa hasil di terminal/pseudo-terminal untuk memastikan tidak ada control sequence aktif dari dokumen.
