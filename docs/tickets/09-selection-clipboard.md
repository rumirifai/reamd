# 09: Seleksi teks dan copy ke clipboard

**What to build:** Pengguna dapat memilih teks yang terlihat dengan drag dan menyalin teks hasil render ke clipboard tanpa kode warna/control terminal.

**Blocked by:** 08, Tambahkan scroll mouse dan scrollbar (Ticket 01 selesai).

**Status:** done

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE` dan `AC-PLATFORM`.

- [x] Drag memilih teks yang terlihat pada area dokumen; test otomatis mencakup pemetaan area dokumen, Unicode, dan seleksi lintas baris ter-wrap. **Acceptance dipindahkan dari checklist Ticket 01.** Bukti spike `.scratch/ticket01/selection-prototype` bukan bukti produk atau terminal nyata.
- [x] `Ctrl+C` menyalin teks hasil render ke clipboard sistem tanpa ANSI/style control codes; smoke test clipboard OS nyata dicatat di bawah. **Acceptance dipindahkan dari checklist Ticket 01.** Clipboard tiruan pada test spike bukan bukti keberhasilan clipboard OS.
- [x] Jika tidak ada pilihan, tampilkan petunjuk dan jangan keluar dari TUI.
- [x] Jika clipboard tidak tersedia atau gagal, tampilkan pesan yang jelas dan pertahankan pilihan; test otomatis memakai adapter tiruan.
- [x] Seleksi dan copy tidak mengubah dokumen sumber atau menyalin escape sequence aktif; smoke test plain text tanpa ANSI pada sistem/terminal nyata dicatat di bawah.

**Keputusan clipboard:** pengguna menyetujui `github.com/atotto/clipboard` di balik adapter internal `internal/clipboard`; persetujuan dicatat di `docs/markdown-terminal-reader-analysis.md`.

**Bukti otomatis:** test TUI mencakup drag dan pemetaan area dokumen, Unicode, baris ter-wrap, pilihan kosong, tampilan seleksi, serta copy plain text melalui adapter tiruan untuk sukses/tidak tersedia/gagal. Fixture karakter kontrol aktual melewati jalur `document.New` dan layout; clipboard tiruan tidak menerima ANSI/escape aktif, dan sumber tetap tidak berubah.

**Smoke test manual (berdasarkan laporan dan tangkapan layar pengguna):** OS Windows, terminal Windows Terminal (versi tidak tercatat), mekanisme `github.com/atotto/clipboard` melalui adapter internal. Pengguna menjalankan aplikasi pada fixture Unicode dan karakter kontrol aktual, menyalin pilihan ke clipboard, lalu memeriksa hasil dengan `Get-Clipboard -Raw`. Teks clipboard menampilkan `café`, `中文`, `🙂`, serta notasi aman untuk ESC (`^[[2J`) dan OSC/BEL (`^[]52;c;TEST^G`); pemeriksaan byte ESC dan BEL mentah menghasilkan `False`. Hash SHA-256 file sebelum dan sesudah sesi sama (`True`). Bukti ini hanya untuk Windows Terminal yang diuji dan tidak menyatakan dukungan lintas OS.
