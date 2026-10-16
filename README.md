# reamd

**A read-only Markdown reader for your terminal.** `reamd` is planned as a TUI: it opens a document in the current terminal, lets you scroll through it, and returns you to the shell when you quit.

> [!NOTE]
> `reamd` is in the planning stage. The commands and features below describe the intended MVP and are not available yet.

**Languages:** English | [Bahasa Indonesia](#reamd-bahasa-indonesia)

## Planned usage

```sh
# Open one Markdown file in the current terminal
reamd README.md

# Choose a theme
reamd --theme transparent README.md

# Show usage or version information
reamd --help
reamd --version
```

With no arguments, `reamd` is planned to show help and usage examples. It accepts one file path per run and does not read from stdin in the MVP. In an interactive terminal it opens a read-only TUI; when output is redirected or piped, it prints plain rendered text instead.

## Planned MVP

- Read a Markdown file and show it in a scrollable, read-only TUI.
- Support the full CommonMark and GitHub Flavored Markdown (GFM) syntax, including headings, lists, quotes, tables, links, and code blocks.
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

Go has been selected as the implementation language. The parser and TUI libraries, installation method, and build commands have not been selected yet. The distribution goal is to provide ready-to-run binaries so users do not need to install Go themselves.

See the [technical analysis and MVP proposal](docs/markdown-terminal-reader-analysis.md) for the REA observations, design decisions, open questions, and acceptance criteria.

See the [command specification](docs/command-spec.md) for the planned arguments, themes, and keyboard controls.

See the [acceptance criteria and Markdown fixtures](docs/acceptance-criteria.md) for how MVP behavior will be verified.

See the [development workflow](docs/development-workflow.md) for the recommended implementation phases and verification gates.

See the [development ticket checklists](docs/tickets/README.md) for agent-sized tasks, prerequisites, and acceptance checks.

---

<a id="reamd-bahasa-indonesia"></a>

# reamd (Bahasa Indonesia)

**Pembaca Markdown read-only di terminal.** `reamd` direncanakan sebagai TUI: dokumen dibuka di terminal yang sedang digunakan, dapat digulir, lalu pengguna kembali ke shell setelah keluar.

> [!NOTE]
> `reamd` masih dalam tahap perencanaan. Contoh command dan fitur di bawah ini adalah rancangan MVP, belum tersedia untuk digunakan.

## Rancangan penggunaan

```sh
# Buka satu file di terminal saat ini
reamd README.md

# Pilih tema
reamd --theme transparent README.md

# Tampilkan bantuan atau versi
reamd --help
reamd --version
```

Jika dijalankan tanpa argumen, `reamd` direncanakan menampilkan bantuan dan contoh penggunaan. MVP menerima satu path file per sesi dan belum membaca stdin. Di terminal interaktif, aplikasi membuka TUI read-only; jika output dialihkan atau di-pipe, aplikasi mencetak teks biasa yang sudah dirender.

## Rancangan MVP

- Membaca file Markdown dan menampilkannya dalam TUI read-only yang dapat digulir.
- Mendukung seluruh sintaks CommonMark dan GitHub Flavored Markdown (GFM), termasuk heading, list, quote, tabel, link, dan code block.
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

Go telah dipilih sebagai bahasa implementasi. Library parser dan TUI, cara instalasi, serta perintah build belum ditentukan. Tujuan distribusinya adalah menyediakan binary siap pakai agar pengguna tidak perlu memasang Go sendiri.

Lihat [analisis teknis dan rancangan MVP](docs/markdown-terminal-reader-analysis.md) untuk hasil observasi REA, keputusan rancangan, pertanyaan terbuka, dan kriteria penerimaan.

Lihat [spesifikasi command](docs/command-spec.md) untuk rancangan argumen, tema, dan kontrol keyboard.

Lihat [kriteria penerimaan dan fixture Markdown](docs/acceptance-criteria.md) untuk cara memeriksa perilaku MVP.

Lihat [workflow pengembangan](docs/development-workflow.md) untuk tahapan implementasi dan gerbang verifikasi yang disarankan.

Lihat [checklist ticket pengembangan](docs/tickets/README.md) untuk tugas agent, prasyarat, dan pemeriksaan penerimaan.
