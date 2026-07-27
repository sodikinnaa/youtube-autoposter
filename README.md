# 🎬 YouTube Auto-Poster CLI (Golang - Clean Architecture)

A lightweight, high-performance YouTube Auto-Poster CLI written in Golang built with **Clean Architecture**. Supports **Interactive User-Friendly CLI (POV User)**, video uploads, custom thumbnails, scheduled publishing, OAuth2 auto-tokens, and an **Interactive Credentials Wizard**.

---

## 📁 Struktur Clean Architecture

```text
youtube-autoposter/
├── internal/
│   ├── cli/                  # 🖥️ Interactive User POV CLI Interface
│   │   └── interactive.go
│   ├── domain/               # Core Domain Entities, Models, & Interfaces
│   │   └── video.go
│   ├── infrastructure/       # External Implementations (OAuth & YouTube API)
│   │   ├── oauth/
│   │   │   ├── authenticator.go
│   │   │   └── wizard.go     # Interactive Credentials Wizard
│   │   └── youtube/
│   │       └── client.go     # YouTube Data API v3 Client & Thumbnail Engine
│   └── usecase/              # Application Business Logic
│       └── upload_video.go   # Upload Orchestration UseCase
├── main.go                   # Entry Point & Dependency Injection (DI)
├── go.mod / go.sum
├── .gitignore
├── README.md
└── youtube-autoposter       # Compiled Binary Executable
```

---

## 📖 PANDUAN PENGGUNAAN BINARY

### 1. Build Binary

```bash
cd youtube-autoposter
go build -o youtube-autoposter .
```

---

### 2. Cara Menggunakan: Mode Interaktif (User POV) ⭐ Recommended

Cukup jalankan binary tanpa argumen/flag (atau dengan flag `-i`):

```bash
./youtube-autoposter
```

Kamu akan dipandu dengan wizard terminal interaktif yang sangat ramah pengguna:

```text
=======================================================
🎬 YOUTUBE AUTO-POSTER - INTERACTIVE WIZARD
=======================================================
Silakan isi detail video di bawah ini (atau tekan Enter untuk default):

🎥 Path File Video (contoh: ./my_video.mp4): ./video_tutorial.mp4

🖼️  Path Custom Thumbnail (opsional, contoh: ./thumb.jpg, tekan Enter jika tidak ada): ./cover.png

📝 Judul Video [default: video_tutorial]: Tutorial Golang Clean Architecture

📄 Deskripsi Video (opsional, tekan Enter jika kosong): Video panduan YouTube Auto-Poster CLI.

🏷️  Tags Video (pisahkan dengan koma, contoh: coding,golang,tutorial): golang,clean-architecture,cli

🔒 Status Privasi Video:
   [1] Private  (Hanya kamu yang bisa lihat - Default)
   [2] Public   (Bisa ditonton siapa saja)
   [3] Unlisted (Hanya yang punya link yang bisa lihat)
Pilihan Status (1/2/3) [default: 1]: 1

⏱️  Jadwalkan Tayang Otomatis? (y/N): n

=======================================================
📋 RINGKASAN KONFIGURASI UPLOAD
=======================================================
🎥 File Video : ./video_tutorial.mp4
🖼️  Thumbnail  : ./cover.png
📝 Judul      : Tutorial Golang Clean Architecture
📄 Deskripsi  : Video panduan YouTube Auto-Poster CLI.
🏷️  Tags       : golang, clean-architecture, cli
🔒 Privasi    : private
=======================================================
Apakah konfigurasi di atas sudah sesuai dan siap upload? (Y/n): y
```

---

### 3. Cara Menggunakan: Mode Flag (Automation / Script POV)

Untuk kebutuhan otomatisasi / cronjob / shell script, kamu bisa tetap mengoperkan flag secara langsung:

```bash
./youtube-autoposter \
  -file ./video_tutorial.mp4 \
  -thumbnail ./cover.png \
  -title "Tutorial Golang Clean Architecture" \
  -description "Deskripsi video otomatis" \
  -tags "golang,tutorial" \
  -privacy private
```

---

## ⚙️ Daftar Flag / Argument Binary

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-i` | Aktifkan Mode Interaktif (User POV) | `false` |
| `-file` | Path file video | *(Wajib jika non-interaktif)* |
| `-thumbnail` | Path file gambar thumbnail (`.png`, `.jpg`, `.webp`) | `""` |
| `-title` | Judul video YouTube | Nama file video tanpa ekstensi |
| `-description` | Deskripsi video | `""` |
| `-tags` | Tags dipisahkan koma | `""` |
| `-category` | Category ID (22=People & Blogs, 28=Sci&Tech, 27=Edu) | `22` |
| `-privacy` | Status Privasi (`public`, `private`, `unlisted`) | `private` |
| `-publish-at` | Waktu tayang otomatis (format RFC3339, misal `2026-08-01T15:00:00Z`) | `""` |
| `-secret` | Path lokasi file OAuth JSON | `client_secret.json` |
| `-token` | Path lokasi simpanan token OAuth | `token.json` |

---

## 🔒 Catatan Keamanan
File `client_secret.json` dan `token.json` **secara otomatis di-ignore oleh `.gitignore`**, sehingga aman dari kebocoran saat di-commit ke Git.
