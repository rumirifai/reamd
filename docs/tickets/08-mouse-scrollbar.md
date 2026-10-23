# 08: Tambahkan scroll mouse dan scrollbar

**What to build:** Pada terminal yang melaporkan dukungan mouse, pengguna dapat menggulir dokumen dengan wheel dan berpindah cepat melalui scrollbar; keyboard tetap tersedia sebagai fallback.

**Blocked by:** 02, Bangun CLI dan vertical slice TUI.

**Status:** done

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE` dan `AC-PLATFORM`.

- [x] Scroll wheel menggulir dokumen secara konsisten.
- [x] Klik track scrollbar memindahkan posisi dokumen ke lokasi yang sesuai.
- [x] Thumb scrollbar mencerminkan posisi dokumen dan dapat di-drag.
- [x] Fitur mouse aktif hanya ketika terminal melaporkan kapabilitas yang sesuai.
- [x] Navigasi keyboard tetap dapat mengakses seluruh fungsi utama saat mouse tidak tersedia.

**Bukti otomatis:** test mencakup wheel naik/turun, klik track di kedua ujung, posisi dan ukuran thumb, drag thumb, penolakan event mouse saat kapabilitas tidak tersedia, fallback navigasi keyboard, serta aktivasi link dengan mouse pada screen simulasi yang melaporkan dukungan mouse.

**Smoke test manual:** Windows 10, aplikasi Windows Terminal, profile Windows PowerShell (nama/versi Windows Terminal tidak dicatat; profile dan aplikasi terminal terlihat pada screenshot pengguna). Pengguna melaporkan telah menjalankan seluruh langkah pengujian, termasuk wheel, scrollbar, terminal dengan mouse, dan fallback keyboard, tanpa menemukan keanehan atau error.

**Hasil:** checklist Ticket 08 terpenuhi berdasarkan test otomatis dan laporan smoke test manual dari pengguna.
