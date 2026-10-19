# 01: Evaluasi stack dan risiko teknis

**Tujuan:** Mengumpulkan bukti kelayakan risiko yang berpotensi mengubah pilihan stack sebelum dependency produksi atau aplikasi utama dibuat. Spike bukan acceptance test aplikasi.

**Blocked by:** Tidak ada.

**Status:** done

**Acuan:** Fase 1 di `docs/development-workflow.md`; risiko stack terkait `AC-MD`, `AC-TUI`, `AC-MOUSE`, `AC-SEC`, dan `AC-PLATFORM`.

## Bukti spike dan keputusan stack

- [x] Bandingkan kandidat parser dan uji corpus yang relevan.
  - [x] Goldmark dan gomarkdown cocok dengan CommonMark 0.31.2: 652/652 per parser dalam mode pembanding yang dicatat.
  - [x] Fixture GFM Goldmark mengenali table, task list, strikethrough, dan autolink; ekstensi gomarkdown yang dicoba tidak mendapati task-list semantics.
  - [x] GFM spec 0.29: 648 contoh diparse; 14 exact-output mismatch pada bundle `extension.GFM`, 0 semantic mismatch sesudah normalisasi terbatas yang didokumentasikan. Linkify URL/email literal disetujui pengguna sebagai ekstensi produk.
  - [x] Probe fallback AST mempertahankan raw HTML, fenced code dengan bahasa tak dikenal, dan source Mermaid sebagai teks.
  - [x] Probe mengubah ESC/BEL menjadi notasi terlihat; ini bukan uji renderer terminal.
- [x] Bandingkan API dan risiko stack kandidat TUI.
  - [x] Simulasi tcell menerima key/mouse injection dan dimensi layar simulasi berubah; Bubble Tea dan tcell cross-compile ke Windows/macOS/Linux amd64 pada spike terdahulu.
  - [x] Bukti versi API, lifecycle yang didokumentasikan, dan trade-off dicatat di evaluasi stack.
- [x] Lengkapi satu uji kelayakan terisolasi bahwa event mouse tcell dapat menghasilkan pilihan teks dan plain text melalui abstraksi clipboard.
  - [x] Jalankan test prototipe di `.scratch/ticket01/selection-prototype` pada Windows, Go 1.27.2: `go test .` — lulus (`ok stackspike`).
  - [x] Batas bukti: test memakai event mouse simulasi dan clipboard tiruan; ini hanya bukti kelayakan spike, bukan bukti drag selection di terminal nyata atau keberhasilan clipboard OS.
- [x] Periksa kelayakan helper browser tanpa menjalankan browser OS.
  - [x] Source helper menggunakan executable/argumen terpisah atau ShellExecute, bukan command shell dari URL.
  - [x] Probe menerima URL absolut HTTP(S), menolak skema lain/malformed, dan opener palsu hanya dipanggil setelah aktivasi; error diteruskan.
- [x] Tinjau metadata versi, lisensi, graph, dan advisory yang berisiko mengubah pilihan.
  - [x] Goldmark + tcell + helper browser: 18 module termasuk root; file lisensi ditemukan untuk 17 dependency.
  - [x] Jenis lisensi yang ditemukan: MIT, Apache-2.0, BSD-2-Clause, dan BSD-3-Clause.
  - [x] Module scan baseline menemukan 3 advisory pada x/text v0.31 dan x/sys v0.38; package scan prototipe terdahulu tidak menemukan jalur rentan.
  - [x] Pasangan patch minimum x/text v0.41 + x/sys v0.44 lulus test/vet/build/verify dan module scan dengan Go 1.25; pengguna menyetujui minimum Go 1.25 serta pasangan patch tersebut pada 19 Oktober 2026.
- [x] Catat risiko pemeliharaan yang tersedia, versi, lisensi, graph, advisory, dan pasangan patch minimum yang disetujui. Tinjauan maintainer/isu satu per satu bukan gerbang Fase 1 kecuali muncul temuan spesifik yang dapat mengubah pilihan dependency; tindak lanjut non-blocking dicatat pada Ticket 10.
- [x] Catat kandidat, bukti, trade-off, serta rekomendasi stack dalam Bahasa Indonesia.
- [x] Tunggu persetujuan pengguna sebelum menambah dependency produksi atau memulai implementasi utama.

## Acceptance implementasi/verifikasi yang dialihkan

Checklist berikut dipindahkan ke ticket yang mengimplementasikan atau memverifikasi perilaku. Item belum terbukti tetap belum lulus; pencatatan ini bukan bukti hasilnya.

- [ ] **Renderer aplikasi, fallback semua node, dan keluaran teks terminal tanpa kontrol aktif:** dipindahkan dari Ticket 01 ke Ticket 03 (fallback/render seluruh konten) dan Ticket 05 (sanitasi output aplikasi); verifikasi output terminal/pseudo-terminal di Ticket 10. Probe fallback dan sanitasi scratch bukan uji renderer aplikasi.
- [ ] **Binding keyboard, resize runtime, dan pemulihan terminal nyata saat keluar/error/interupsi:** dipindahkan dari Ticket 01 ke Ticket 02 untuk vertical slice dan Ticket 10 untuk verifikasi terminal/pseudo-terminal serta lintas OS. API dan simulasi bukan pengganti pengujian runtime.
- [ ] **Smoke test TUI lintas OS/terminal:** dipindahkan dari Ticket 01 ke Ticket 02 (platform tersedia untuk vertical slice) dan Ticket 10 (matriks sasaran); belum dijalankan.
- [ ] **Drag selection aplikasi, clipboard OS plain text tanpa ANSI, dan fallback clipboard:** dipindahkan dari Ticket 01 ke Ticket 09 untuk implementasi/acceptance dan Ticket 10 untuk smoke OS. Test simulasi di atas hanya bukti spike, tidak meluluskan acceptance produk.
- [ ] **Integrasi browser OS, pesan error TUI, jalur parse/render tanpa opener otomatis, dan smoke terminal:** dipindahkan dari Ticket 01 ke Ticket 07 (gate aktivasi, jalur parse/render, pesan kegagalan) dan Ticket 10 (browser OS nyata/lintas platform). Validator URL dan opener tiruan yang lulus tetap bukti spike, bukan acceptance integrasi.
- [ ] **Graph dependency setelah package clipboard produksi dipilih:** dipindahkan dari Ticket 01 ke Ticket 09 dan Ticket 10; evaluasi graph terpilih belum dilakukan karena package produksi belum dipilih.
- [ ] **Peninjauan maintainer/isu satu per satu:** tindak lanjut non-blocking pada Ticket 10, kecuali ditemukan isu spesifik yang berpotensi mengubah pilihan dependency. Statusnya belum dilakukan dan tidak menghalangi penutupan Fase 1 berdasarkan bukti yang tersedia.

**Bukti dokumentasi:** [evaluasi stack dan risiko teknis](../stack-evaluation.md) memuat corpus CommonMark/GFM, probe fallback/sanitasi, simulasi tcell, test prototipe seleksi, browser gate, graph, lisensi, dan advisory beserta batas klaimnya. Stack Goldmark + `tcell/v2`, helper browser, Linkify, minimum Go 1.25, serta pasangan patch minimum telah disetujui. Ticket 01 ditutup karena seluruh pekerjaan Fase 1 yang dapat mengubah rekomendasi memiliki bukti dan persetujuan; Ticket 02 dibuka sebagai ready-for-agent tetapi tidak dikerjakan dalam sesi ini. Root repo tetap tanpa `go.mod` atau dependency produksi.

**Gerbang penutupan terpenuhi:** risiko stack yang dapat mengubah keputusan memiliki bukti yang dapat diperiksa dan persetujuan yang diperlukan sudah diberikan. Acceptance aplikasi/runtime/OS di atas tetap terbuka pada ticket pemiliknya dan tidak dinyatakan lulus oleh spike.
