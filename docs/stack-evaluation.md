# Evaluasi stack dan risiko teknis — Ticket 01

Tanggal evaluasi: 16 Oktober 2026; administrasi penutupan Ticket 01 diperbarui 19 Oktober 2026. Tujuan: mengumpulkan bukti sebelum menambah dependency produksi atau memulai implementasi. Semua kode eksperimen dan module Go dibuat di scratch/Temp, bukan sebagai kode aplikasi di repo. Tidak ada `go.mod` atau dependency yang ditambahkan ke repo utama.

## Ringkasan keputusan

**Disetujui pengguna:** 16 Oktober 2026. Stack yang disetujui untuk memulai implementasi adalah Goldmark + `tcell/v2`. Minimum Go 1.24 yang disetujui saat itu diperbarui menjadi **Go 1.25 pada 19 Oktober 2026**, setelah audit menemukan advisory pada dependency transitif dan pasangan patch minimum Go 1.25 lulus pemeriksaan scratch. Pengguna juga menyetujui Linkify untuk URL/email polos sebagai ekstensi perilaku produk di atas profil GFM spec 0.29. Persetujuan ini bukan bukti seluruh acceptance Ticket 01 lulus; belum ada dependency produksi atau aplikasi yang ditambahkan.

- **Parser:** pilih `goldmark` sebagai kandidat pertama untuk parser/AST terminal. Paket CommonMark + GFM bawaan mencakup tabel, task list, strikethrough, dan autolink; uji corpus CommonMark lulus 652/652. Renderer HTML default tidak boleh dipakai sebagai renderer aplikasi: ia menghilangkan raw HTML; gunakan AST dan konversi raw-HTML node menjadi teks, lalu sanitasi karakter kontrol secara terpisah. Adapter tersebut belum dibuat atau diuji.
- **TUI:** pilih `tcell/v2` untuk spike implementasi berikutnya dibanding Bubble Tea. Alasannya: API layar berbasis sel dan simulasi yang bisa diinjeksi event-nya cocok untuk layout grid, uji input/resize, serta kontrol mouse. Trade-off: framework lebih rendah tingkatnya, sehingga layout, navigasi, bantuan, dan event loop lebih banyak menjadi tanggung jawab aplikasi. Bubble Tea tetap kandidat valid dengan API model/update/view yang lebih ringkas, tetapi tidak ada bukti bahwa ia mempermudah seleksi teks dan scrollbar yang diwajibkan.
- **Clipboard:** jangan menetapkan implementasi final dahulu. Bungkus di antarmuka kecil yang mengembalikan error dan mempertahankan seleksi. `atotto/clipboard` memiliki API lintas platform sederhana tetapi rilis terakhir yang terlihat berumur lama; Linux memerlukan `xclip` atau `xsel`. Clipboard OSC 52 dari tcell dapat ditolak terminal dan API `SetClipboard` tidak memberikan konfirmasi sukses, sehingga tidak cukup sebagai satu-satunya mekanisme untuk memenuhi pesan kegagalan yang disyaratkan.
- **Browser:** `github.com/pkg/browser` layak menjadi helper kandidat; source yang diperiksa menjalankan proses dengan executable dan argumen terpisah di macOS/Linux, dan memakai ShellExecute di Windows—bukan menyusun command shell. Validasi skema `http`/`https` dan aksi pengguna harus tetap dilakukan aplikasi sebelum pemanggilan. Tidak ada browser yang diluncurkan dalam spike ini.
- **Versi Go:** minimum proyek yang berlaku disetujui **Go 1.25** pada 19 Oktober 2026. Tcell langsung menyatakan `go 1.24.0`; Goldmark dan gomarkdown menyatakan `go 1.22`. Dependency perbaikan minimum `golang.org/x/text v0.41.0` dan `golang.org/x/sys v0.44.0` menyatakan Go 1.25. Go 1.27.2 milik pengguna memenuhi syarat.

Acceptance aplikasi tetap belum dibuktikan pada Fase 1; kriteria tersebut dicatat pada ticket implementasi/verifikasi pemiliknya. Persetujuan stack tidak mengubah kontrak produk di `command-spec.md`, dan tidak membuat `go.mod` atau dependency produksi.

## Kandidat parser

| Kandidat | Bukti/dukungan | Lisensi, Go, dependency | Penilaian |
|---|---|---|---|
| `github.com/yuin/goldmark` `v1.8.6` | Menyatakan CommonMark 0.31.2, AST yang dapat diperluas, dan ekstensi GFM bawaan `GFM` (table, strikethrough, linkify, task list). Uji corpus lokal: 652/652 cocok ketika HTML renderer disetel `WithXHTML` + `WithUnsafe`. Uji fitur GFM lokal melihat table, checkbox task, strikethrough, dan autolink. | MIT; module `go 1.22`; dokumentasi menyatakan hanya standard library sebagai dependency. Published 3 Sep 2026; versi dan lisensi dilihat di pkg.go.dev. | Kandidat utama. Renderer HTML default sengaja menyembunyikan raw HTML sebagai `<!-- raw HTML omitted -->`, sehingga kehilangan isi yang acceptance criteria minta tampil sebagai teks aman. `WithUnsafe` hanya dipakai pada uji kesesuaian; jangan dipakai untuk menulis output terminal. Implementasikan renderer AST sendiri dan uji fallback aman sebelum menganggap AC-MD/AC-SEC lulus. |
| `github.com/gomarkdown/markdown` `v0.0.0-20261011113952-f83953680332` | Mode `parser.CommonMark` menghasilkan keluaran HTML persis untuk 652/652 contoh CommonMark 0.31.2 pada probe. `CommonExtensions` lokal mengenali table, strikethrough, dan autolink, tetapi tidak menghasilkan checkbox/task semantics. AST menyediakan node `HTMLBlock`/`HTMLSpan` dan leaf berisi content/literal untuk renderer khusus. | BSD-2-Clause; module `go 1.22`; dokumentasi menyebut hanya standard library. Published 11 Okt 2026. | Pembanding CommonMark yang kuat, namun paket GFM bawaan yang diuji belum lengkap untuk task list. Menggabungkan mode strict CommonMark dengan ekstensi GFM tidak dibuktikan; perlu parser/ekstensi tambahan atau implementasi sendiri. Bukan rekomendasi pertama untuk seluruh target GFM. |

Korpus yang diunduh adalah [CommonMark 0.31.2 `spec.json`](https://spec.commonmark.org/0.31.2/spec.json). Probe membandingkan HTML sesudah trim spasi tepi. Ketidakcocokan awal Goldmark 72/652 terjadi saat renderer default membuang raw HTML; setelah konfigurasi keluaran spec mentah, hasilnya 0/652. Ini menguji parser + keluaran referensi, **bukan** keamanan renderer terminal. GFM fixture yang diuji hanya contoh sempit; belum ada GFM conformance corpus penuh atau fallback node yang tidak dirender.

## Kandidat TUI

| Kandidat | Bukti kemampuan dan platform | Lisensi, Go, dependency | Trade-off/risiko |
|---|---|---|---|
| `github.com/gdamore/tcell/v2` `v2.13.10` | Dokumentasi API menyatakan key/mouse/resize, terminal Windows dan POSIX, pemulihan lewat `Stop`/`Fini`, Unicode/grapheme, dan `SimulationScreen`. Probe scratch dengan `SimulationScreen` berhasil menerima key dan mouse yang diinjeksi, mengamati ukuran 40×10 setelah perubahan ukuran, serta menyimpan plain text ke clipboard simulasi. Library dan Bubble Tea berhasil dikompilasi untuk `windows/amd64`, `darwin/amd64`, dan `linux/amd64`. | Apache-2.0; module `go 1.24.0`; pure Go menurut upstream; sekitar 29 import langsung pada halaman pkg.go.dev (bukan jumlah graph proyek). Published 6 Mei 2026. | Lebih banyak logika UI/lifecycle harus ditulis sendiri. Simulasi membuktikan API mock dapat diuji, bukan perilaku terminal nyata. Clipboard tcell adalah protokol terminal yang dapat diabaikan/diblokir; `SetClipboard` tidak mengembalikan error. Tidak membuktikan restore saat panic/sinyal maupun drag selection aplikasi. |
| `github.com/charmbracelet/bubbletea` `v1.3.10` | API model/update/view, key, mouse event (termasuk drag jika mouse cell-motion diaktifkan), window-size messages, alternate screen, signal handling, `Run`, serta restore terminal didokumentasikan. Kompilasi silang ke tiga OS/arsitektur di atas lulus. Dokumentasi mencatat Windows tidak mengirim resize via `SIGWINCH`; resize nyata Windows perlu diuji. | MIT; module `go 1.24.0`; pemindaian module scratch membawa banyak module transitif Charmbracelet, terminal, warna, dan Unicode. Published 17 Sep 2025. | Lebih tinggi tingkat abstraksinya dan cocok untuk event loop, tapi tidak menyediakan seleksi teks/scrollbar siap pakai; itu tetap fitur aplikasi. Mengaktifkan mouse motion dapat mengambil alih seleksi native terminal. Tidak ada smoke test interaktif atau uji restore terminal yang dijalankan. |

Kompilasi silang memverifikasi source dan build tags kandidat, tidak memverifikasi binary `reamd`, terminal tertentu, atau runtime di macOS/Linux. Probe tcell tidak menguji restore state fisik. `SimulationScreen.SetSize` mengubah ukuran layar simulasi namun tidak mengantrekan `EventResize` secara otomatis pada skenario yang dicoba; test awal yang menunggu event itu terblokir dan dihentikan. Probe yang diperbaiki memeriksa dimensi baru dan tidak mengklaim telah menguji notifikasi resize runtime.

**Drag selection:** kedua pendekatan mengirim mouse event kepada aplikasi, jadi seleksi yang diminta harus diimplementasikan aplikasi. Mode mouse terminal dapat menelan seleksi native; perlu percobaan manual di Windows Terminal dan terminal sasaran lain, termasuk apakah tombol modifier dapat menjadi fallback. Belum terbukti untuk kedua kandidat.

## Clipboard

- `github.com/atotto/clipboard` `v0.1.4`, BSD-3-Clause, Go module metadata lama/tidak menyediakan directive modern yang menjadi batas aman. pkg.go.dev menunjukkan publish 24 Feb 2021. README menyatakan macOS, Windows, dan Linux/Unix; Linux memerlukan executable sistem `xclip` atau `xsel`. API mengembalikan error, berguna untuk pesan gagal, tetapi kondisi desktop tanpa utilitas harus ditangani.
- tcell memiliki `Screen.SetClipboard([]byte)`/`GetClipboard`, mengirim permintaan clipboard melalui terminal; dokumentasinya menyebut terminal dapat menolak permintaan. Tidak ada nilai return keberhasilan untuk `SetClipboard`.
- Tidak ada penulisan atau pembacaan clipboard OS sungguhan, agar tidak mengubah clipboard pengguna. Uji simulasi tcell hanya menyimpan data di objek test.

Rekomendasi: gunakan interface `WriteText(string) error` agar model seleksi tidak bergantung pada library. Sebelum mengunci package, jalankan manual di sesi Windows/macOS/Linux, termasuk tanpa backend clipboard; pastikan ANSI tidak ikut dan seleksi tetap tersedia saat gagal.

## Membuka browser

`github.com/pkg/browser` versi pseudo `v0.0.0-20240102092130-5ac0b6a4141c` (BSD-2-Clause, `go 1.14`, publish 2 Jan 2024) menyediakan `OpenURL`. Source yang diunduh diperiksa: macOS memanggil `open` lewat `exec.Command`; Linux memilih `xdg-open`, `x-www-browser`, atau `www-browser`, juga lewat `exec.Command`; Windows memanggil `windows.ShellExecute`. URL menjadi argumen/data API, bukan bagian dari string shell. Library hanya membuka URL, bukan memvalidasi protokol atau mendeteksi aksi pengguna.

Tidak ada opener dijalankan. Belum ada pengujian URL adversarial atau interaksi terminal. Implementasi wajib parse URL, hanya mengizinkan skema `http`/`https`, memanggil helper hanya atas aktivasi langsung (klik/Enter), menangani error, serta tidak membuka link saat parse/render dokumen.

## Pemeliharaan, versi, dan graph dependency

Pemeriksaan sumber dilakukan 16 Okt 2026 dari pkg.go.dev/module metadata dan module yang benar-benar diunduh pada scratch. Versi bukan pin produksi dan perlu ditinjau lagi saat dependency akan dipilih.

| Kandidat | Versi yang diperiksa | Lisensi | Go directive module |
|---|---|---|---:|
| Goldmark | `v1.8.6`, terbit 2026-09-03 | MIT | 1.22 |
| gomarkdown | `v0.0.0-20261011113952-f83953680332`, terbit 2026-10-11 | BSD-2-Clause | 1.22 |
| Bubble Tea | `v1.3.10`, terbit 2025-09-17 | MIT | 1.24 |
| tcell | `v2.13.10`, terbit 2026-05-06 | Apache-2.0 | 1.24 |
| atotto/clipboard | `v0.1.4`, terbit 2021-02-24 | BSD-3-Clause | metadata module lama |
| pkg/browser | `v0.0.0-20240102092130-5ac0b6a4141c`, terbit 2024-01-02 | BSD-2-Clause | 1.14 |

`go list -m all` pada scratch gabungan parser, kedua TUI, clipboard, dan browser menghasilkan 39 module (root + dependency graph gabungan; bukan jumlah khusus stack final). Goldmark/gomarkdown menyatakan standard library only; Bubble Tea menarik graph lebih besar, sementara tcell juga membawa terminal/encoding dependencies. Graph perlu dihitung ulang setelah kandidat final dipilih. Tidak dilakukan audit CVE, lisensi dependency transitif satu per satu, atau penilaian kualitas maintainer/response time.

## Percobaan yang dijalankan dan hasil

1. `git --no-optional-locks status --short --branch` sebelum pekerjaan: `main...origin/main`; working tree sudah memiliki perubahan `docs/markdown-terminal-reader-analysis.md` dari sesi sebelumnya. Perubahan itu dipertahankan dan diperbarui, tidak direset.
2. `"C:/Users/Kyouka Ruri/.g/versions/1.27.2/bin/go.exe" version` dan `go env GOVERSION GOPATH GOPROXY`: Go `1.27.2 windows/amd64` terinstal; executable tidak ada di `PATH` shell agent, maka dipanggil dengan path penuh.
3. Di scratch: `go mod init stackspike`; unduh Goldmark, gomarkdown, Bubble Tea, tcell, atotto/clipboard, pkg/browser, serta CommonMark spec JSON. Semua module dan file corpus tetap di Temp.
4. `go test -v` corpus: goldmark 0/652 mismatches dengan `WithXHTML` + `WithUnsafe`; gomarkdown 0/652 dengan `parser.CommonMark` dan renderer default CommonMark. Pengaturan raw untuk goldmark hanya untuk baseline corpus.
5. `go test -v` GFM Goldmark: table, checkbox task, strikethrough, dan autolink ditemukan; raw HTML menjadi komentar “raw HTML omitted” dan konten literal tidak dipertahankan oleh renderer default.
6. `go test -run 'TestGomarkdownCommonExtensions|TestGoldmarkGFMAndSafeHTML' -v`: gomarkdown mengenali table, strikethrough, autolink, tetapi tidak menghasilkan checkbox GFM task list pada konfigurasi extension yang dicoba.
7. `go test -run 'TestTCellSimulationEvents' -v`: lulus untuk key/mouse injection, perubahan ukuran layar simulasi dan plain-text clipboard simulasi. Tidak menguji OS clipboard atau terminal nyata.
8. `GOOS=windows/darwin/linux GOARCH=amd64 go build github.com/charmbracelet/bubbletea github.com/gdamore/tcell/v2`: seluruh kompilasi silang berhasil.
9. `go list -m -json` memeriksa versi, tanggal dan Go directive; `go list -m all | wc -l` menghitung 39 module pada graph scratch gabungan. Source browser helper diperiksa; tidak diluncurkan.
10. Satu probe tcell awal timeout saat menunggu `EventResize` yang tidak otomatis dikirim SimulationScreen pada alur itu. Test diperbaiki; hasil final yang dilaporkan adalah test simulasi no. 7, bukan klaim event resize fisik.

Tidak dijalankan: CommonMark/GFM full conformance untuk seluruh mode GFM, renderer AST terminal yang aman, fallback node tak dikenal, sanitasi ESC/OSC aktual, Bubble Tea interactive test, clipboard nyata, browser launch, terminal restore, atau smoke test Windows Terminal/Terminal.app/GNOME Terminal. Tidak dibuat go.mod/dependency di repo dan tidak ada aplikasi produk.

## Risiko dan acceptance yang belum terbukti

| Kriteria | Status bukti pada spike |
|---|---|
| `AC-MD` | CommonMark corpus untuk parser-pasangan lulus sesuai mode di atas; belum ada GFM corpus penuh, raw HTML fallback teks aman, Mermaid fallback, atau test seluruh node. Belum lulus seluruh AC. |
| `AC-TUI` | API keyboard/resize/mouse terlihat dan tcell simulasi terbatas lulus; tidak ada aplikasi/TUI nyata, resize real, keyboard binding produk, atau restore terminal. Belum lulus. |
| `AC-MOUSE` | tcell menerima event mouse dalam simulasi. Seleksi teks/scrollbar belum ada; clipboard OS dan fallback gagal belum diuji; browser tidak diluncurkan. Belum lulus. |
| `AC-SEC` | Tidak ada pengujian byte ESC/OSC melalui terminal renderer; Goldmark mode unsafe hanya digunakan untuk conformance corpus dan tidak boleh digunakan pada input nyata. Belum lulus. |
| `AC-PLATFORM` | Dua library TUI berhasil cross-compile tiga OS amd64. Tidak ada app build/test CI atau smoke test terminal pada OS mana pun. Belum lulus. |

## Persetujuan dan batas gerbang

Keputusan awal 16 Oktober 2026 adalah minimum Go 1.24, Goldmark, `tcell/v2`, clipboard package final ditunda dengan interface error dan seleksi dipertahankan saat gagal, serta `github.com/pkg/browser` sebagai helper kandidat dengan validasi HTTP(S) dan aktivasi eksplisit. Pada 19 Oktober, pengguna menyetujui minimum Go 1.25 setelah bukti advisory dan menyetujui Linkify untuk URL/email literal sebagai ekstensi produk di atas GFM spec 0.29.

Persetujuan rekomendasi bukan bukti acceptance teknis. Jangan menganggap Ticket 01 selesai, atau menambahkan dependency produksi, sebelum checklist teknis yang tersisa memiliki bukti. Sesi ini tidak memulai Ticket 02.

## Pembaruan lanjutan — 19 Oktober 2026

### Keadaan repo dan lingkungan

- `git --no-optional-locks status --short --branch` menunjukkan `main...origin/main` dengan perubahan lokal yang sudah ada pada `README.md`, `docs/markdown-terminal-reader-analysis.md`, dan `docs/tickets/01-stack-spike.md`, serta `docs/stack-evaluation.md` sebagai file untracked. Perubahan tersebut dipertahankan; tidak ada aplikasi atau modul produksi yang dibuat.
- Pencarian manifest/CI hingga kedalaman 3 tidak menemukan `go.mod`, `go.sum`, `go.work`, atau workflow GitHub Actions di repo.
- `go version` berhasil: `go1.27.2 windows/amd64`; `go env` melaporkan `GOPROXY=https://proxy.golang.org,direct`. Jadi Go tersedia walaupun sebelumnya executable tidak ada di PATH shell agent, dan unduhan dependency dapat dilakukan.
- `go mod download -json github.com/yuin/goldmark@v1.8.6` berhasil. Output menunjuk ke cache global pengguna `C:\Users\Kyouka Ruri\go\pkg\mod`; tidak ada file module yang ditambahkan ke repo. Ini membuktikan akses unduhan untuk satu module pada saat pemeriksaan, bukan jaminan ketersediaan jaringan di lingkungan CI atau untuk semua host dependency.
- Prototipe sebelumnya disalin pengguna ke `.scratch/ticket01/` (folder ini tetap berisi module dan file test terisolasi, bukan kode aplikasi). Test tidak dapat dijalankan langsung dari Temp karena batas working directory tool. Pada percobaan pertama sandbox WSL gagal start; setelah persetujuan untuk menjalankan pemeriksaan terbatas tanpa sandbox, test dapat dijalankan dari scratch proyek. Detail hasil pengulangan ada di bagian berikut.

### Status sub-check dan bukti

| Area | Sub-check | Hasil | Bukti / batasan |
|---|---|---|---|
| Parser | CommonMark 0.31.2 untuk Goldmark dan gomarkdown | Lulus | Diulang: 652/652 masing-masing pada mode pembanding terdokumentasi. |
| Parser | GFM fixture: tabel, task list, strikethrough, autolink pada Goldmark | Lulus terbatas | Fixture sempit lulus; bukan bukti conformance penuh. |
| Parser | GFM spec 0.29 resmi | Lulus untuk profil semantik yang diuji; exact output berbeda | 648 contoh diparse. Bundle `extension.GFM` menghasilkan 14 output mismatch; normalization test menghasilkan 0 semantic mismatches dengan aturan eksplisit: nested `strong` diratakan, sintaks spasi sebelum hard-break diabaikan, dan autolink ditoleransi hanya untuk HTTP(S)/email dengan href sama dengan teks. Pengguna menyetujui Linkify literal URL/email sebagai ekstensi produk. Exact HTML tetap berbeda dari spec. |
| Parser | Fallback AST: raw HTML, fenced code tak dikenal, Mermaid | Lulus terbatas | Probe string output scratch mempertahankan sample tersebut; belum mencakup setiap tipe node atau renderer aplikasi. |
| Keamanan | Sanitasi ESC/BEL menjadi notasi terlihat | Lulus di probe | Input berisi byte ESC/BEL; output probe tidak mengandung byte kontrol tersebut. Belum diuji lewat terminal fisik. |
| TUI | Injeksi key/mouse dan dimensi SimulationScreen tcell | Lulus terbatas | Test simulasi diulang; tidak membuktikan event resize nyata atau TUI produk. |
| TUI | Compile tcell dan Bubble Tea untuk Windows/macOS/Linux amd64 | Lulus dari spike terdahulu | Hanya library compile silang; tidak menunjukkan aplikasi final/runtime. |
| TUI | Lifecycle normal/error/interupsi dan smoke terminal nyata | Belum terbukti | Tidak ada terminal interaktif atau aplikasi produk yang diuji; acceptance dialihkan ke Ticket 02/10. |
| Clipboard | Drag selection, clipboard OS plain text tanpa ANSI, fallback error | Simulasi kelayakan lulus; acceptance produk belum terbukti | Test prototipe memakai event mouse simulasi dan fake clipboard; bukan terminal nyata atau clipboard OS. Acceptance dialihkan ke Ticket 09/10. |
| Browser | Pemanggilan helper tidak menyusun shell command | Source diperiksa, terbatas | Source `pkg/browser` terdahulu; opener nyata tidak dijalankan. |
| Browser | Validasi URL dan activation gate | Lulus di probe | URL absolut HTTP(S) diterima; skema lain/malformed ditolak; opener palsu hanya dipanggil pada aktivasi dan error diteruskan. Tidak menguji integrasi parse/render atau browser OS; acceptance dialihkan ke Ticket 07/10. |
| Dependency | Graph Goldmark + tcell + helper browser | Lulus untuk trio terpilih | 18 module termasuk root (17 dependency); tidak menyertakan package clipboard yang belum dipilih. |
| Dependency | File dan jenis lisensi seluruh dependency graph trio | Lulus untuk inventaris | File lisensi ditemukan untuk 17 dependency; jenis lisensi MIT, Apache-2.0, BSD-2-Clause, BSD-3-Clause. |
| Dependency | Pemeliharaan dan versi transitif | Sebagian, keputusan toolchain disetujui | `go list -m -u all` menampilkan update bagi sejumlah transitif, sedangkan tidak ada update bagi tiga dependency langsung pada saat cek. govulncheck module scan baseline menemukan 3 advisory; scan package pada prototipe tidak menemukan jalur rentan. Pengguna menyetujui Go 1.25 dan pasangan patch minimum pada 19 Oktober 2026; pemeliharaan upstream masih perlu ditinjau lebih jauh. |
| Lingkungan | Go dan download dependency | Lulus | Go 1.27.2 tersedia dan unduhan module berhasil. Test dijalankan dari scratch setelah tool mengizinkan perintah tanpa sandbox. |

### Pengulangan test prototipe — 19 Oktober 2026

Working directory: `.scratch/ticket01/`, yang disediakan pengguna; Go `1.27.2 windows/amd64`; perintah dijalankan tanpa sandbox setelah sandbox WSL gagal start dan izin terminal diberikan.

- `go test ./...` — lulus (`ok stackspike`).
- `go test -v ./...` — lulus. `TestCommonMarkCorpus`: Goldmark 0/652 mismatch dan gomarkdown 0/652 mismatch; `TestGoldmarkGFMAndSafeHTML`, `TestGomarkdownCommonExtensions`, dan `TestTCellSimulationEvents` lulus. Log gomarkdown tetap mencatat tidak ada semantik checkbox task-list pada ekstensi yang dicoba; test lulus sebagai probe pembanding, bukan bukti task list didukung.
- `go list -m all` dan `go list -m all | wc -l` pada module gabungan — hasil `39`, termasuk root; ini graph pembanding lama, bukan stack terpilih.
- `go test -run TestGoldmarkAgainstOfficialGFMSpecExamples -v` — 648 contoh dari corpus [GFM spec 0.29](https://raw.githubusercontent.com/github/cmark-gfm/master/test/spec.txt) diparse; corpus bertanggal 2019-04-06 dan berlisensi CC-BY-SA 4.0. Bundle `extension.GFM` menghasilkan 14 exact-output mismatch; profile Table+TaskList+Strikethrough tanpa Linkify menghasilkan 9. Pengguna menyetujui Linkify URL/email literal sebagai ekstensi proyek di atas GFM 0.29.
- Test normalisasi semantik membandingkan seluruh contoh. Hasil: 0 semantic mismatch setelah (1) meratakan elemen `<strong>` bertingkat, (2) mengabaikan spasi Markdown sebelum `<br />`, dan (3) mengabaikan wrapper anchor tambahan hanya bila href HTTP(S)/mailto valid dan sama dengan isi teks. Normalisasi ini tidak mengklaim HTML byte-identik dan tidak membuka browser; aturan tersebut tercatat eksplisit dalam test scratch. Profil yang diterima untuk proyek adalah GFM + Linkify, sesuai persetujuan pengguna.
- `go test -v ./...` — test CommonMark, fallback AST/sanitasi, validator URL/opener palsu, dan tcell simulation lulus. Test GFM hanya mengukur serta mencatat perbedaan.
- `go list -m -u all` pada graph baseline: update tersedia untuk beberapa transitif (misalnya x/text `v0.31.0` ke `v0.43.0`, x/sys `v0.38.0` ke `v0.49.0`, dan lainnya); tidak ada versi baru yang ditampilkan untuk Goldmark, tcell, dan pkg/browser pada pemeriksaan tersebut.
- `go run golang.org/x/vuln/cmd/govulncheck@latest -scan=module` pada graph baseline gagal dengan 3 temuan: `GO-2026-6629` (`x/text`, fixed `v0.41.0`), `GO-2026-5970` (`x/text`, fixed `v0.39.0`), dan `GO-2026-5024` (`x/sys`, fixed `v0.44.0`). Database tool `govulncheck v1.9.0` terakhir diperbarui 8 Okt 2026.
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` pada module gabungan awal — `No vulnerabilities found` untuk jalur kode probe yang dianalisis. Ini tidak meniadakan temuan scan seluruh module graph atau membuktikan keamanan aplikasi yang belum ada.
- Uji versi terbaru perbaikan di `.scratch/ticket01/approved-stack-graph/`: x/text `v0.43.0` dan x/sys `v0.49.0`; `go test ./...`, `go vet ./...`, `go build ./...`, `go mod verify`, dan scan module lulus/tidak menemukan vulnerability. Namun Go directive naik dari 1.24 ke 1.26.
- Uji pasangan patch minimum di `.scratch/ticket01/patched-go125/`: x/text `v0.41.0` dan x/sys `v0.44.0`; `go.mod` memakai `go 1.25.0`. `go mod tidy`, `go test ./...`, `go vet ./...`, `go build ./...`, `go mod verify`, serta govulncheck module scan berhasil; tidak ada vulnerability ditemukan dan graph berisi 18 module termasuk root. Versi x/text `v0.41.0` menutup kedua advisory x/text yang terdeteksi, dan x/sys `v0.44.0` menutup advisory Windows yang ditemukan.
- Pengguna menyetujui minimum Go 1.25 dan penggunaan pasangan patch minimum pada 19 Oktober 2026. Root repo tetap belum memiliki `go.mod` atau dependency produksi; keputusan ini berlaku saat module aplikasi dibuat, tidak menambahkan file production sekarang.
- Module terpisah di `.scratch/ticket01/approved-stack-graph/` mengunci kandidat langsung Goldmark `v1.8.6`, tcell `v2.13.10`, dan pkg/browser pseudo-version yang disetujui sebagai helper. `go mod tidy`, `go list -m all`, `go mod download all`, `go test ./...`, `go vet ./...`, `go build ./...`, dan `go mod verify` berhasil. `go list -m all` menghasilkan 18 module termasuk root; clipboard dikecualikan karena package final belum dipilih.
- `python audit_licenses.py` menemukan file lisensi pada seluruh 17 dependency non-root. Metadata lisensi teridentifikasi sebagai MIT (Goldmark, go-colorful, go-sixel, uniseg, quant), Apache-2.0 (tcell, gdamore/encoding), BSD-2-Clause (pkg/browser), dan BSD-3-Clause (golang.org/x modules).

Pemeriksaan tanggal module menunjukkan dependency graph baseline mencakup rilis lama, antara lain `soniakeys/quant` (2018), `golang.org/x/xerrors` (2019), `golang.org/x/crypto` (2021), `go-sixel` (2022), dan `golang.org/x/net` (2023); ini sinyal untuk peninjauan pemeliharaan, bukan bukti proyek tersebut tidak dipelihara. `go list -m -u all` sudah dipakai untuk melihat update yang tersedia dan `govulncheck` untuk advisory graph, tetapi belum dilakukan peninjauan respons maintainer/isu satu per satu. Pengguna menyetujui minimum Go 1.25 pada 19 Oktober 2026; patch minimum diuji di Go 1.25, sedangkan versi terbaru yang diuji membutuhkan Go 1.26 dan tidak dipilih. Graph 18 module hanya mewakili parser/TUI/helper browser, tidak mencakup package clipboard.

`safePlainText` dan URL gate adalah probe scratch, bukan kode aplikasi. Test membuktikan perilaku yang diisolasi, bukan interaksi renderer/TUI nyata. Tidak ada clipboard OS, browser opener, atau terminal fisik yang diuji. `.scratch/ticket01/` dipertahankan sesuai permintaan pengguna; tidak ada module produksi di root repo.

Stack yang disetujui tetap Goldmark + `tcell/v2`; pilihan library tidak diubah. Pengguna menyetujui Linkify untuk URL/email literal sebagai perilaku proyek dan menaikkan minimum Go menjadi 1.25 agar dapat memakai patch transitif minimum yang diuji. GFM spec 0.29 telah diparse penuh dan profil semantik hasil normalisasi lulus, meski output HTML exact berbeda. Clipboard tetap memakai abstraksi ber-error dengan pemilihan package final ditunda. TUI nyata, clipboard OS, browser OS, maintainer review penuh, dan graph clipboard belum terbukti dan tetap menjadi acceptance/tindak lanjut pada ticket pemilik, bukan blocker Fase 1. Bukti kelayakan risiko Fase 1 dan persetujuan yang diperlukan lengkap; Ticket 01 ditutup dan Ticket 02 siap, tanpa memulai implementasinya.

### Prototipe seleksi drag tcell — 19 Oktober 2026

- Dibuat prototipe terisolasi di `.scratch/ticket01/selection-prototype/` untuk memetakan event tombol mouse tcell menjadi rentang sel, mengekstrak plain text, lalu memanggil abstraksi `WriteText` yang memakai `Screen.SetClipboard`.
- Pemeriksaan terakhir: Windows, Go 1.27.2; folder `.scratch/ticket01/selection-prototype`; command `go test .`; hasil lulus: `ok stackspike`.
- Test memakai event mouse simulasi dan clipboard tiruan untuk menguji rentang plain text serta ketiadaan byte ESC/CSI. Ini bukti kelayakan spike saja, bukan bukti drag selection di terminal nyata atau keberhasilan clipboard OS. Acceptance aplikasi dipindahkan ke Ticket 09/10.
- Belum ada sesi interaktif nyata. `SetClipboard` tcell tidak memberikan status keberhasilan; prototipe tidak membuktikan seleksi drag di terminal atau clipboard OS. Petunjuk menjalankan test dan mencoba terminal tersedia di README prototipe.

## Sumber upstream

- [Goldmark di pkg.go.dev](https://pkg.go.dev/github.com/yuin/goldmark)
- [gomarkdown di pkg.go.dev](https://pkg.go.dev/github.com/gomarkdown/markdown)
- [CommonMark 0.31.2 spec dan corpus](https://spec.commonmark.org/0.31.2/)
- [Bubble Tea di pkg.go.dev](https://pkg.go.dev/github.com/charmbracelet/bubbletea)
- [tcell v2 di pkg.go.dev](https://pkg.go.dev/github.com/gdamore/tcell/v2)
- [atotto/clipboard di pkg.go.dev](https://pkg.go.dev/github.com/atotto/clipboard)
- [pkg/browser di pkg.go.dev](https://pkg.go.dev/github.com/pkg/browser)
