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

**Verifikasi:** checklist acceptance criteria lengkap dengan bukti otomatis/manual; jangan menandai test yang belum dijalankan sebagai lulus.
