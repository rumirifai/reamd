# 02: Bangun CLI dan vertical slice TUI

**What to build:** Pengguna dapat menjalankan `reamd <path>` untuk membuka satu file di terminal interaktif, membaca dan menggulir tampilan dasar, keluar kembali ke shell, atau menerima teks biasa saat output dialihkan.

**Blocked by:** Tidak ada; Ticket 01 selesai.

**Status:** done

**Acuan:** Fase 2-3 di `docs/development-workflow.md`; `AC-CLI` dan bagian dasar `AC-TUI`/`AC-THEME`.

- [x] Buat fondasi Go sesuai keputusan stack yang disetujui dan jaga pemisahan CLI, pemuatan file, model dokumen, layout, dan TUI.
- [x] Dukung satu path file, `--help`/`-h`, `--version`/`-v`, dan `--theme` dengan validasi sesuai `docs/command-spec.md`.
- [x] Saat tanpa argumen, tampilkan bantuan dan contoh dalam Bahasa Inggris lalu keluar dengan kode 0.
- [x] Validasi jumlah argumen, opsi tema, file tidak ditemukan, file tidak dapat dibaca, input non-UTF-8/biner, pesan stderr, dan kode keluar 1 atau 2 sesuai jenis error.
- [x] Di terminal interaktif, tampilkan konten dasar dalam TUI yang dapat digulir dengan keyboard; `q`, `Esc`, dan `?` berfungsi. **Dipindahkan dari checklist Ticket 01:** uji keyboard adalah acceptance aplikasi; API/simulasi spike bukan bukti interaktif.
- [x] Tangani ukuran minimum 40×10 dan resize; pulihkan terminal setelah keluar normal, error, panic, dan interupsi. **Dipindahkan dari checklist Ticket 01:** resize runtime dan pemulihan terminal nyata milik vertical slice, bukan gerbang pemilihan stack.
- [x] Saat output dialihkan atau dipipe, jangan ambil alih layar TUI dan jangan keluarkan ANSI secara default.
- [x] Verifikasi path dengan spasi/Unicode serta pastikan file tidak berubah setelah dibuka.
- [x] Jalankan smoke test vertical slice pada platform yang tersedia dan catat platform/terminal serta batasannya; matriks seluruh OS/terminal tetap milik Ticket 10. **Dipindahkan dari checklist Ticket 01:** smoke aplikasi bukan bukti spike.

## Bukti penyelesaian

- Otomatis: test mencakup CLI/status keluar, file valid dan invalid, output pipe, path Unicode/spasi, file tetap utuh, navigasi/resize, ukuran minimum, dan lifecycle restore melalui layar simulasi.
- Smoke interaktif: Windows 11/amd64, pseudo-terminal Windows ConPTY melalui `pywinpty`, ukuran awal 40×10. File 200 baris berhasil dibuka; `j`, PageDown, End, resize ke 60×20 dan 30×8, pesan terminal kecil, lalu Esc dan exit code 0 teramati. Smoke tambahan dengan file path Unicode/spasi menunjukkan konten tampil, `?` membuka bantuan, `q` keluar dengan exit code 0, dan urutan restore layar terlihat. Ini bukan smoke di Windows Terminal GUI/macOS/Linux.
- Batas verifikasi: permission-denied nyata dan sinyal OS nyata tidak disimulasikan di PTY; jalur loader permission dan pemulihan event interruption/panic diuji otomatis dengan batasan platform yang dicatat dalam test. CI lintas OS telah ditambahkan, tetapi workflow belum dijalankan oleh layanan CI.

**Verifikasi:** kriteria Ticket 02 lulus pada platform yang tersedia; keterbatasan lintas OS/terminal dicatat di atas dan menjadi cakupan Ticket 10.
