# 05: Render code dan amankan konten Markdown

**What to build:** Code block dan Mermaid tetap terbaca, sementara raw HTML dan karakter kontrol dokumen tidak pernah dijalankan sebagai instruksi terminal.

**Blocked by:** 02, Bangun CLI dan vertical slice TUI.

**Status:** done

**Acuan:** Fase 4 di `docs/development-workflow.md`; `AC-MD`, `AC-RENDER`, `AC-SEC`; `FX-005`, `FX-007`, `FX-008`, dan `FX-009`.

- [x] Render fenced code dan inline code sebagai konten kode; wrap baris code panjang dan beri penanda kelanjutan.
- [x] Tampilkan blok Mermaid dengan penanda diagram dan source code, tanpa merender diagram.
- [x] Tampilkan raw HTML sebagai teks aman; jangan mengeksekusi tag atau isinya. **Dipindahkan dari checklist Ticket 01:** probe fallback parser bukan uji renderer aplikasi.
- [x] Netralisasi ESC, BEL, ANSI/OSC, dan karakter kontrol sebelum output terminal agar tidak dapat membersihkan layar, mengubah judul terminal, atau mengubah clipboard. **Dipindahkan dari checklist Ticket 01:** sanitasi melalui output aplikasi adalah acceptance produk, bukan spike.
- [x] Tambahkan fixture test yang berisi byte kontrol aktual, bukan hanya representasi teks `\\x1b`; verifikasi jalur output pipe dan layar TUI simulasi tidak mengeluarkan kontrol aktif. **Dipindahkan dari checklist Ticket 01:** probe safePlainText scratch belum melewati renderer/output terminal.
- [x] Pastikan tidak ada network request ketika dokumen atau gambar dirender; pemeriksaan memakai server HTTP lokal dan tidak mengirim request ke situs luar.

**Bukti:** test FX-005/007/008 memuat fixture bernama dan melewati document/layout; FX-009 memberikan byte ESC, BEL, C1, dan DEL aktual melalui renderer pipe serta layar simulasi tcell. Test memastikan notasi tampak aman dan tidak ada byte kontrol dalam output. Test jaringan lokal memastikan render tautan dan gambar tidak membuat request.

**Verifikasi:** test sanitasi mencakup kedua jalur keluaran aplikasi; tidak dilakukan smoke test terminal GUI nyata pada ticket ini.
