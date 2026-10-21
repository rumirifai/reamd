# 07: Navigasi anchor dan aktivasi link

**What to build:** Pengguna dapat menjelajahi link dengan keyboard atau klik, melompat ke heading internal, dan membuka link web hanya atas tindakan eksplisit.

**Blocked by:** 03, Render elemen inti CommonMark (Ticket 01 selesai).

**Status:** in-progress

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE`; `FX-004`.

- [x] `Tab` dan `Shift+Tab` memindahkan fokus ke link; `Enter` mengaktifkan link yang sedang fokus.
- [x] Klik link mengaktifkan link tersebut pada terminal simulasi yang mendukung mouse.
- [x] Anchor `#...` menuju heading dengan aturan slug gaya GitHub, termasuk tanda baca, spasi, dan heading duplikat; target yang hilang menampilkan pesan.
- [x] Link HTTP/HTTPS membuka browser hanya sebagai hasil `Enter` atau klik eksplisit. Opener OS memakai proses langsung tanpa shell dan dapat diganti dengan opener tiruan untuk test.
- [x] Protokol selain HTTP/HTTPS tidak dijalankan atau dibuka.
- [x] Kegagalan opener menampilkan pesan yang jelas dan membiarkan navigasi/keluar TUI tetap tersedia. **Dipindahkan dari checklist Ticket 01:** pesan error adalah acceptance TUI, belum diuji pada spike.
- [x] Parsing, render, pemuatan, dan perpindahan fokus tidak memanggil opener; render lokal tidak membuat request jaringan.
- [x] Test otomatis fixture FX-004, opener tiruan, fokus keyboard, aktivasi `Enter`, klik simulasi, protokol, dan kegagalan opener lulus.
- [ ] Smoke test browser OS nyata setelah aktivasi eksplisit, mencatat OS, terminal, cara aktivasi, dan hasil sebenarnya. **Belum dijalankan.**

**Verifikasi:** uji fixture anchor dan gate aktivasi otomatis; smoke test browser manual mencatat OS, terminal, aktivasi, serta hasil sebenarnya.

**Bukti implementasi:** model dokumen mempertahankan target link sebagai metadata privat yang dibuang saat layout teks biasa; layout memetakan rentang sel link, TUI menangani fokus/klik, dan opener hanya dipanggil dari aktivasi. Test FX-004 menguji anchor dan protokol tanpa membuka browser. `go test ./...` lulus pada pengujian implementasi; verifikasi final dicatat saat penyelesaian. Smoke test browser OS belum dijalankan, sehingga status ticket tetap `in-progress`.
