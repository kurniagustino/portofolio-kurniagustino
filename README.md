# Portfolio Kurnia Gustino

Situs portfolio pribadi berbasis **Go + Fiber + GORM + Tailwind CSS**, dengan
panel admin untuk mengelola seluruh isi situs (proyek, blog, pengalaman,
pendidikan, testimoni, pesan masuk, dan pengaturan situs) tanpa menyentuh kode.

## Fitur

- **Homepage dinamis** — hero, layanan, proyek, riwayat kerja, pendidikan,
  testimoni, dan tulisan terbaru semuanya datang dari database.
- **Halaman publik** — `/resume` (bisa dicetak jadi PDF), `/blog`,
  `/blog/:slug`, `/projects/:slug` (studi kasus masalah → solusi → hasil).
- **Form kontak** dengan honeypot anti-bot + rate-limit (5 pesan/10 menit/IP);
  pesan masuk terlihat di tab **Pesan** pada admin.
- **SEO** — canonical URL, Open Graph & Twitter Card, JSON-LD
  (Person, WebSite, CreativeWork, BlogPosting), `/sitemap.xml` dinamis, dan
  `/robots.txt`.
- **Halaman 404** sendiri untuk rute yang tidak dikenal.
- **Keamanan** — JWT dari env `JWT_SECRET`, rate-limit login
  (10 percobaan/15 menit/IP), password admin di-hash bcrypt.
- **Tailwind CSS produksi** (file statis33 KB, bukan CDN dev).

## Menjalankan

```bash
# 1. Siapkan file .env (lihat bagian bawah)
# 2. Jalankan
go run .
```

Server berjalan di `http://localhost:3000` (ubah dengan env `PORT`).
Migrasi tabel, slug backfill, dan seed data awal (proyek + pengalaman kerja)
berjalan otomatis saat pertama kali start — hanya bila tabel masih kosong.

### Variabel environment (file `.env`)

| Variabel      | Wajib | Keterangan                                              |
| ------------- | ----- | ------------------------------------------------------- |
| `DATABASE_URL`| ya    | DSN MySQL/TiDB, contoh: `user:pass@tcp(host:4000)/db?...`|
| `JWT_SECRET`  | ya    | Secret untuk sesi admin (sudah diisi, jangan dibagikan)  |
| `SITE_URL`    | tidak | URL produksi, mis. `https://kurniagustino.com` — dipakai canonical/OG/sitemap. Bila kosong memakai alamat request. |
| `PORT`        | tidak | Port server (default `3000`)                            |

## Admin

Buka `/login`. Kelola proyek, pesan masuk, pengalaman kerja, pendidikan &
sertifikasi, testimoni, blog, menu navigasi, dan pengaturan situs (foto profil,
CV, WhatsApp, link sosial) dari satu panel.

> Catatan: foto profil, `cv.pdf`, dan nomor WhatsApp diisi lewat **Site
> Settings**; field studi kasus (masalah/solusi/hasil, tahun, kategori, peran)
> diisi lewat edit proyek.

## Mengubah template / styling

Template HTML di-embed ke dalam binary Go (`templates/embed.go`), jadi
**rebuild setelah mengubah template**: `go build`.

Tailwind dibangun dari file statis (bukan CDN). Setelah mengubah class di
template, bangun ulang CSS-nya:

```bash
npx tailwindcss@3.4.17 -i static/src/input.css -o static/css/tailwind.css --minify
```

Hasil build (`static/css/tailwind.css`) ikut di-commit supaya deploy tidak
membutuhkan Node.js.

## Struktur

```
main.go            # konfigurasi Fiber + seluruh route
handlers/          # handler halaman, API admin, CRUD generik, SEO, auth
models/            # model GORM + helper slug
database/          # koneksi, AutoMigrate, seed otomatis
templates/         # semua HTML (di-embed)
public/            # aset statis (gambar proyek, cv.pdf)
static/            # CSS hasil build Tailwind
media/             # hasil unggahan admin (foto, gambar proyek) — tidak di-commit
```
