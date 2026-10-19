# 03: Render elemen inti CommonMark

**What to build:** Dokumen Markdown dasar tampil terstruktur dan mudah dibaca, dengan urutan, nesting, penekanan, serta wrapping yang benar.

**Blocked by:** 02, Bangun CLI dan vertical slice TUI.

**Status:** ready-for-agent

**Acuan:** Fase 4 di `docs/development-workflow.md`; `AC-MD`, `AC-RENDER`; `FX-001`, `FX-002`, dan bagian relevan `FX-006`.

- [ ] Render heading, paragraf, emphasis, inline code, line break, dan horizontal rule dengan hierarki yang terlihat.
- [ ] Render ordered/unordered list bertingkat dan blockquote dengan indentasi lanjutan yang konsisten.
- [ ] Pastikan seluruh jenis konten yang dikenali parser dirender atau diberi fallback aman; tidak ada isi yang hilang diam-diam. **Dipindahkan dari checklist Ticket 01:** ini acceptance renderer aplikasi, bukan bukti pemilihan stack; probe AST spike terbatas tidak menguji renderer.
- [ ] Bungkus paragraf, list, dan quote sesuai lebar tampilan Unicode, bukan jumlah byte.
- [ ] Tambahkan test berbasis fixture untuk urutan konten, nesting, wrapping, serta karakter aksen/CJK/emoji.
- [ ] Jalankan conformance corpus CommonMark yang sesuai parser terpilih dan catat versi/profilnya. Ticket 01 mencatat bukti parser spike Goldmark 652/652 CommonMark 0.31.2 di `docs/stack-evaluation.md`; bukti itu bukan kelulusan acceptance renderer/layout ticket ini.

**Verifikasi:** fixture terkait dan conformance test yang dipilih lulus; layout dicoba pada lebar terminal 40, 80, dan 120 kolom.
