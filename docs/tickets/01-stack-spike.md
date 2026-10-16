# 01: Evaluasi stack dan risiko teknis

**What to build:** Hasil evaluasi yang membuktikan pilihan parser, framework TUI, clipboard, dan pembukaan browser dapat memenuhi kebutuhan MVP, sebelum dependency produksi atau aplikasi utama dibuat.

**Blocked by:** Tidak ada.

**Status:** ready-for-agent

**Acuan:** Fase 1 di `docs/development-workflow.md`; `AC-MD`, `AC-TUI`, `AC-MOUSE`, `AC-SEC`, dan `AC-PLATFORM`.

- [ ] Bandingkan kandidat parser untuk seluruh CommonMark + GFM, conformance corpus, raw HTML aman, dan fallback yang tidak membuang konten.
- [ ] Bandingkan kandidat TUI untuk input keyboard, resize, mouse, pemulihan terminal, serta Windows/macOS/Linux.
- [ ] Uji risiko drag selection, copy clipboard tanpa ANSI, dan fallback saat clipboard tidak tersedia.
- [ ] Uji pembukaan hanya URL HTTP/HTTPS setelah aktivasi pengguna; jangan membangun command shell dari URL.
- [ ] Tinjau pemeliharaan, lisensi, dependency transitif, dan persyaratan versi Go.
- [ ] Catat kandidat, bukti percobaan, trade-off, risiko yang tersisa, dan rekomendasi stack dalam Bahasa Indonesia.
- [ ] Tunggu persetujuan pengguna sebelum menambahkan dependency produksi atau memulai implementasi utama.

**Verifikasi:** hasil spike menunjukkan kriteria penting bisa dipenuhi atau mencatat keputusan produk yang perlu diminta. Jangan menyatakan kemampuan yang belum diuji sebagai terverifikasi.
