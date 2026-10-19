# 10: Verifikasi lintas platform

**What to build:** Bukti bahwa seluruh kriteria MVP yang telah diimplementasikan lulus di Windows, macOS, dan Linux, dengan keterbatasan terminal yang dicatat secara jujur.

**Blocked by:** 04, 05, 06, 07, 08, dan 09 selesai.

**Status:** blocked

**Acuan:** Fase 6 di `docs/development-workflow.md`; seluruh `docs/acceptance-criteria.md`, khususnya `AC-PLATFORM`.

- [ ] Jalankan test otomatis parsing, layout, input file, sanitasi, help/version, status keluar, dan output non-TTY pada Windows, macOS, serta Linux.
- [ ] Jalankan build/test CI pada ketiga OS dan catat hasil yang benar-benar diperoleh.
- [ ] Jalankan test terminal/pseudo-terminal untuk keyboard, resize, exit, dan pemulihan terminal jika didukung stack terpilih.
- [ ] Smoke test Windows Terminal, Terminal.app, dan GNOME Terminal atau emulator Linux yang kompatibel.
- [ ] Periksa tema, zoom, mouse, scrollbar, seleksi/copy, anchor, pembukaan link, error, dan interupsi sesuai kemampuan terminal.
- [ ] Catat OS, terminal, hasil PASS/FAIL, command test, dan keterbatasan yang belum terselesaikan.
- [ ] Pastikan seluruh kriteria wajib lulus atau ada blocker rilis yang disetujui dan terdokumentasi.
- [ ] Verifikasi output teks hasil render melalui terminal/pseudo-terminal: byte ESC/OSC/karakter kontrol dari fixture tidak boleh menjadi kontrol terminal aktif. **Dipindahkan dari checklist Ticket 01:** probe sanitizer scratch tidak melewati renderer aplikasi.
- [ ] Uji keyboard interaktif, resize runtime, pemulihan terminal nyata setelah keluar/error/interupsi, dan smoke test lintas OS/terminal; catat hasil per OS. **Dipindahkan dari checklist Ticket 01:** dokumentasi API dan simulator bukan bukti runtime/OS.
- [ ] Smoke test drag selection aplikasi, clipboard OS plain text tanpa ANSI, dan fallback clipboard pada sistem desktop yang tersedia; catat OS/terminal. **Dipindahkan dari checklist Ticket 01:** test event simulasi dan fake clipboard bukan acceptance produk.
- [ ] Smoke test integrasi browser OS setelah aktivasi pengguna serta verifikasi pesan error TUI; catat hasil sebenarnya. **Dipindahkan dari checklist Ticket 01:** opener tiruan bukan browser OS atau pesan TUI.
- [ ] Setelah package clipboard produksi dipilih, periksa ulang graph dependency, lisensi, versi Go, dan advisory graph final. **Dipindahkan dari checklist Ticket 01:** graph terdahulu tidak mencakup package clipboard yang belum dipilih.
- [ ] Tindak lanjut non-blocking: tinjau maintainer/isu satu per satu untuk dependency terpilih; eskalasi menjadi blocker hanya bila ditemukan temuan spesifik yang berpotensi mengubah pilihan dependency. **Asal:** checklist Ticket 01; belum dilakukan, bukan klaim lulus.

**Verifikasi:** checklist acceptance criteria lengkap dengan bukti otomatis/manual; jangan menandai test yang belum dijalankan sebagai lulus.
