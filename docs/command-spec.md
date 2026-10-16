# Spesifikasi command `reamd`

Dokumen ini menetapkan kontrak penggunaan command dan kontrol TUI untuk MVP `reamd`. Dokumen ini melengkapi [analisis teknis](markdown-terminal-reader-analysis.md), bukan menggantikannya.

## Status

Keputusan di bawah berasal dari sesi tanya-jawab produk. Ini adalah spesifikasi perilaku yang dituju, bukan klaim bahwa command sudah diimplementasikan.

## Tujuan

`reamd` membuka satu dokumen sebagai pembaca Markdown read-only di terminal yang sedang digunakan. MVP memakai Go dan menargetkan CommonMark + GitHub Flavored Markdown (GFM). Aplikasi tidak mengedit atau menyimpan perubahan ke file.

## Sintaks command

```text
reamd [--theme <auto|dark|light|transparent>] <path>
reamd --help
reamd --version
```

### Contoh

```sh
# Buka dokumen dengan tema default
reamd README.md

# Pilih tema secara eksplisit
reamd --theme dark README.md
reamd --theme light README.md
reamd --theme transparent README.md

# Lihat bantuan atau versi
reamd --help
reamd --version
```

`reamd` tanpa argumen menampilkan bantuan dan contoh penggunaan, lalu selesai dengan sukses. MVP menerima **tepat satu path file**. File dibaca sebagai Markdown tanpa membatasi ekstensi, sehingga nama file tidak harus berakhiran `.md` atau `.markdown`.

## Argumen dan opsi

| Argumen/opsi | Perilaku |
|---|---|
| `<path>` | Wajib untuk membuka dokumen. Path relatif dihitung dari direktori kerja saat ini. Path absolut juga diterima. |
| `--theme <value>` | Memilih `auto`, `dark`, `light`, atau `transparent`. Default: `auto`. Nilai lain menghasilkan pesan kesalahan dan kode keluar non-zero. |
| `--help`, `-h` | Menampilkan sintaks, opsi, serta contoh command. Tidak membuka TUI. |
| `--version`, `-v` | Menampilkan versi aplikasi. Tidak membuka TUI. |

MVP tidak menerima banyak path sekaligus dan tidak membaca dokumen dari stdin. Input stdin dapat dipertimbangkan setelah pengalaman file tunggal stabil.

## Perilaku saat dijalankan

### Terminal interaktif

Jika path valid dan proses berjalan di terminal interaktif, `reamd` membuka TUI baca-saja di terminal tersebut. TUI menggunakan area layar terminal, bukan membuka jendela atau proses terminal baru.

- Dokumen dapat digulir dan mengikuti ukuran layar terminal.
- TUI menampilkan petunjuk singkat bahwa `?` membuka bantuan pintasan.
- `q` atau `Esc` menutup TUI dan mengembalikan pengguna ke shell.
- File tidak diubah. Link tidak dibuka otomatis; link web hanya dibuka setelah pengguna mengaktifkannya secara eksplisit.
- Jika jendela lebih kecil dari 40 kolom × 10 baris, tampilkan pesan ringkas untuk memperbesar terminal dan tetap sediakan cara keluar.
- Terminal harus dipulihkan ke kondisi normal ketika TUI ditutup, termasuk setelah terjadi kesalahan selama sesi.

### Output dialihkan atau di-pipe

Jika output tidak menuju terminal interaktif, `reamd <path>` tidak mencoba mengambil alih layar. Aplikasi mencetak versi teks biasa yang dirender ke stdout, tanpa warna/ANSI secara default. Ini memungkinkan hasil disimpan atau diteruskan ke program lain.

Contoh:

```sh
reamd README.md > README-rendered.txt
reamd README.md | less
```

Perilaku pipe berlaku untuk output. Input stdin tetap di luar cakupan MVP.

## Pintasan TUI

| Tombol | Perilaku |
|---|---|
| `Up` / `Down` | Gulir satu baris ke atas/bawah. |
| `j` / `k` | Gulir satu baris ke bawah/atas. |
| `PageUp` / `PageDown` | Gulir kira-kira satu tinggi layar. |
| `Home` / `End` | Pindah ke awal/akhir dokumen. |
| `q` / `Esc` | Keluar dari TUI dan kembali ke shell. |
| `?` | Tampilkan atau tutup panel bantuan pintasan. |
| `Tab` / `Shift+Tab` | Pindah fokus ke link berikutnya/sebelumnya. |
| `Enter` | Aktifkan link yang sedang fokus. Link internal melompat ke heading; link web membuka browser default jika protokolnya HTTP atau HTTPS. |
| `Ctrl+Plus` / `Ctrl+-` | Naikkan/turunkan zoom satu langkah, jika kombinasi tombol diteruskan terminal. |
| `+` / `-` | Alternatif untuk menaikkan/menurunkan zoom. |
| `Ctrl+0` | Kembalikan zoom ke nilai awal. |
| `Ctrl+C` | Salin teks yang sedang dipilih ke clipboard sistem; jika tidak ada pilihan, tampilkan petunjuk dan jangan keluar. |

Terminal tertentu menangani kombinasi Ctrl sendiri atau mengodekan `Ctrl+Plus` dengan cara berbeda. Karena itu, `+` dan `-` disediakan sebagai alternatif. Bantuan TUI perlu menjelaskan binding yang tersedia.

### Mouse dan clipboard

- Scroll wheel menggulir dokumen.
- Klik pada link mengaktifkan link. Link eksternal hanya membuka URL `http` atau `https`, melalui browser default, sebagai akibat langsung dari klik pengguna. Link dengan protokol lain tidak dibuka.
- Klik link anchor internal seperti `[Intro](#intro)` melompat ke heading tujuan. Slug heading mengikuti gaya GitHub: huruf kecil, spasi menjadi `-`, tanda baca dihapus, dan judul duplikat mendapat akhiran.
- Drag memilih teks yang terlihat. `Ctrl+C` menyalin teks hasil render ke clipboard sistem tanpa ANSI/style control codes.
- Jika akses clipboard gagal, tampilkan pesan yang menjelaskan kegagalan dan pertahankan pilihan agar pengguna dapat memakai cara copy dari terminal.
- Scrollbar dapat diklik untuk berpindah posisi dan thumb-nya dapat di-drag.
- Fitur mouse aktif pada terminal yang melaporkan dukungan mouse. Semua fungsi utama tetap memiliki kontrol keyboard.

### Zoom

- Nilai awal: 100%.
- Satu tekanan `+` atau `Ctrl+Plus`: naik 10 poin persentase.
- Satu tekanan `-` atau `Ctrl+-`: turun 10 poin persentase.
- Batas: 70% sampai 150%.
- `Ctrl+0`: reset ke 100%.

Zoom mengubah kepadatan layout, penekanan heading, jarak, dan wrapping di dalam grid karakter terminal. Zoom **tidak** mengubah ukuran font terminal sebenarnya. Saat mencapai batas, tekanan tambahan tidak mengubah nilai lebih jauh.

## Tema

| Tema | Perilaku |
|---|---|
| `auto` (default) | Mengikuti warna foreground/background terminal dan kemampuan warnanya. Tidak mengasumsikan semua terminal menggunakan latar gelap. |
| `dark` | Memakai palet teks dan aksen yang dirancang untuk latar gelap. |
| `light` | Memakai palet teks dan aksen yang dirancang untuk latar terang. |
| `transparent` | Tidak menggambar warna latar belakang, sehingga latar terminal tetap terlihat. Warna teks dan aksen masih boleh digunakan. |

`NO_COLOR` menonaktifkan warna ANSI. Terminal tanpa dukungan warna mendapat tampilan yang tetap terbaca melalui struktur, indentasi, dan penanda teks. Saat stdout dialihkan, warna dimatikan secara default, terlepas dari tema yang dipilih.

## Aturan render yang memengaruhi command

- Format sasaran: CommonMark + GFM.
- Tabel ditampilkan sebagai tabel jika muat. Jika tidak, tampilkan baris sebagai susunan label/nilai agar isi tidak terpotong.
- Baris code block yang lebih panjang dari lebar terminal di-wrap otomatis dan diberi penanda visual bahwa baris berlanjut.
- Link menampilkan label dan URL; program tidak membukanya otomatis.
- Gambar menampilkan alt text dan URL; program tidak melakukan fetch.
- Mermaid ditampilkan sebagai source code dengan penanda bahwa itu diagram. Render diagram bukan bagian MVP.
- Raw HTML tidak dieksekusi. Tag dan kontennya ditampilkan sebagai teks aman, tanpa mengirim control sequence ke terminal.
- Konten file tidak boleh dapat menyuntikkan ANSI/OSC atau karakter kontrol yang mengubah terminal.

## Bantuan dan kesalahan

`--help` dan pemanggilan tanpa argumen menampilkan bantuan ringkas yang memuat sintaks, opsi tema, pintasan TUI utama, dan contoh penggunaan.

Bahasa `--help`, panel bantuan, dan pesan error default adalah Inggris.

Jika path tidak ada, bukan file yang dapat dibaca, bukan teks UTF-8 yang dapat diproses, atau tampak seperti file biner, tampilkan pesan yang menyebut masalahnya. Kesalahan argumen/opsi serta kesalahan membaca file tidak boleh meninggalkan terminal dalam mode TUI. Kesalahan ditulis ke stderr. Kode keluar: `0` untuk sukses/help/version, `1` untuk gagal membaca atau menjalankan, dan `2` untuk penggunaan atau opsi yang salah.

## Di luar cakupan MVP

- Mengedit atau menyimpan dokumen.
- Membuka banyak dokumen dan berpindah antar-tab.
- Membaca dokumen dari stdin.
- Menjalankan raw HTML atau mengambil konten jaringan.
- Merender Mermaid atau mengambil gambar.
- Mengubah ukuran font aplikasi terminal.
- Menetapkan metode instalasi/distribusi dan perintah build, sampai rancangan packaging dipilih.
