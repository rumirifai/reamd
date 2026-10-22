# Analisis: Markdown Live Preview dan rancangan `reamd`

Tanggal observasi: 16 Oktober 2026  
Target: <https://markdownlivepreview.com/>  
Metode: REA `capture_browser_scenario`, Chromium 155.0.8059.41, headless, viewport 1280×720, tanpa storage awal. Evidence: `ev_0fd79ea5b05880b1f34d44b03d4067e2936468697433403b73e1b624c8be7637`.

## Ringkasan

`reamd` direncanakan sebagai pembaca Markdown **read-only** dengan antarmuka TUI. Saat dijalankan di terminal interaktif, aplikasi mengambil alih area terminal untuk menampilkan dokumen dan menerima input keyboard. Pengguna dapat menggulir dokumen, mengubah skala tampilan aplikasi, lalu keluar kembali ke shell. `reamd` tidak menjadi editor Markdown pada MVP.

Website yang dianalisis adalah aplikasi client-side dengan editor Monaco di sebelah kiri dan hasil render HTML di sebelah kanan. DOM hasil observasi menunjukkan `#output` berisi HTML hasil Markdown; contoh default mencakup heading, emphasis, list, image, link, blockquote, tabel, fenced code, Mermaid, dan inline code. Header memiliki Open, Copy, Export PDF, Reset, Sync scroll, dan Dark mode. File Open menerima `.md`, `.markdown`, dan `text/markdown`. UI menyebut `markedjs/marked`; halaman menampilkan diagram Mermaid sebagai SVG.

Untuk `reamd`, kebutuhan utamanya adalah parsing Markdown dan penyajian yang nyaman di terminal. Monaco, editor dua-pane, HTML/CSS, dan ekspor PDF bukan bagian dari MVP. TUI tetap menggunakan terminal yang sedang dibuka, bukan membuka jendela browser atau terminal baru.

## Observasi dan batas bukti

| Temuan | Dasar | Kepastian |
|---|---|---|
| Website menampilkan editor Monaco dan preview berdampingan. | DOM hasil REA | Teramati |
| Preview berupa HTML semantik di `#output.content.markdown-body`. | DOM hasil REA | Teramati |
| Contoh mencakup heading, emphasis, list, blockquote, tabel, code, link, image, dan Mermaid. | DOM hasil REA | Teramati |
| File Open menerima `.md`, `.markdown`, dan `text/markdown`; terdapat juga kontrol Copy, Export PDF, Reset, Sync scroll, dan Dark mode. | DOM hasil REA | UI teramati; fungsi kontrol tidak diuji |
| Render lokal melalui CLI dapat dilakukan tanpa meniru layanan web. | Struktur yang terlihat dan tujuan proyek | Inferensi kuat |

Capture hanya memuat halaman awal dan sample bawaan. Tombol, upload file, perubahan isi editor, source bundle, dan request saat interaksi tidak diuji. Karena itu, perilaku sinkronisasi scroll, sanitasi, konfigurasi parser, dan ekspor PDF website tidak disimpulkan dari observasi ini.

## Keputusan produk saat ini

| Topik | Keputusan atau status |
|---|---|
| Bentuk aplikasi | TUI pembaca read-only, pilihan A dari sesi tanya-jawab |
| Editing dokumen | Tidak termasuk MVP; file diedit dengan editor lain |
| Format Markdown | CommonMark + GFM sebagai target MVP |
| Input | Tepat satu path file; stdin tidak termasuk MVP |
| Command tanpa argumen | Menampilkan bantuan dan contoh penggunaan |
| Tabel lebar | Rekomendasi: tabel biasa jika muat, tampilan bertumpuk jika terminal terlalu sempit |
| Zoom | Mulai 100%, langkah 10%, batas 70%-150%, `=`/`+` memperbesar, `-` memperkecil, `0` mereset; mengubah layout aplikasi, bukan ukuran font terminal |
| Tema | `auto` (default), `dark`, `light`, `transparent`; rincian ada di spesifikasi command |
| Bahasa | Dokumen kerja dalam Bahasa Indonesia; README untuk pengguna GitHub disiapkan dalam Bahasa Inggris dan Bahasa Indonesia |
| Bahasa pemrograman | Go dipilih sebagai bahasa implementasi |
| Parser dan TUI | Goldmark dan `tcell/v2` disetujui untuk implementasi setelah evaluasi Ticket 01 |
| Toolchain minimum | Go 1.25 disetujui; pasangan patch minimum dependency transitif yang diperiksa dicatat di evaluasi stack |
| Linkify dan browser | URL/email polos dikenali sebagai Linkify; helper `github.com/pkg/browser` disetujui dengan validasi HTTP(S) dan aktivasi eksplisit |
| Clipboard | Adapter internal dengan error; pengguna menyetujui `github.com/atotto/clipboard` sebagai implementasi produksi pada Ticket 09 |
| Distribusi | Sasaran: pengguna memasang melalui command terminal dan menjalankan aplikasi tanpa setup developer |

Spesifikasi command dan kontrol TUI yang telah dikonfirmasi tersedia di [command-spec.md](command-spec.md). Kriteria penerimaan dan fixture untuk memeriksa seluruh batas MVP tersedia di [acceptance-criteria.md](acceptance-criteria.md).

## CommonMark + GFM dan cakupan format

CommonMark adalah spesifikasi aturan dasar Markdown yang membantu parser berbeda memahami sintaks umum secara konsisten. GFM (GitHub Flavored Markdown) memperluasnya dengan format yang lazim di GitHub, antara lain tabel, task list, strikethrough, dan autolink.

Menyatakan format yang didukung berarti mendokumentasikan sintaks mana yang dirender secara khusus dan apa fallback-nya jika sintaks tidak didukung. Memilih CommonMark + GFM memberi dasar yang luas, tetapi fitur aplikasi lain seperti diagram Mermaid atau footnote tidak otomatis menjadi bagian GFM. MVP perlu memiliki daftar perilaku yang eksplisit:

| Konten | Perilaku MVP yang direncanakan |
|---|---|
| Heading, emphasis, paragraf, line break, horizontal rule | Render dengan hierarki, indentasi, dan wrapping terminal |
| Ordered/unordered list dan nesting | Render dengan marker dan indentasi yang konsisten |
| Blockquote | Render dengan indentasi/penanda quote |
| GFM table dan task list | Render tabel dan status task; tabel beralih ke tampilan bertumpuk jika ruang terlalu sempit |
| Fenced code dan inline code | Render sebagai teks kode; syntax highlighting belum menjadi syarat MVP |
| Link dan autolink | Tampilkan label dan alamat; link HTTP(S) hanya terbuka setelah aktivasi eksplisit. Anchor internal memakai gaya slug GitHub |
| Image | Tampilkan alt text serta URL; tidak mengunduh gambar |
| Raw HTML di Markdown | Tidak dieksekusi; tag ditampilkan sebagai teks aman |
| Mermaid | Tampilkan source code dengan penanda diagram; render diagram menjadi fitur lanjutan |
| Ekstensi lain di luar CommonMark/GFM | Tidak dijanjikan dukungannya; fallback harus aman dan didokumentasikan |

Kriteria dan fixture untuk perilaku ini ada di [acceptance-criteria.md](acceptance-criteria.md). Mermaid tetap berupa source fallback pada MVP.

## Pengalaman TUI yang direncanakan

Alur utama:

1. Pengguna menjalankan `reamd README.md`.
2. Aplikasi membaca file dan menampilkan dokumen dalam layar TUI pada terminal yang sama.
3. Pengguna menggulir dengan keyboard atau mouse; `q`/`Esc` menutup tampilan dan mengembalikan kontrol ke shell.
4. Pengguna dapat menavigasi link dengan keyboard atau mouse, memilih teks dan menyalinnya, serta menggunakan scrollbar mouse di terminal yang mendukungnya.
5. Pengguna dapat memperbesar atau memperkecil **skala tampilan di dalam aplikasi** dengan tombol `=`/`+` dan `-`; `=` dan `-` tidak memerlukan modifier. Kombinasi Ctrl dapat menjadi alternatif jika terminal meneruskannya.

Zoom TUI tidak dapat menjamin perubahan ukuran font terminal. Ukuran font dikontrol aplikasi terminal dan pintasan zoom-nya bervariasi. Zoom `reamd` berarti menyesuaikan penekanan heading, kepadatan, jarak, dan wrapping konten dalam grid karakter yang sama. Binding utama `=`/`-` tidak membutuhkan modifier; `+` dan kombinasi Ctrl tetap menjadi alternatif bila terminal meneruskannya.

Saat stdout bukan terminal interaktif, aplikasi mengeluarkan render teks biasa tanpa kontrol TUI atau warna secara default. Ini membuat output dapat dipipe ke program lain dan disimpan ke file. Membaca dokumen dari stdin tidak termasuk MVP.

## Arsitektur yang disarankan

`CLI -> file loader -> parser Markdown (AST) -> model dokumen -> layout berdasarkan ukuran terminal dan zoom -> TUI renderer -> terminal`

AST (*Abstract Syntax Tree*) adalah struktur data yang mewakili elemen Markdown, seperti heading, paragraf, tabel, dan code block. Dengan struktur ini, program tidak perlu menebak bentuk dokumen dari teks mentah menggunakan regex.

- Pisahkan parsing, layout, dan tampilan TUI supaya tiap bagian dapat diuji dan diubah sendiri.
- Pilih parser dengan dukungan seluruh CommonMark/GFM yang sesuai dengan bahasa pemrograman final.
- Jangan melakukan network access untuk merender dokumen, mengunduh gambar, atau membuka link secara otomatis. Link HTTP(S) hanya dibuka lewat browser default setelah tindakan pengguna.
- Netralisasi karakter kontrol terminal dari isi file agar dokumen tidak bisa menyuntikkan ANSI/OSC atau mengubah state terminal.
- Jangan mengeksekusi raw HTML.
- Deteksi ukuran terminal dan hitung lebar karakter Unicode dengan benar; jumlah karakter tidak selalu sama dengan jumlah kolom visual.
- TUI hanya digunakan saat berjalan di terminal interaktif. Saat output dialihkan, gunakan renderer teks biasa.
- Hormati `NO_COLOR`, terminal tanpa dukungan ANSI, dan fallback ASCII bila karakter box-drawing tidak tersedia.

### Kandidat bahasa pemrograman

Go dipilih sebagai bahasa implementasi; minimum toolchain Go 1.25 telah disetujui. Goldmark disetujui sebagai parser Markdown dan `tcell/v2` sebagai framework TUI berdasarkan bukti spike Ticket 01. Program dapat dibuild menjadi binary untuk sistem operasi sasaran, sehingga pengguna tidak perlu memasang Go untuk menjalankan binary yang sudah didistribusikan. Tampilan akhir tetap ditentukan oleh desain TUI, layout terminal, dan parser. Pada Ticket 09, pengguna menyetujui `github.com/atotto/clipboard` di balik adapter internal yang mengembalikan error; persetujuan dependency tidak dengan sendirinya membuktikan clipboard OS lulus smoke test.

Keputusan Ticket 09: clipboard sistem ditulis sebagai plain text melalui adapter `internal/clipboard`, dengan `github.com/atotto/clipboard` sebagai package produksi. Test aplikasi menyuntikkan adapter tiruan untuk menguji sukses, tidak tersedia, dan gagal tanpa mengakses clipboard pengguna.

## Risiko teknis dan mitigasi

- **Perbedaan lebar Unicode**: gunakan perhitungan lebar terminal yang memahami karakter Unicode, bukan panjang string biasa.
- **Wrapping di dalam list dan quote**: layout perlu mempertahankan indentasi konteks saat baris berlanjut.
- **Kontrol terminal berbahaya**: buang atau tampilkan dengan aman escape sequence dari dokumen; jangan meneruskannya ke terminal.
- **Tabel sempit**: beralih ke tampilan bertumpuk agar informasi tidak terpotong.
- **Perbedaan dukungan parser**: simpan fixture untuk CommonMark/GFM dan dokumentasikan fallback ekstensi.
- **Pintasan zoom**: terminal dapat menangani kombinasi tombol tertentu sendiri; sediakan binding alternatif dan opsi keluar.
- **Mode TUI dan pipe**: deteksi apakah proses memiliki terminal interaktif; gunakan keluaran teks biasa jika tidak.
- **Mode raw terminal**: ketika TUI memakai input tombol langsung, terminal harus dipulihkan ke kondisi awal saat keluar normal maupun saat terjadi error.

## Tahapan pengembangan

Urutan implementasi dan gerbang verifikasi dijabarkan di [development-workflow.md](development-workflow.md). Dokumen workflow itu menjadi panduan operasional agar setiap fase selesai sebelum dependensinya dibangun lebih jauh.

## Kriteria penerimaan awal

- `reamd README.md` membuka TUI baca-saja di terminal yang sama, tanpa mengubah file atau melakukan network access.
- TUI dapat digulir, `q` mengembalikan kontrol ke shell, dan terminal dipulihkan setelah aplikasi keluar.
- Kombinasi zoom dan alternatifnya, jika dibutuhkan, mengubah skala layout aplikasi tanpa mengklaim mengubah ukuran font terminal.
- Konten CommonMark/GFM yang disepakati memiliki fixture dan render yang dapat diperiksa.
- Tabel menyesuaikan ruang dan beralih ke tampilan bertumpuk saat terminal terlalu sempit.
- Input Unicode dan wrapping diuji pada lebar terminal berbeda.
- Kontrol terminal dari isi Markdown tidak lolos sebagai escape sequence; link dan gambar tidak diambil otomatis.
- Saat output diarahkan ke pipe, tidak ada layar TUI atau ANSI secara default.
- Parser, layout, dan TUI dapat diuji secara terpisah.

## Evaluasi stack — Ticket 01

Hasil spike, versi yang diperiksa, kandidat parser/TUI/clipboard/browser, hasil corpus, batas keamanan, dependency, trade-off, dan acceptance yang belum terbukti dicatat di [evaluasi stack](stack-evaluation.md).

Ringkasnya: pengguna menyetujui Goldmark + `tcell/v2`, helper `github.com/pkg/browser` dengan validasi HTTP(S) dan aktivasi eksplisit, serta clipboard melalui abstraksi ber-error dengan package final ditunda. Minimum Go 1.24 yang disetujui pada 16 Oktober 2026 diperbarui menjadi Go 1.25 pada 19 Oktober 2026 setelah audit menemukan advisory pada dependency transitif; pasangan perbaikan minimum berhasil diuji di scratch. Go 1.27.2 terpasang. Kedua parser cocok dengan corpus CommonMark 0.31.2 pada konfigurasi pembanding. Goldmark GFM bundle diuji pada 648 contoh GFM spec 0.29 dan menghasilkan 14 exact-output mismatch; setelah normalisasi semantik yang tercatat, hasilnya 0 semantic mismatch. Pengguna menyetujui Linkify URL/email polos sebagai ekstensi produk. Renderer HTML default Goldmark menghilangkan raw HTML sehingga aplikasi tetap memerlukan renderer AST aman. Kedua kandidat TUI berhasil cross-compile untuk Windows/macOS/Linux amd64; test simulasi tcell terbatas lulus. Pemeriksaan terakhir prototipe seleksi di Windows dengan Go 1.27.2 menjalankan `go test .` dan lulus (`ok stackspike`), memakai event mouse simulasi serta clipboard tiruan; ini bukan bukti drag selection di terminal nyata atau keberhasilan clipboard OS. Belum ada smoke test terminal nyata, clipboard/browser OS, atau restore terminal. Ticket 01 telah ditutup setelah risiko yang dapat mengubah pilihan stack memiliki bukti dan persetujuan; Ticket 02 berstatus ready-for-agent dan belum dimulai. Tidak ada dependency produksi atau `go.mod` yang ditambahkan ke root repo.

## Kesimpulan

Tujuan produk berubah dari renderer teks langsung menjadi TUI pembaca read-only. Keputusan ini tidak mengubah inti parsing Markdown, tetapi menambah pengelolaan keyboard, scroll, ukuran terminal, zoom layout, serta pemulihan state terminal. TUI dapat dibuat tanpa menanam browser. MVP tetap berfokus pada CommonMark/GFM, pembacaan yang nyaman, tabel responsif, dan perilaku aman terhadap isi file. Ticket 01 telah memenuhi gerbang Fase 1 dan ditutup dengan stack yang disetujui. Acceptance aplikasi yang belum diuji—termasuk renderer aman, interaksi TUI nyata, clipboard OS, browser OS, dan smoke test lintas platform—tetap terbuka pada ticket implementasi/verifikasi terkait dan tidak dinyatakan lulus berdasarkan spike.
