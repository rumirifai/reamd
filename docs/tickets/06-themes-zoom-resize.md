# 06: Tambahkan tema, zoom, dan resize

**What to build:** Pengguna dapat membaca dokumen pada ukuran terminal yang berubah, memakai empat tema yang disepakati, serta mengubah kepadatan layout dengan zoom aplikasi.

**Blocked by:** 03, Render elemen inti CommonMark.

**Status:** in-progress

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-TUI`, `AC-THEME`, `AC-RENDER`.

- [x] Dukung tema `auto`, `dark`, `light`, dan `transparent`; default `auto` memakai warna terminal, tema `dark`/`light` memakai palet foreground/background, dan `transparent` mempertahankan keduanya.
- [x] Hormati `NO_COLOR`; output tanpa warna tetap memiliki hierarki melalui bold, underline, marker, dan indentasi.
- [x] Dukung zoom awal 100%, langkah 10%, rentang 70%-150%, alternatif `+`/`-`, dan reset `Ctrl+0`.
- [x] Zoom mengubah lebar wrapping, pemadatan jarak, dan penekanan heading; tidak mengubah klaim ukuran font emulator.
- [x] Tata ulang konten saat ukuran terminal berubah, termasuk tabel responsif.
- [x] Verifikasi otomatis bahwa konten penting tidak hanya dapat dibedakan lewat warna.

**Verifikasi:** test otomatis memeriksa empat tema, `NO_COLOR`, marker struktural, batas/langkah/reset zoom, wrapping, tabel responsif pada lebar 40 dan 20 kolom, serta pemulihan terminal event loop. Resize TUI menerapkan ulang lebar dan tinggi pada setiap `EventResize`.

**Smoke test manual:** belum dijalankan. Sesi berjalan di Windows Terminal (`WT_SESSION` tersedia), tetapi proses shell yang tersedia tidak menyediakan kendali input interaktif ke sesi TUI. Status ticket tetap `in-progress` sampai smoke test manual benar-benar dilakukan.
