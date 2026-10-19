# 07: Navigasi anchor dan aktivasi link

**What to build:** Pengguna dapat menjelajahi link dengan keyboard atau klik, melompat ke heading internal, dan membuka link web hanya atas tindakan eksplisit.

**Blocked by:** 03, Render elemen inti CommonMark (Ticket 01 selesai).

**Status:** blocked

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE`; `FX-004`.

- [ ] `Tab` dan `Shift+Tab` memindahkan fokus ke link; `Enter` mengaktifkan link yang sedang fokus.
- [ ] Klik link mengaktifkan link tersebut pada terminal yang mendukung mouse.
- [ ] Anchor `#...` menuju heading dengan aturan slug gaya GitHub, termasuk heading duplikat.
- [ ] Link HTTP/HTTPS membuka browser default hanya setelah `Enter` atau klik eksplisit pengguna.
- [ ] Protokol selain HTTP/HTTPS tidak dijalankan atau dibuka.
- [ ] Kegagalan membuka browser menampilkan pesan yang jelas dan membiarkan TUI tetap dapat digunakan. **Dipindahkan dari checklist Ticket 01:** pesan error adalah acceptance TUI, belum diuji pada spike.
- [ ] Tidak ada link atau request jaringan yang dibuka otomatis saat file dimuat atau dirender.
- [ ] Uji jalur parse/render agar tidak memicu opener, serta uji browser OS hanya sesudah aktivasi eksplisit. **Dipindahkan dari checklist Ticket 01:** validator URL dan opener tiruan lulus di spike, tetapi tidak menguji integrasi aplikasi atau browser OS.

**Verifikasi:** uji fixture anchor dan gate aktivasi otomatis; smoke test browser manual mencatat OS, terminal, aktivasi, serta hasil sebenarnya.
