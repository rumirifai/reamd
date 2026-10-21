# 04: Lengkapi GFM dan tabel responsif

**What to build:** Sintaks GFM yang disepakati terbaca, termasuk tabel, task list, strikethrough, dan autolink; tabel tetap menyampaikan semua nilai pada terminal sempit.

**Blocked by:** 03, Render elemen inti CommonMark.

**Status:** done

**Acuan:** Fase 4 di `docs/development-workflow.md`; `AC-MD`, `AC-RENDER`; `FX-003` dan `FX-010`.

- [x] Render tabel sebagai tabel saat kolom muat dan sebagai pasangan label/nilai ketika ruang terlalu sempit.
- [x] Pastikan seluruh isi sel tetap terbaca pada layout bertumpuk; jangan memotong nilai diam-diam.
- [x] Render status task checked dan unchecked dengan pembeda yang tidak hanya mengandalkan warna.
- [x] Dukungan parser mencakup strikethrough dan autolink sesuai profil GFM yang dipilih.
- [x] Tambahkan test fixture untuk tabel, alignment, task list, strikethrough, autolink, dan beberapa lebar terminal.
- [x] Gunakan bukti conformance GFM Goldmark v1.8.6 yang sudah tercatat di `docs/stack-evaluation.md`; profil GFM spec 0.29 + Linkify dan normalisasi semantik yang disetujui, 0 semantic mismatch dari 648 contoh. Corpus tidak diulang sesuai lingkup Ticket 04.

## Bukti rendering

- Test aplikasi di `internal/document/document_test.go` memeriksa tabel/alignment, task checked dan unchecked, strikethrough, autolink URL/email, serta tabel pada 40, 80, dan 120 kolom. Lebar 40 memakai pasangan label/nilai; lebar 80 dan 120 memakai grid. Pemeriksaan memastikan teks sel penuh tetap ada sesudah wrapping Unicode.
- `layout.BuildTables` menjalankan layout tabel dari data sel yang sama untuk TUI dan output pipe; output pipe tetap teks biasa tanpa ANSI. `internal/app` memiliki test pipe/TUI untuk memastikan parser menghasilkan konten bersama.
- Bukti parser terdahulu: `docs/stack-evaluation.md`, bagian pembaruan 19 Oktober 2026 dan pengulangan test prototipe GFM. Bukti tersebut berlaku pada Goldmark `v1.8.6`, profil GFM yang disepakati, dan normalisasi semantik tercatat; bukti itu adalah conformance parser, bukan pengganti test rendering aplikasi di atas.

**Verifikasi:** acceptance criteria dan fixture GFM terkait lulus pada lebar representatif 40, 80, dan 120 kolom.
