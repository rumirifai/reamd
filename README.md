# reamd

**A read-only Markdown reader for your terminal.** `reamd` opens one document in the current terminal, lets you scroll through it, and returns you to the shell when you quit.

> [!NOTE]
> `reamd` is in early development. The command line, file loading, scrolling TUI, and plain pipe output work. Markdown rendering, themes, links, mouse input, and clipboard support are not implemented yet.

**Languages:** English | [Bahasa Indonesia](#reamd-bahasa-indonesia)

## Usage

```sh
# Open one Markdown file in the current terminal
reamd README.md

# Choose a theme (validated now, palette work is still in progress)
reamd --theme transparent README.md

# Show usage or version information
reamd --help
reamd --version
```

With no arguments, `reamd` shows help and usage examples and exits with code 0. It accepts one file path per run and does not read from stdin. In an interactive terminal it opens a read-only TUI; when output is redirected or piped, it prints plain text without ANSI styling.

Current keys: `Up`/`Down` and `j`/`k` scroll one line, `PageUp`/`PageDown` scroll one screen, `Home`/`End` jump to the start or end, `?` opens help, and `q` or `Esc` quits. Terminals smaller than 40 columns by 10 rows show a short resize message.

## Current status

The reader parses and displays core CommonMark headings, paragraphs, emphasis, inline code, breaks, lists, and blockquotes in a scrollable TUI or as plain rendered text when output is redirected. Other parsed constructs remain safe text fallbacks until their renderer is implemented.

```sh
# Run from source
go run ./cmd/reamd README.md

# Build a local binary
go build -o reamd ./cmd/reamd
```

## Planned MVP

- Read a Markdown file and show it in a scrollable, read-only TUI.
- Support CommonMark and GitHub Flavored Markdown (GFM), including headings, lists, quotes, tables, links, and code blocks. Bare URLs and email addresses are recognized as autolinks.
- Adapt layout to terminal width. Switch wide tables to a stacked layout when needed.
- Provide in-app zoom by adjusting layout and wrapping. This will not change the terminal emulator's actual font size.
- Support `auto`, `dark`, `light`, and `transparent` themes.
- Support keyboard and mouse navigation, text selection, and copying rendered text to the system clipboard.
- Show image alt text and URLs without downloading images. Open HTTP(S) links in the default browser only after explicit user activation; never open links automatically.
- Show Mermaid source as a fallback; diagram rendering is a later possibility.
- Offer plain text output when writing to a pipe, without ANSI styling by default.

## Design principles

- Keep Markdown parsing, terminal layout, and the TUI separate.
- Do not execute embedded HTML or pass terminal control sequences from document content through to the terminal.
- Do not make network requests while rendering a document. Explicit activation may open an HTTP(S) link in the default browser.

## Project status

Go 1.25 is the approved minimum project toolchain. Goldmark is the approved Markdown parser and tcell is the approved TUI library. The current vertical slice uses tcell, `uniseg` for display width, and `x/term` for terminal detection; Goldmark is not wired in yet because Markdown rendering belongs to the next ticket. The clipboard package, installation method, and release process have not been selected. The distribution goal is to provide ready-to-run binaries so users do not need to install Go themselves.

See the [technical analysis and MVP proposal](docs/markdown-terminal-reader-analysis.md) for the REA observations, design decisions, open questions, and acceptance criteria.

See the [command specification](docs/command-spec.md) for the planned arguments, themes, and keyboard controls.

See the [acceptance criteria and Markdown fixtures](docs/acceptance-criteria.md) for how MVP behavior will be verified.

See the [development workflow](docs/development-workflow.md) for the recommended implementation phases and verification gates.

See the [development ticket checklists](docs/tickets/README.md) for agent-sized tasks, prerequisites, and acceptance checks.

---

<a id="reamd-bahasa-indonesia"></a>

# reamd (Bahasa Indonesia)

**Pembaca Markdown read-only di terminal.** `reamd` membuka satu dokumen di terminal yang sedang digunakan, dapat digulir, lalu pengguna kembali ke shell setelah keluar.

> [!NOTE]
> `reamd` masih dalam pengembangan awal. Command line, pemuatan file, TUI yang dapat digulir, dan output teks biasa ke pipe sudah berjalan. Render Markdown, tema, link, mouse, dan clipboard belum diimplementasikan.

## Penggunaan

```sh
# Buka satu file di terminal saat ini
reamd README.md

# Pilih tema (validasi sudah ada, palet masih dikerjakan)
reamd --theme transparent README.md

# Tampilkan bantuan atau versi
reamd --help
reamd --version
```

Jika dijalankan tanpa argumen, `reamd` menampilkan bantuan dan contoh penggunaan lalu keluar dengan kode 0. MVP menerima satu path file per sesi dan belum membaca stdin. Di terminal interaktif, aplikasi membuka TUI read-only; jika output dialihkan atau di-pipe, aplikasi mencetak teks biasa tanpa ANSI.

Tombol yang tersedia saat ini: `Up`/`Down` dan `j`/`k` menggulir satu baris, `PageUp`/`PageDown` menggulir satu layar, `Home`/`End` menuju awal atau akhir, `?` membuka bantuan, dan `q` atau `Esc` keluar. Terminal di bawah 40 kolom x 10 baris menampilkan pesan untuk memperbesar terminal.

## Status saat ini

Reader sudah merender elemen inti CommonMark berupa heading, paragraf, penekanan, inline code, line break, list, dan blockquote di TUI yang dapat digulir atau sebagai teks render biasa saat output dialihkan. Elemen parser lain tetap menjadi fallback teks aman sampai renderer khususnya dibuat.

```sh
# Jalankan dari source
go run ./cmd/reamd README.md

# Build binary lokal
go build -o reamd ./cmd/reamd
```

## Rancangan MVP

- Membaca file Markdown dan menampilkannya dalam TUI read-only yang dapat digulir.
- Mendukung CommonMark dan GitHub Flavored Markdown (GFM), termasuk heading, list, quote, tabel, link, dan code block. URL dan alamat email polos dikenali sebagai autolink.
- Menyesuaikan layout dengan lebar terminal. Tabel lebar beralih ke tampilan bertumpuk jika diperlukan.
- Menyediakan zoom di dalam aplikasi dengan mengubah layout dan wrapping. Fitur ini tidak mengubah ukuran font terminal yang sebenarnya.
- Mendukung tema `auto`, `dark`, `light`, dan `transparent`.
- Mendukung navigasi keyboard dan mouse, seleksi teks, serta copy teks hasil render ke clipboard sistem.
- Menampilkan alt text dan URL gambar tanpa mengunduh gambar. Link HTTP(S) dibuka di browser default hanya setelah diaktifkan pengguna, tidak otomatis.
- Menampilkan source Mermaid sebagai fallback; render diagram dapat dipertimbangkan kemudian.
- Menghasilkan teks biasa saat output diarahkan ke pipe, tanpa ANSI secara default.

## Prinsip teknis

- Pisahkan parsing Markdown, layout terminal, dan TUI.
- Jangan mengeksekusi HTML di dalam dokumen atau meneruskan escape sequence dari konten ke terminal.
- Jangan melakukan network access saat merender dokumen. Aktivasi eksplisit dapat membuka link HTTP(S) di browser default.

## Status proyek

Go 1.25 adalah versi minimum toolchain proyek yang disetujui. Goldmark disetujui sebagai parser Markdown dan tcell sebagai library TUI. Vertical slice saat ini memakai tcell, `uniseg` untuk lebar tampilan, dan `x/term` untuk deteksi terminal; Goldmark belum dihubungkan karena render Markdown menjadi lingkup ticket berikutnya. Package clipboard, cara instalasi, serta proses rilis belum ditentukan. Tujuan distribusinya adalah menyediakan binary siap pakai agar pengguna tidak perlu memasang Go sendiri.

Lihat [analisis teknis dan rancangan MVP](docs/markdown-terminal-reader-analysis.md) untuk hasil observasi REA, keputusan rancangan, pertanyaan terbuka, dan kriteria penerimaan.

Lihat [spesifikasi command](docs/command-spec.md) untuk rancangan argumen, tema, dan kontrol keyboard.

Lihat [kriteria penerimaan dan fixture Markdown](docs/acceptance-criteria.md) untuk cara memeriksa perilaku MVP.

Lihat [workflow pengembangan](docs/development-workflow.md) untuk tahapan implementasi dan gerbang verifikasi yang disarankan.

Lihat [checklist ticket pengembangan](docs/tickets/README.md) untuk tugas agent, prasyarat, dan pemeriksaan penerimaan.
