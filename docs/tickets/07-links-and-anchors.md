# 07: Navigasi anchor dan aktivasi link

**What to build:** Pengguna dapat menjelajahi link dengan keyboard atau klik, melompat ke heading internal, dan membuka link web hanya atas tindakan eksplisit.

**Blocked by:** 01, Evaluasi stack dan risiko teknis; 03, Render elemen inti CommonMark.

**Status:** blocked

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE`; `FX-004`.

- [ ] `Tab` dan `Shift+Tab` memindahkan fokus ke link; `Enter` mengaktifkan link yang sedang fokus.
- [ ] Klik link mengaktifkan link tersebut pada terminal yang mendukung mouse.
- [ ] Anchor `#...` menuju heading dengan aturan slug gaya GitHub, termasuk heading duplikat.
- [ ] Link HTTP/HTTPS membuka browser default hanya setelah `Enter` atau klik eksplisit pengguna.
- [ ] Protokol selain HTTP/HTTPS tidak dijalankan atau dibuka.
- [ ] Kegagalan membuka browser menampilkan pesan yang jelas dan membiarkan TUI tetap dapat digunakan.
- [ ] Tidak ada link atau request jaringan yang dibuka otomatis saat file dimuat atau dirender.

**Verifikasi:** uji fixture anchor internal dan smoke test aktivasi browser manual; catat terminal dan hasil sebenarnya.
