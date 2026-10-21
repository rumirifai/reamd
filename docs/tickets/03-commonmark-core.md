# 03: Render elemen inti CommonMark

**What to build:** Dokumen Markdown dasar tampil terstruktur dan mudah dibaca, dengan urutan, nesting, penekanan, serta wrapping yang benar.

**Blocked by:** 02, Bangun CLI dan vertical slice TUI.

**Status:** done

**Acuan:** Fase 4 di `docs/development-workflow.md`; `AC-MD`, `AC-RENDER`; `FX-001`, `FX-002`, dan bagian relevan `FX-006`.

- [x] Render heading, paragraf, emphasis, inline code, line break, dan horizontal rule dengan hierarki yang terlihat.
- [x] Render ordered/unordered list bertingkat dan blockquote dengan indentasi lanjutan yang konsisten.
- [x] Pastikan konten parser dirender atau diberi fallback aman; test mencakup HTML literal, strikethrough, link/image URL, dan source fenced block.
- [x] Bungkus paragraf, list, dan quote sesuai lebar tampilan Unicode, bukan jumlah byte.
- [x] Tambahkan test fixture FX-001, FX-002, FX-006 untuk urutan, nesting, wrapping pada 40/80/120 kolom, serta aksen/CJK/emoji.
- [x] Verifikasi versi Goldmark `v1.8.6` di `go.mod`. Conformance corpus CommonMark 0.31.2 Goldmark lulus 652/652 pada Ticket 01, tercatat di `docs/stack-evaluation.md`; corpus tidak diulang sesuai lingkup Ticket 03.

## Bukti penyelesaian

- Parser AST, model dokumen, layout Unicode, dan TUI tetap terpisah; output TUI dan pipe menggunakan baris hasil model yang sama.
- Pipe mencetak plain text tanpa ANSI; test integrasi memeriksa output pipe dan input yang diterima runner TUI.
- Pemeriksaan kontrol tetap melalui `layout.Sanitize` pada tahap layout dan output pipe.
- Fallback link menampilkan alamat, image menampilkan alt text dan alamat, raw HTML tampil literal, serta fenced block menampilkan isi sebagai teks biasa. Fallback bukan renderer khusus untuk fitur tersebut.
- Pemeriksaan smoke test TUI manual bukan bagian checklist Ticket 03; test parser/layout dan alur output otomatis lulus.

**Verifikasi:** fixture terkait dan conformance test yang dipilih lulus; layout dicoba pada lebar terminal 40, 80, dan 120 kolom.
