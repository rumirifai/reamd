# 06: Tambahkan tema, zoom, dan resize

**What to build:** Pengguna dapat membaca dokumen pada ukuran terminal yang berubah, memakai empat tema yang disepakati, serta mengubah kepadatan layout dengan zoom aplikasi.

**Blocked by:** 03, Render elemen inti CommonMark.

**Status:** blocked

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-TUI`, `AC-THEME`, `AC-RENDER`.

- [ ] Dukung tema `auto`, `dark`, `light`, dan `transparent`; default `auto` mengikuti kemampuan dan warna terminal.
- [ ] Hormati `NO_COLOR`; output tanpa warna tetap memiliki hierarki yang terbaca.
- [ ] Dukung zoom awal 100%, langkah 10%, rentang 70%-150%, alternatif `+`/`-`, dan reset `Ctrl+0`.
- [ ] Zoom mengubah layout aplikasi, penekanan heading, jarak, dan wrapping; jangan mengklaim mengubah ukuran font emulator.
- [ ] Tata ulang konten segera saat ukuran terminal berubah, termasuk tabel responsif.
- [ ] Verifikasi bahwa konten penting tidak hanya dapat dibedakan lewat warna.

**Verifikasi:** test layout/tema otomatis lulus dan smoke test manual mencatat hasil pada tema serta resize yang benar-benar dicoba.
