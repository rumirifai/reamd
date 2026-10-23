# 11: Poles visual TUI sesuai arahan Ocean

**What to build:** Rapikan tampilan TUI yang sudah berfungsi agar mendukung pembacaan Markdown yang tenang, jelas, dan nyaman, dengan menerapkan `docs/DESIGN.md`.

**Blocked by:** 10, Verifikasi lintas platform.

**Status:** ready-for-agent

**Acuan:** `docs/DESIGN.md`; Fase 7 di `docs/development-workflow.md`; `AC-RENDER`, `AC-THEME`, `AC-TUI`, `AC-MOUSE`, dan `AC-PLATFORM`.

**Prasyarat:** Ticket 10 berstatus `done`. Ticket 10 mencatat Terminal.app/macOS sebagai blocker rilis yang diterima dan tidak menghalangi Ticket 11; keterbatasan itu tetap terbuka dan tidak boleh disebut lulus.

- [ ] Audit tampilan TUI aktual terhadap arahan suasana, kepadatan, warna, hierarki Markdown, status interaksi, dan lebar terminal dalam `docs/DESIGN.md`; catat masalah visual spesifik sebelum mengubah kode.
- [ ] Perbaiki hirarki heading, paragraf/list/quote, inline dan fenced code, link/fokus, seleksi, scrollbar, bantuan, serta pesan status/error hanya pada bagian yang terbukti kurang jelas atau tidak konsisten dengan arah desain.
- [ ] Terapkan aksen Ocean secara hemat dan sesuai tema; jangan memaksakan latar gelap pada `auto` atau `transparent`, dan jangan bergantung pada warna untuk membedakan informasi penting.
- [ ] Pertahankan layout ringkas yang memprioritaskan isi dokumen. Jangan menambah dashboard, kartu/panel dekoratif, judul figlet, animasi, atau elemen yang mengurangi ruang baca.
- [ ] Pastikan tampilan tetap terbaca pada 40×10 dan lebar 40, 80, serta 120 kolom; gunakan lebar tampilan Unicode yang sudah dipakai layout.
- [ ] Pastikan tema `auto`, `dark`, `light`, `transparent`, `NO_COLOR`, dan terminal tanpa warna/warna terbatas memiliki fallback yang masuk akal sesuai spesifikasi.
- [ ] Tambahkan atau perbarui test untuk perubahan style/state yang dibuat. Hindari snapshot ANSI yang mewajibkan warna identik lintas terminal; uji struktur, atribut/style yang relevan, dan fallback.
- [ ] Ulangi test otomatis yang relevan dan catat pemeriksaan visual/smoke test yang benar-benar dilakukan. Perubahan ini terjadi setelah Ticket 10, sehingga hasil sebelum polish tidak otomatis membuktikan tampilan akhir.
- [ ] Perbarui README hanya jika klaim status atau penggunaan yang terlihat pengguna berubah; README tetap berbahasa Inggris saja.

**Verifikasi:** test relevan lulus; tinjauan TUI pada ukuran dan tema representatif dicatat; smoke test manual hanya ditandai lulus bila benar-benar dijalankan. Jika OS/terminal tidak tersedia, catat sebagai belum diuji dan jangan mengklaim cakupan lintas platform.

## Batas pekerjaan

- Ticket ini adalah polish visual, bukan penambahan fitur produk atau perubahan kontrak interaksi.
- `docs/command-spec.md` dan `docs/acceptance-criteria.md` tetap menjadi sumber kebenaran perilaku dan syarat lulus; `docs/DESIGN.md` memberi arahan visual selama tidak berkonflik.
- Jika implementasi mengungkap konflik desain dengan spesifikasi produk, hentikan bagian yang bergantung pada keputusan itu dan tanyakan kepada pengguna.
- Jangan mulai Ticket 12 atau mengubah metode distribusi.

## Bukti penyelesaian

Agent mencatat masalah visual yang ditemukan, pilihan perbaikannya, file yang berubah, test/command dan hasilnya, tema/ukuran terminal yang diperiksa, serta smoke test manual dan keterbatasan lingkungan. Tandai `done` hanya berdasarkan bukti tersebut.
