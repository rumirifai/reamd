# Checklist ticket pengembangan `reamd`

Ticket di folder ini adalah unit tugas untuk agentic AI. Ticket bersifat lokal dan menjadi bagian dokumentasi proyek; ticket tidak diterbitkan ke GitHub Issues.

Urutan fase dan gerbang proyek tetap mengacu pada [workflow pengembangan](../development-workflow.md). Perilaku produk dan cara verifikasinya tetap mengacu pada [spesifikasi command](../command-spec.md) dan [kriteria penerimaan](../acceptance-criteria.md).

## Cara menggunakan

1. Pilih ticket berstatus `ready-for-agent` yang blocker-nya sudah selesai.
2. Mulai satu session dengan satu ticket sebagai ruang lingkup.
3. Minta agent membaca `AGENTS.md`, dokumen acuan yang disebut ticket, dan ticket itu sendiri.
4. Agent memeriksa keadaan repo, lalu mengerjakan hanya hasil yang tercantum.
5. Tandai checklist berdasarkan bukti. Jangan menandai smoke test manual lulus jika belum benar-benar dijalankan.
6. Setelah blocker terpenuhi, perbarui status ticket berikutnya dari `blocked` menjadi `ready-for-agent`.

Fase 0 workflow bukan ticket produk tersendiri. Checklist persiapan fase 0 diulang sebelum setiap ticket: periksa perubahan repo, pastikan acceptance criteria dan keputusan yang relevan jelas, lalu batasi tugas pada satu ticket.

## Urutan dan dependency

| Ticket | Fase | Judul | Blocked by |
|---|---:|---|---|
| [01](01-stack-spike.md) | 1 | Evaluasi stack dan risiko teknis | Tidak ada |
| [02](02-cli-vertical-slice.md) | 2-3 | Bangun CLI dan vertical slice TUI | 01 |
| [03](03-commonmark-core.md) | 4 | Render elemen inti CommonMark | 02 |
| [04](04-gfm-responsive-tables.md) | 4 | Lengkapi GFM dan tabel responsif | 03 |
| [05](05-code-and-safe-content.md) | 4 | Render code dan amankan konten Markdown | 02 |
| [06](06-themes-zoom-resize.md) | 5 | Tambahkan tema, zoom, dan resize | 03 |
| [07](07-links-and-anchors.md) | 5 | Navigasi anchor dan aktivasi link | 01, 03 |
| [08](08-mouse-scrollbar.md) | 5 | Tambahkan scroll mouse dan scrollbar | 02 |
| [09](09-selection-clipboard.md) | 5 | Seleksi teks dan copy ke clipboard | 01, 08 |
| [10](10-cross-platform-verification.md) | 6 | Verifikasi lintas platform | 04-09 |
| [11](11-distribution.md) | 7 | Siapkan distribusi dan instalasi | 10 |

Ticket 04, 05, 06, 07, dan 08 dapat dikerjakan paralel setelah blocker masing-masing selesai. Ticket 09 menunggu dukungan scroll mouse sebagai fondasi pemetaan seleksi di area dokumen.

## Aturan ticket

- Status `ready-for-agent` berarti dapat dimulai; `blocked` berarti masih menunggu ticket yang tercantum.
- Selesaikan checklist acceptance dan verifikasi dalam ticket sebelum menandainya selesai.
- Perubahan perilaku harus tetap sesuai sumber kebenaran di `docs/command-spec.md`; jika ada konflik, berhenti dan minta keputusan.
- Fitur di luar MVP tidak ditambahkan tanpa persetujuan.
- Perbarui ticket dan dokumen acuan bila pekerjaan menemukan keputusan atau acceptance criteria yang perlu diklarifikasi.
