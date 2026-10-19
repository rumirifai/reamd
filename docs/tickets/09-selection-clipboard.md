# 09: Seleksi teks dan copy ke clipboard

**What to build:** Pengguna dapat memilih teks yang terlihat dengan drag dan menyalin teks hasil render ke clipboard tanpa kode warna/control terminal.

**Blocked by:** 08, Tambahkan scroll mouse dan scrollbar (Ticket 01 selesai).

**Status:** blocked

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE` dan `AC-PLATFORM`.

- [ ] Drag memilih teks yang terlihat pada area dokumen. **Acceptance dipindahkan dari checklist Ticket 01.** Bukti kelayakan spike yang tersedia: `.scratch/ticket01/selection-prototype`, Windows/Go 1.27.2, `go test .` lulus (`ok stackspike`) memakai event mouse simulasi; tidak membuktikan drag selection pada terminal nyata.
- [ ] `Ctrl+C` menyalin teks hasil render ke clipboard sistem tanpa ANSI/style control codes; verifikasi pada clipboard OS nyata. **Acceptance dipindahkan dari checklist Ticket 01.** Clipboard tiruan pada test spike bukan bukti keberhasilan clipboard OS.
- [ ] Jika tidak ada pilihan, tampilkan petunjuk dan jangan keluar dari TUI.
- [ ] Jika clipboard tidak tersedia atau gagal, tampilkan pesan yang jelas dan pertahankan pilihan. **Dipindahkan dari checklist Ticket 01:** fallback clipboard produk belum diuji.
- [ ] Seleksi dan copy tidak mengubah dokumen sumber atau menyalin escape sequence aktif; uji plain text tanpa ANSI pada sistem/terminal nyata.

**Verifikasi:** jalankan smoke test clipboard pada sistem desktop yang tersedia dan catat OS, terminal, serta hasilnya; jangan klaim dukungan lintas OS yang belum diuji. Test simulasi di `.scratch/ticket01/selection-prototype` hanya bukti kelayakan spike, tidak memenuhi acceptance produk.
