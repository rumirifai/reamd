# Workflow pengembangan `reamd`

Dokumen ini menjabarkan cara mengubah rancangan `reamd` menjadi aplikasi Go secara bertahap. Setiap fase memiliki hasil yang dapat diperiksa dan gerbang sebelum fase berikutnya. Ini adalah workflow yang direncanakan, bukan laporan bahwa implementasi atau CI sudah tersedia.

## Istilah singkat

| Istilah | Arti sederhana |
|---|---|
| Dependency/library | Kode siap pakai yang ditambahkan ke proyek untuk menyediakan kemampuan tertentu. |
| Spike | Percobaan kecil untuk memastikan pilihan teknis bekerja sebelum dipakai dalam aplikasi utama. |
| Vertical slice | Potongan aplikasi kecil yang sudah berjalan dari input sampai tampilan, walau fiturnya masih terbatas. |
| CI | Layanan yang menjalankan build dan test secara otomatis ketika perubahan dikirim ke repository. |
| Conformance corpus | Kumpulan contoh resmi untuk memeriksa apakah parser mengikuti aturan Markdown yang dipilih. |
| Pseudo-terminal | Terminal tiruan yang dipakai test otomatis untuk mengirim tombol dan memeriksa tampilan TUI. |
| Smoke test | Pemeriksaan manual singkat untuk memastikan aplikasi benar-benar berjalan di lingkungan nyata. |

## Acuan dan keputusan produk

Sebelum mengerjakan fitur, baca:

1. `AGENTS.md` untuk aturan kerja agent.
2. `docs/command-spec.md` untuk kontrak command, TUI, tema, input, dan interaksi.
3. `docs/acceptance-criteria.md` untuk kondisi lulus dan fixture.
4. `docs/markdown-terminal-reader-analysis.md` untuk latar riset dan alasan teknis.

Jika ada ketidaksesuaian, jangan memilih diam-diam. Tunjukkan bagian yang bertentangan dan minta keputusan pengguna sebelum mengubah perilaku produk. Gunakan grill-me/question mode untuk keputusan produk atau teknis yang belum dijawab; rangkum hasil yang disepakati sebelum implementasi terkait.

## Ringkasan fase

Checklist tugas agent per fase tersedia di [docs/tickets/README.md](tickets/README.md). Workflow ini menjelaskan urutan dan gerbang fase; ticket tersebut memecah pekerjaan menjadi unit yang dapat dikerjakan dalam satu session.

| Fase | Fokus | Gerbang selesai |
|---|---|---|
| 0 | Persiapan tugas dan ruang lingkup | Acceptance criteria yang dikerjakan dan keputusan terbuka teridentifikasi |
| 1 | Spike teknis dan pemilihan library | Parser, TUI, clipboard, dan pembukaan link memenuhi kebutuhan kritis |
| 2 | Fondasi Go dan CLI | Aplikasi dapat dibangun, menampilkan bantuan/error, dan membaca satu file |
| 3 | Vertical slice TUI | File tampil di terminal interaktif, bisa digulir, keluar dengan aman, dan punya fallback pipe |
| 4 | Parsing dan layout Markdown | Fixture dan conformance test CommonMark/GFM lulus |
| 5 | Interaksi TUI | Tema, zoom, resize, mouse, link, dan clipboard memenuhi spesifikasi |
| 6 | Verifikasi lintas platform | Test otomatis dan smoke test manual selesai pada matriks sasaran |
| 7 | Distribusi | Binary dan petunjuk instalasi siap sesuai metode yang dipilih |

## Fase 0: siapkan pekerjaan agent

Untuk setiap tugas coding:

1. Periksa `git status` jika repo sudah memakai Git; jangan menimpa perubahan yang bukan bagian tugas.
2. Pilih satu tujuan kecil, misalnya `AC-CLI` atau satu kelompok fixture `FX-*`, bukan seluruh MVP sekaligus.
3. Catat acceptance criteria yang akan dipenuhi dan file yang diperkirakan berubah.
4. Baca bagian spesifikasi yang berhubungan dengan tugas.
5. Jika ada keputusan terbuka, ajukan pertanyaan sebelum mengerjakan kode yang bergantung padanya.

**Gerbang:** tujuan tugas punya hasil yang dapat diuji, tidak bertentangan dengan spesifikasi, dan tidak mencakup fitur di luar MVP tanpa persetujuan.

## Fase 1: spike teknis dan pemilihan library

Spike adalah percobaan kecil untuk menguji risiko sebelum memilih dependency atau membangun arsitektur besar. Go sudah dipilih, tetapi parser, TUI, clipboard, dan integrasi browser belum diputuskan.

Bandingkan kandidat dan uji hal berikut:

- **Parser Markdown:** CommonMark + GFM penuh, conformance corpus, link internal, raw HTML yang dapat ditampilkan sebagai teks aman, dan kemampuan menjaga isi yang tidak punya renderer khusus.
- **Framework TUI:** key handling, resize, mouse wheel, scrollbar, drag selection, pemulihan terminal, dan dukungan Windows/macOS/Linux.
- **Clipboard:** copy teks hasil render tanpa ANSI pada ketiga sistem operasi, termasuk perilaku ketika clipboard tidak tersedia.
- **Browser default:** hanya buka HTTP/HTTPS setelah aksi pengguna; argumen URL tidak boleh dieksekusi sebagai perintah shell.
- **Layout terminal:** lebar Unicode, wrapping list/code, tabel stacked, tema, dan output non-interaktif.
- **Kesehatan dependency:** pemeliharaan, lisensi, dependency transitif, dan kecocokan versi Go.

> Jangan memilih library hanya karena populer. Nilai kandidat berdasarkan fixture dan AC yang relevan. Catat pilihan, alasan, keterbatasan, dan versi yang diuji di dokumen terkait sebelum mengikat implementasi ke API library tersebut.

Agent menyajikan rekomendasi dependency beserta trade-off dalam Bahasa Indonesia dan menunggu persetujuan pengguna sebelum menambahkan dependency produksi.

**Gerbang:** bukti spike menunjukkan kebutuhan kritis bisa dilakukan atau ada keputusan produk yang disepakati untuk mengubah scope. Jangan mulai implementasi luas jika mouse selection, clipboard, atau pemulihan terminal belum terverifikasi.

## Fase 2: fondasi Go dan CLI

Setelah stack disetujui:

1. Buat `go.mod` dengan versi Go minimum yang dipilih dan dependency langsung yang diperlukan saja.
2. Tambahkan entry point executable, parsing argumen, `--help`, `--version`, validasi satu path, `--theme`, dan kode keluar sesuai `docs/command-spec.md`.
3. Pisahkan command parsing, pembacaan file, model dokumen, layout, dan tampilan agar logika bisa dites tanpa membuka TUI.
4. Tambahkan CI lintas Windows, macOS, dan Linux sedini mungkin.
5. Dokumentasikan perintah build/test hanya setelah perintah tersebut berjalan di repo.

Path yang tidak ada, tidak dapat dibaca, atau konten bukan UTF-8/biner harus menghasilkan pesan yang jelas sebelum TUI mengambil alih terminal.

**Gerbang:** bantuan, versi, opsi salah, path valid, path invalid, dan status keluar memenuhi `AC-CLI`.

## Fase 3: vertical slice TUI

Implementasikan alur terkecil dari command sampai pengalaman yang bisa dipakai:

1. Baca satu file.
2. Tampilkan dokumen dasar di TUI dalam terminal yang sama.
3. Dukung scroll keyboard, resize, `q`/`Esc`, bantuan `?`, dan TUI minimal 40×10.
4. Pulihkan mode terminal ketika pengguna keluar, menerima interupsi, atau terjadi error.
5. Saat stdout bukan terminal interaktif, render teks biasa tanpa TUI/warna secara default.

> Pertahankan pemisahan tiga bagian: parser menghasilkan struktur dokumen, layout menghitung baris/kolom, dan TUI menangani input serta menggambar tampilan.

**Gerbang:** skenario dasar `reamd <path>` lulus, output pipe tidak mencoba membuka TUI, dan terminal tidak tertinggal dalam mode raw/alternate screen.

## Fase 4: parsing dan layout Markdown

Implementasikan parser terhadap model dokumen dan layout terminal. Gunakan fixture bernama di `docs/acceptance-criteria.md` sebagai dasar, lalu saat coding pindahkan contoh menjadi file test yang dapat dibaca test otomatis.

- Jalankan conformance corpus CommonMark dan GFM sesuai versi library yang dipilih.
- Tambahkan unit test untuk node Markdown, wrapping, nested indent, tabel, tema, dan lebar Unicode.
- Pastikan raw HTML tidak dieksekusi dan ditampilkan sebagai teks aman.
- Pastikan karakter kontrol dinetralkan; jangan kirim escape sequence dari dokumen ke terminal.
- Jangan fetch gambar, membaca jaringan saat render, atau membuka link otomatis.
- Uji lebar 40, 80, dan 120 kolom, termasuk tabel sempit, code panjang, kombinasi aksen, CJK, dan emoji.

**Gerbang:** `AC-MD`, `AC-RENDER`, `AC-SEC`, dan fixture yang berkaitan lulus tanpa mengandalkan snapshot warna yang identik.

## Fase 5: interaksi TUI

Kerjakan perilaku interaktif dalam unit kecil dan hubungkan masing-masing dengan acceptance criteria:

1. Binding keyboard dan panel bantuan.
2. Zoom 70%-150%, langkah 10%, alternatif `+`/`-`, reset `Ctrl+0`.
3. Resize, tabel responsif, dan tema `auto`, `dark`, `light`, `transparent` dengan `NO_COLOR`.
4. Mouse wheel dan scrollbar.
5. Seleksi teks dan copy teks hasil render ke clipboard dengan `Ctrl+C`.
6. Navigasi link via `Tab`/`Shift+Tab`/`Enter`, klik anchor internal bergaya GitHub, dan buka HTTP/HTTPS eksternal hanya setelah aktivasi pengguna.

Jika browser atau clipboard tidak tersedia, tampilkan pesan yang jelas dan jangan crash atau meninggalkan TUI dalam kondisi rusak.

**Gerbang:** tiap interaksi lulus test otomatis yang sesuai atau smoke test manual pada terminal yang benar-benar mendukung kapabilitas tersebut.

## Fase 6: verifikasi lintas platform

### Otomatis

- Jalankan parser, layout, input file, sanitasi, status keluar, help/version, dan test tanpa TUI pada CI Windows/macOS/Linux.
- Gunakan test terminal/pseudo-terminal untuk alur keyboard dan restore terminal jika library test mendukungnya.
- Verifikasi dependency dan perubahan `go.mod`/`go.sum`.

### Manual

- Jalankan smoke test Windows Terminal, Terminal.app, serta GNOME Terminal atau emulator Linux kompatibel.
- Periksa keyboard, scroll wheel, resize, zoom, tema, seleksi/copy, navigasi anchor, pembukaan link, serta exit setelah error/interupsi.
- Catat platform/terminal dan hasilnya. Jangan menyebut kapabilitas manual lulus jika belum benar-benar dicoba.

**Gerbang:** seluruh kriteria wajib `docs/acceptance-criteria.md` lulus atau ada kegagalan yang disepakati secara eksplisit sebagai blocker sebelum rilis.

## Fase 7: distribusi dan instalasi

Sasaran produk adalah pengguna dapat memasang dari terminal dan menjalankan tanpa memasang Go. Sebelum mengumumkan instalasi:

1. Pilih metode distribusi yang disetujui, seperti binary GitHub Release dan/atau package manager.
2. Build dan uji binary untuk setiap OS/arsitektur yang diumumkan.
3. Sertakan versi, checksum, dan instruksi install/uninstall yang benar-benar diuji.
4. Jangan menganggap `go install` memenuhi sasaran pengguna umum karena command itu membutuhkan toolchain Go pada mesin pengguna.
5. Perbarui README, help/version, serta instruksi agent dengan metode dan command yang sudah diverifikasi.

**Gerbang:** fresh install di mesin/lingkungan bersih menghasilkan command `reamd` yang berjalan, dan setiap klaim dukungan release memiliki artifact yang sesuai.

## Siklus kerja agent untuk satu perubahan

Untuk setiap unit kerja:

1. Sebutkan acceptance ID yang dituju dan keputusan yang sudah ada.
2. Jika ada ambiguitas, tanyakan dahulu dan tunggu persetujuan.
3. Buat perubahan terkecil yang memenuhi kriteria; hindari refactor yang tidak dibutuhkan.
4. Tambah atau perbarui test/fixture bersamaan dengan perilaku yang diubah.
5. Jalankan formatter, test terkait, seluruh test, vet, dan build yang tersedia.
6. Periksa diff untuk perubahan tak sengaja dan pastikan dokumen kontrak tetap selaras.
7. Laporkan file yang berubah, acceptance ID lulus/gagal, command yang dijalankan, hasil, serta pemeriksaan manual yang belum dilakukan.

Setelah `go.mod` tersedia, baseline check yang direncanakan adalah:

| Tujuan | Command |
|---|---|
| Format file Go yang diubah | `gofmt -w <file.go>` |
| Semua test | `go test ./...` |
| Pemeriksaan statis Go | `go vet ./...` |
| Build semua package | `go build ./...` |

Sesuaikan command jika struktur proyek atau CI menetapkan cara lain; `AGENTS.md` hanya boleh mencatat command yang sudah diverifikasi.

## Pemeliharaan dokumentasi

- Perubahan pada sintaks atau perilaku command harus memperbarui `docs/command-spec.md`.
- Perubahan pada kondisi lulus atau fixture harus memperbarui `docs/acceptance-criteria.md`.
- Perubahan dependency/arsitektur harus memperbarui analisis teknis dan alasan pemilihan.
- Perubahan instalasi atau penggunaan harus memperbarui kedua bahasa di `README.md`.
- Perubahan command build/test yang telah diverifikasi harus dicatat di `AGENTS.md`.
