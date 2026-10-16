# 11: Siapkan distribusi dan instalasi

**What to build:** Pengguna dapat memasang `reamd` melalui metode yang dipilih dan menjalankan binary tanpa memasang Go.

**Blocked by:** 10, Verifikasi lintas platform.

**Status:** blocked

**Acuan:** Fase 7 di `docs/development-workflow.md`; sasaran distribusi pada `README.md` dan `docs/acceptance-criteria.md`.

- [ ] Minta dan catat persetujuan metode distribusi yang dipilih, seperti binary GitHub Release atau package manager.
- [ ] Build artifact untuk setiap OS/arsitektur yang diumumkan dan uji binary hasil build.
- [ ] Sertakan versi dan checksum yang dapat diverifikasi.
- [ ] Uji instruksi install/uninstall pada lingkungan bersih yang tidak memiliki toolchain Go.
- [ ] Jangan mengandalkan `go install` sebagai metode utama bagi pengguna yang tidak memasang Go.
- [ ] Perbarui README dalam Bahasa Inggris dan Indonesia dengan metode serta command yang benar-benar diuji.

**Verifikasi:** fresh install berhasil menjalankan `reamd`; semua klaim platform memiliki artifact yang diuji.
