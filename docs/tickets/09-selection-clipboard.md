# 09: Seleksi teks dan copy ke clipboard

**What to build:** Pengguna dapat memilih teks yang terlihat dengan drag dan menyalin teks hasil render ke clipboard tanpa kode warna/control terminal.

**Blocked by:** 01, Evaluasi stack dan risiko teknis; 08, Tambahkan scroll mouse dan scrollbar.

**Status:** blocked

**Acuan:** Fase 5 di `docs/development-workflow.md`; `AC-MOUSE` dan `AC-PLATFORM`.

- [ ] Drag memilih teks yang terlihat pada area dokumen.
- [ ] `Ctrl+C` menyalin teks hasil render ke clipboard sistem tanpa ANSI/style control codes.
- [ ] Jika tidak ada pilihan, tampilkan petunjuk dan jangan keluar dari TUI.
- [ ] Jika clipboard tidak tersedia atau gagal, tampilkan pesan yang jelas dan pertahankan pilihan.
- [ ] Seleksi dan copy tidak mengubah dokumen sumber atau menyalin escape sequence aktif.

**Verifikasi:** jalankan smoke test clipboard pada sistem desktop yang tersedia dan catat OS, terminal, serta hasilnya; jangan klaim dukungan lintas OS yang belum diuji.
