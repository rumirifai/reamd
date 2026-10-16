# 08: Tambahkan scroll mouse dan scrollbar

**What to build:** Pada terminal yang melaporkan dukungan mouse, pengguna dapat menggulir dokumen dengan wheel dan berpindah cepat melalui scrollbar; keyboard tetap tersedia sebagai fallback.

**Blocked by:** 02, Bangun CLI dan vertical slice TUI.

**Status:** blocked

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE` dan `AC-PLATFORM`.

- [ ] Scroll wheel menggulir dokumen secara konsisten.
- [ ] Klik track scrollbar memindahkan posisi dokumen ke lokasi yang sesuai.
- [ ] Thumb scrollbar mencerminkan posisi dokumen dan dapat di-drag.
- [ ] Fitur mouse aktif hanya ketika terminal melaporkan kapabilitas yang sesuai.
- [ ] Navigasi keyboard tetap dapat mengakses seluruh fungsi utama saat mouse tidak tersedia.

**Verifikasi:** smoke test manual pada terminal dengan dukungan mouse dan pemeriksaan fallback keyboard pada terminal tanpa dukungan tersebut.
