# 06: Tambahkan tema, zoom, dan resize

**What to build:** Pengguna dapat membaca dokumen pada ukuran terminal yang berubah, memakai empat tema yang disepakati, serta mengubah kepadatan layout dengan zoom aplikasi.

**Blocked by:** 03, Render elemen inti CommonMark.

**Status:** done

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-TUI`, `AC-THEME`, `AC-RENDER`.

- [x] Dukung tema `auto`, `dark`, `light`, dan `transparent`; default `auto` memakai warna terminal, tema `dark`/`light` memakai palet foreground/background, dan `transparent` mempertahankan keduanya.
- [x] Hormati `NO_COLOR`; output tanpa warna tetap memiliki hierarki melalui bold, underline, marker, dan indentasi.
- [x] Dukung zoom awal 100%, langkah 10%, rentang 70%-150%, `=`/`+` untuk memperbesar, `-` untuk memperkecil, dan `0` untuk reset tanpa modifier.
- [x] Zoom mengubah lebar wrapping, pemadatan jarak, dan penekanan heading; tidak mengubah klaim ukuran font emulator.
- [x] Tata ulang konten saat ukuran terminal berubah, termasuk tabel responsif.
- [x] Verifikasi otomatis bahwa konten penting tidak hanya dapat dibedakan lewat warna.

**Verifikasi:** test otomatis memeriksa empat tema, `NO_COLOR`, marker struktural, batas/langkah/reset zoom, wrapping, tabel responsif pada lebar 40 dan 20 kolom, serta pemulihan terminal event loop. Resize TUI menerapkan ulang lebar dan tinggi pada setiap `EventResize`.

**Smoke test manual (laporan pengguna):** Windows, Windows Terminal (versi tidak dicatat). Tema `auto`, `dark`, `light`, `transparent`, dan `NO_COLOR` dicoba. Zoom `=`/`+`, `-`, dan `0` dilaporkan berfungsi; binding lama `Ctrl+0` gagal dan diganti dengan `0` tanpa modifier. Setelah resize sempit, tabel berubah menjadi susunan label/nilai dan seluruh konten tetap terlihat. Pengguna mengonfirmasi keluar dengan `q` memulihkan terminal; prompt PowerShell kembali normal dan dapat digunakan.
