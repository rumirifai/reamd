# 10: Verifikasi lintas platform

**What to build:** Bukti bahwa seluruh kriteria MVP yang telah diimplementasikan lulus di Windows, macOS, dan Linux, dengan keterbatasan terminal yang dicatat secara jujur.

**Blocked by:** 04, 05, 06, 07, 08, dan 09 selesai.

**Status:** done

**Acuan:** Fase 6 di `docs/development-workflow.md`; seluruh `docs/acceptance-criteria.md`, khususnya `AC-PLATFORM`.

**Prasyarat:** Ticket 04, 05, 06, 07, 08, dan 09 berstatus `done` saat Ticket 10 dimulai. Bukti smoke test Ticket 06 dan 07 yang dilaporkan pengguna dicatat di masing-masing ticket.

- [x] Jalankan test otomatis parsing, layout, input file, sanitasi, help/version, status keluar, dan output non-TTY pada Windows, macOS, serta Linux.
- [x] Jalankan build/test CI pada ketiga OS dan catat hasil yang benar-benar diperoleh.
- [x] Jalankan test terminal/pseudo-terminal untuk keyboard, resize, exit, dan pemulihan terminal jika didukung stack terpilih.
- [x] Smoke test Windows Terminal berdasarkan laporan smoke manual pengguna pada Ticket 06–09.
- [ ] Smoke test Terminal.app pada macOS. **Belum diuji; pemilik proyek menerima ini sebagai blocker rilis yang terdokumentasi dan tidak menghalangi Ticket 11 visual polish. Ini bukan hasil lulus. (23 Oktober 2026)**
- [x] Smoke test GNOME Terminal pada Ubuntu WSL 2/WSLg berdasarkan laporan manual pengguna 23 Oktober 2026.
- [x] Periksa tema, zoom, keyboard, resize, dan keluar normal pada Windows Terminal berdasarkan laporan Ticket 06.
- [x] Periksa anchor, fokus link, pembukaan browser setelah aktivasi, dan kelanjutan TUI setelah browser ditutup di Windows Terminal berdasarkan laporan Ticket 07.
- [x] Periksa seleksi teks dan clipboard OS plain text di Windows Terminal berdasarkan laporan Ticket 09.
- [x] Periksa wheel, scrollbar, dan fallback keyboard berdasarkan laporan smoke manual Ticket 08: Windows 10, Windows Terminal, profile Windows PowerShell; versi terminal tidak dicatat.
- [x] Periksa tema/zoom, keyboard, scroll, resize, keluar normal, anchor/fokus `Tab`/`Enter`, dan copy/paste di GNOME Terminal berdasarkan laporan manual pengguna.
- [x] Uji kondisi opener browser gagal dan clipboard tidak tersedia di GNOME Terminal berdasarkan laporan pengguna: hasil dilaporkan sesuai prosedur; teks pesan persis tidak direkam.
- [x] Konfirmasi terminal kembali normal dan shell tetap dapat digunakan setelah interupsi SIGTERM di GNOME Terminal berdasarkan laporan pengguna 23 Oktober 2026.
- [x] Catat OS, terminal, hasil PASS/FAIL, command test, dan keterbatasan yang belum terselesaikan.
- [x] Pastikan seluruh kriteria wajib lulus atau ada blocker rilis yang disetujui dan terdokumentasi. Smoke test Terminal.app/macOS yang belum dilakukan diterima pemilik proyek sebagai blocker rilis yang tidak menghalangi Ticket 11; lihat catatan pada poin smoke test macOS.
- [x] Verifikasi output teks hasil render melalui terminal/pseudo-terminal: byte ESC/OSC/karakter kontrol dari fixture tidak menjadi kontrol terminal aktif. Pengguna menguji fixture berisi byte kontrol aktual melalui GNOME Terminal/PTTY; tampilan memperlihatkan notasi aman, rekaman menunjukkan payload OSC 52 mentah tidak ada (`False`) dan notasi aman ada (`True`). **Dipindahkan dari checklist Ticket 01:** probe sanitizer scratch tidak melewati renderer aplikasi.
- [x] Smoke test keyboard, resize runtime, dan keluar normal pada GNOME Terminal berdasarkan laporan pengguna; Windows keyboard/resize/keluar normal tercatat di Ticket 06 dan diuji ulang lewat PTY.
- [x] Uji pemulihan terminal nyata setelah interupsi SIGTERM di GNOME Terminal berdasarkan laporan pengguna: proses dihentikan dan shell kembali normal. Pengujian manual Terminal.app/macOS tidak dilakukan dan sudah diterima sebagai blocker rilis pada poin smoke test macOS di atas.
- [x] Smoke test drag selection aplikasi dan clipboard OS plain text tanpa ANSI pada Windows Terminal berdasarkan Ticket 09.
- [x] Smoke test fallback clipboard saat utilitas clipboard tidak tersedia di GNOME Terminal berdasarkan laporan pengguna; detail pesan persis tidak direkam.
- [x] Smoke test integrasi browser OS setelah aktivasi pengguna pada Windows Terminal berdasarkan Ticket 07.
- [x] Verifikasi manual kondisi pesan error TUI saat opener gagal di GNOME Terminal berdasarkan laporan pengguna; detail pesan persis tidak direkam.
- [x] Setelah package clipboard produksi dipilih, periksa ulang graph dependency, lisensi, versi Go, dan advisory graph final. **Dipindahkan dari checklist Ticket 01:** graph terdahulu tidak mencakup package clipboard yang belum dipilih.
- [ ] Tindak lanjut non-blocking: tinjau maintainer/isu satu per satu untuk dependency terpilih; eskalasi menjadi blocker hanya bila ditemukan temuan spesifik yang berpotensi mengubah pilihan dependency. **Asal:** checklist Ticket 01; belum dilakukan, bukan klaim lulus.

## Bukti verifikasi 22–23 Oktober 2026

### Otomatis dan CI

- Mesin lokal: Windows amd64, Go `1.27.2`; terminal proses berjalan di Windows Terminal (`WT_SESSION` tersedia). `gofmt -w cmd internal` lulus; `go test ./...`, `go test -count=1 ./...`, `go vet ./...`, `go build ./...`, `git diff --check`, `go mod verify`, dan `go list -m all` lulus. Test suite mencakup parser/layout, loader file dan path Unicode/spasi, sanitasi byte kontrol aktual, help/version/status keluar, output pipe tanpa ANSI, event TUI simulasi, resize, serta pemulihan layar.
- CI GitHub Actions pada commit `ed52c9b6e8b020fa6cbc06bb03ab924acfa7eb1c` (sama dengan `HEAD` saat pemeriksaan), run [39709363263](https://github.com/rumirifai/reamd/actions/runs/39709363263), selesai 22 Oktober 2026 02:05 UTC: Windows `success` ([job 119596742408](https://github.com/rumirifai/reamd/actions/runs/39709363263/job/119596742408)); macOS `success` ([job 119596742182](https://github.com/rumirifai/reamd/actions/runs/39709363263/job/119596742182)); Linux/Ubuntu `success` ([job 119596742380](https://github.com/rumirifai/reamd/actions/runs/39709363263/job/119596742380)). Build, test, dan vet lulus di semua job; format lulus di macOS/Linux dan memang dilewati di Windows sesuai workflow.
- PTY Windows nyata melalui `pywinpty` menjalankan binary hasil build pada 40×10, mengirim `j` dan `?`, mengubah ukuran ke 60×20, lalu mengirim `q`. Proses keluar dengan status 0; byte pemulihan layar terminal teramati. Batas bukti: tangkapan PTY tidak menunjukkan teks dokumen yang dirender, sehingga percobaan ini bukan bukti bahwa notasi sanitasi terlihat benar di PTY dan bukan smoke test manual GUI.
- Test aplikasi yang memuat byte ESC, BEL, OSC, C1, dan DEL aktual lulus untuk output pipe dan screen tcell simulasi (`go test -run 'TestRunInitializedHandlesKeysAndRestoresScreen|TestRunInitializedRestoresOnInterrupt|TestRunInitializedRestoresAfterResize|TestWithRestoreFinishesOnNormalExit|TestWithRestoreFinishesOnError|TestWithRestoreFinishesBeforeRepanic|TestDrawSanitizesActualTerminalControlBytes|TestRunPipeSanitizesControlCharacters' ./internal/...`). Bukti simulasi/pipe tidak menggantikan inspeksi output render pada terminal nyata.
- Dependency: `go list -m all` menunjukkan graph 18 modul termasuk modul utama; `go mod verify` melaporkan semua modul terverifikasi. Versi minimum di `go.mod` adalah Go `1.25.0`; versi `go` seluruh 17 modul eksternal tidak lebih tinggi. Berkas lisensi seluruh modul eksternal diperiksa: MIT, Apache-2.0, atau BSD-style; tidak ditemukan lisensi tanpa teks lisensi sumber. `go.sum` tidak diubah oleh pekerjaan ini.
- Advisory: `govulncheck` dengan Go `1.27.2` awalnya melaporkan `GO-2026-6809` pada standard library Go, diperbaiki di `1.27.3`, tanpa temuan pada kode/package yang dipakai. Pemindaian ulang `GOTOOLCHAIN=go1.27.3 go run golang.org/x/vuln/cmd/govulncheck@latest ./...` melaporkan tidak ada vulnerability. Toolchain `1.27.3` hanya diunduh/dipakai untuk pemindaian, tidak mengubah `go.mod` atau source.

### Smoke test manual dan batas lingkungan

- Windows 10 / Windows Terminal: laporan dan screenshot manual pengguna pada Ticket 06–09 dihitung sebagai smoke test nyata. Ticket 06 mencatat tema `auto`, `dark`, `light`, `transparent`, `NO_COLOR`, zoom `=`, `+`, `-`, `0`, resize tabel, dan pemulihan prompt setelah `q`. Ticket 07 mencatat fokus link dengan `Tab`, aktivasi yang membuka Chrome ke `https://example.com/docs`, TUI tetap berjalan setelah browser ditutup, serta klik link. Ticket 08 mencatat wheel, scrollbar, dan fallback keyboard pada profile Windows PowerShell. Ticket 09 mencatat drag/copy clipboard plain text untuk Unicode dan kontrol aktual; byte ESC/BEL tidak ada pada hasil `Get-Clipboard -Raw`, dan hash file tidak berubah. Profile Windows PowerShell terlihat pada screenshot; versi Windows Terminal tidak dicatat. Pengujian ini diterima sebagai laporan pengguna, bukan diulang manual pada sesi ini.
- Linux: smoke test manual pengguna pada 23 Oktober 2026 di GNOME Terminal. Identitas lingkungan dari screenshot: Ubuntu 24.04.1 LTS, kernel `6.8.14.1-microsoft-standard-WSL2` x86_64, GNOME Terminal `3.58.0` menggunakan VTE `0.84.0`, melalui WSLg di host Windows 10. Pengguna melaporkan pemeriksaan utama TUI aman, termasuk keyboard, scrolling, resize, copy/paste, `Tab`/`Enter`, opener gagal, clipboard tidak tersedia, serta interupsi `kill -TERM`. Pengguna mengonfirmasi terminal pulih dan shell dapat digunakan setelah SIGTERM; screenshot menunjukkan perintah shell berhasil setelah interupsi. Untuk fixture kontrol aktual, screenshot menunjukkan byte kontrol dirender sebagai notasi aman; hasil pemeriksaan rekaman PTY: payload OSC 52 mentah `False`, notasi aman `True`. Teks pesan error persis dari opener/clipboard tidak direkam.
- Terminal.app/macOS belum diuji. Pemilik proyek menerima kekurangan smoke test manual ini sebagai blocker rilis yang tidak menghalangi Ticket 11. CI macOS tidak dihitung sebagai smoke test desktop dan item Terminal.app tetap unchecked.
- Fitur tema/zoom/resize/keluar normal, aktivasi link, mouse/scrollbar, seleksi/copy, opener gagal, fallback clipboard, dan pemulihan shell setelah SIGTERM kini memiliki laporan manual pada Windows dan/atau GNOME di Ubuntu WSL 2. Pemulihan setelah error fatal/panic tidak diuji manual; test otomatis tetap lulus.
- Verifikasi fixture kontrol aktual pada GNOME Terminal/PTTY melengkapi test output pipe dan screen simulasi. Bukti yang dilaporkan menunjukkan notasi aman terlihat dan payload OSC 52 mentah tidak ada dalam rekaman.

**Status saat ini:** `done` untuk ruang lingkup verifikasi Ticket 10. CI lintas OS, test otomatis, PTY lifecycle Windows, audit dependency, smoke test manual Windows Terminal dan GNOME Terminal di Ubuntu WSL 2, pemulihan setelah SIGTERM, serta sanitasi byte kontrol aktual melalui render terminal/PTTY memiliki bukti yang dicatat di atas. Smoke test Terminal.app/macOS belum dilakukan dan tidak dinyatakan lulus; pemilik proyek menerima kekurangan ini sebagai blocker rilis yang terdokumentasi dan menyetujui bahwa hal itu tidak menghalangi Ticket 11. Tindak lanjut non-blocking maintainer/dependency masih terbuka dan tidak mengubah hasil verifikasi.
