# Ocean — Arahan Desain TUI `reamd`

Dokumen ini mengadaptasi arah visual Ocean dari `cola-runner/awesome-tui-design` untuk pembaca Markdown `reamd`. Spesifikasi produk tetap menjadi sumber kebenaran untuk perilaku, tema, ukuran terminal minimum, dan aksesibilitas. Arahan ini tidak mengubah kontrak produk.

## Arah visual

- **Suasana:** tenang, fokus, profesional, dan mendukung membaca dalam waktu lama.
- **Kepadatan:** seimbang. Prioritaskan isi dokumen dan hierarki yang jelas daripada hiasan atau ruang kosong yang berlebihan.
- **Karakter:** aksen laut biru dan hijau lembut memberi identitas; struktur teks tetap menjadi pembeda utama.
- **Audiens dan konteks:** pengguna yang membaca Markdown di terminal.

## Prinsip produk yang mengikat

- TUI harus tetap dapat digunakan pada ukuran minimal **40 kolom × 10 baris**. Rekomendasi terminal yang lebih lebar boleh diberikan, tetapi bukan syarat operasi.
- Tema `auto`, `dark`, `light`, dan `transparent` mengikuti `docs/command-spec.md`. Ocean bukan alasan untuk memaksakan latar gelap.
- Warna bergantung pada kapabilitas terminal. Hormati `NO_COLOR`, terminal dengan warna terbatas, dan terminal tanpa warna. Tidak boleh mensyaratkan 256 warna atau TrueColor.
- Informasi penting tidak boleh dibedakan hanya melalui warna. Gunakan struktur, penanda teks, atau gaya terminal yang tersedia sebagai fallback.
- Isi Markdown adalah fokus utama. Jangan menambahkan elemen dekoratif yang mengurangi ruang baca atau mengubah dokumen sumber.
- Fitur utama harus tetap dapat digunakan dengan keyboard. Mouse hanya aktif saat terminal melaporkan kapabilitas yang sesuai.

## Warna

Warna Ocean berikut adalah referensi untuk tema gelap dan terminal yang mendukung warna, bukan warna wajib untuk semua tema.

| Peran | Referensi | Penggunaan yang disarankan |
|---|---|---|
| Latar dalam | `#0b1929` | Hanya pada tema yang memang memakai latar gelap; jangan diterapkan pada `auto` atau `transparent` secara paksa. |
| Teks utama | `#c8d6e5` | Teks dokumen pada palet gelap jika kontrasnya memadai. |
| Aksen utama | `#48dbfb` | Fokus, petunjuk interaksi, atau penanda seleksi. |
| Aksen pendamping | `#55efc4` | Penekanan sekunder yang jarang digunakan. |
| Teks redup | `#576574` | Informasi sekunder hanya jika tetap terbaca pada latar aktual. |

Jangan mengasumsikan setiap terminal memakai warna latar atau foreground tertentu. Gunakan warna foreground/background terminal ketika tema `auto` memintanya. Jika warna atau kontras yang sesuai tidak tersedia, gunakan reverse video, bold, underline, indentasi, atau penanda teks yang didukung.

## Tipografi dan hierarki dokumen

- Gunakan font terminal pengguna; jangan menggambar judul figlet.
- Isi dokumen harus paling menonjol. Heading memakai hierarki ukuran/penekanan yang sudah dihasilkan renderer, bukan dekorasi aplikasi yang mengambil banyak ruang.
- Bold dan warna aksen boleh membantu hierarki jika terminal mendukungnya.
- Link dapat memakai underline bila tersedia; status fokus link harus tetap terlihat tanpa warna.
- Code tetap mudah dibedakan melalui gaya teks/struktur, bukan warna saja.

## Struktur dan ornamen

- Pertahankan layout pembaca dokumen yang ringkas. Panel, kartu, tab, tombol, field input, tabel status aplikasi, dan dashboard bukan pola umum untuk layar dokumen.
- Gunakan garis tipis atau pemisah yang sudah relevan dengan konten. Hindari menambah border hanya untuk dekorasi.
- Sudut membulat `╭─╮│╰─╯`, pemisah `·`, dan simbol `◆`/`◇` adalah pilihan aksen, bukan keharusan. Pastikan fallback karakter yang sederhana tersedia bila glyph tidak didukung.
- Jangan menambahkan animasi, spinner, atau transisi untuk tampilan dokumen statis. Perubahan layar harus tetap stabil dan tidak mengganggu membaca.

## Interaksi dan status

- Status fokus, seleksi, bantuan, dan pesan error harus mudah dibedakan serta tidak bergantung pada warna.
- Seleksi teks harus tampak jelas pada tema gelap, terang, transparan, dan tanpa warna. Gunakan gaya terminal yang tersedia; jangan mengubah teks dokumen atau menambahkan karakter dekoratif ke teks yang disalin.
- Pesan status dan petunjuk harus ringkas agar tidak mengambil ruang dari dokumen, khususnya pada terminal kecil.
- Kegagalan mouse atau clipboard harus menampilkan pesan yang jelas tanpa merusak sesi TUI.

## Lebar dan ruang

- Layout harus tetap terbaca mulai 40 kolom dan 10 baris, sesuai spesifikasi produk.
- Hindari padding tetap, panel dengan ruang kosong besar, dan jarak antarelemen yang mengorbankan area baca.
- Perhitungan lebar harus mengikuti lebar tampilan Unicode, bukan jumlah byte atau rune.
- Pada terminal sempit, pertahankan seluruh isi melalui wrapping atau fallback layout yang sudah ditentukan produk.

## Penerapan dan rujukan

Gunakan arahan Ocean sebagai identitas aksen yang lembut, dengan prioritas berikut:

1. `docs/command-spec.md` untuk perilaku command, tema, kontrol, dan batas TUI.
2. `docs/acceptance-criteria.md` untuk syarat lulus dan fallback.
3. Dokumen ini untuk suasana dan pilihan visual yang tidak bertentangan dengan keduanya.

Jika ada konflik, spesifikasi produk dan acceptance criteria berlaku. Perubahan kontrak produk memerlukan keputusan terpisah dan tidak dapat dilakukan hanya berdasarkan arahan desain ini.
