# 04: Lengkapi GFM dan tabel responsif

**What to build:** Sintaks GFM yang disepakati terbaca, termasuk tabel, task list, strikethrough, dan autolink; tabel tetap menyampaikan semua nilai pada terminal sempit.

**Blocked by:** 03, Render elemen inti CommonMark.

**Status:** blocked

**Acuan:** Fase 4 di `docs/development-workflow.md`; `AC-MD`, `AC-RENDER`; `FX-003` dan `FX-010`.

- [ ] Render tabel sebagai tabel saat kolom muat dan sebagai pasangan label/nilai ketika ruang terlalu sempit.
- [ ] Pastikan seluruh isi sel tetap terbaca pada layout bertumpuk; jangan memotong nilai diam-diam.
- [ ] Render status task checked dan unchecked dengan pembeda yang tidak hanya mengandalkan warna.
- [ ] Dukungan parser mencakup strikethrough dan autolink sesuai profil GFM yang dipilih.
- [ ] Tambahkan test fixture untuk tabel, alignment, task list, strikethrough, autolink, dan beberapa lebar terminal.
- [ ] Jalankan conformance test GFM sesuai versi parser dan catat profil yang dipakai.

**Verifikasi:** acceptance criteria dan fixture GFM terkait lulus pada lebar representatif 40, 80, dan 120 kolom.
