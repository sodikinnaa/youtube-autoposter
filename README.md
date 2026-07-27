# 🎬 YouTube Auto-Poster CLI (Golang - Clean Architecture)

A lightweight, high-performance YouTube Auto-Poster CLI written in Golang built with **Clean Architecture**, supporting video uploads, custom thumbnails, scheduled publishing, OAuth2 auto-tokens, and an **Interactive Setup Wizard**.

---

## 📁 Struktur Clean Architecture

```text
youtube-autoposter/
├── internal/
│   ├── domain/               # Core Domain Entities, Models, & Interfaces
│   │   └── video.go
│   ├── infrastructure/       # External Implementations (OAuth & YouTube API)
│   │   ├── oauth/
│   │   │   ├── authenticator.go
│   │   │   └── wizard.go     # Interactive CLI Setup Wizard
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

## 📖 PANDUAN LENGKAP PENGGUNAAN BINARY

### 1. Build Binary (Satu Kali)

Buka terminal di folder `youtube-autoposter`, lalu jalankan:

```bash
go build -o youtube-autoposter .
```

*File binary `./youtube-autoposter` akan dibuat.*

---

### 2. Jalankan Binary (Interactive Setup Wizard)

Jika kamu pertama kali menjalankan binary dan **belum** memiliki file `client_secret.json`, jalankan:

```bash
./youtube-autoposter -file ./sample.mp4
```

Aplikasi akan otomatis mendeteksi file credentials belum ada dan membuka **Wizard Interaktif**:

```text
⚠️ File credentials 'client_secret.json' tidak ditemukan. Membuka wizard penyiapan...

=======================================================
🧙 WIZARD SETUP GOOGLE OAUTH CREDENTIALS
=======================================================
File 'client_secret.json' belum ditemukan.

Pilih cara penyiapan credentials:
  [1] Input Client ID & Client Secret secara manual
  [2] Masukkan path file JSON credentials yang sudah didownload
  [3] Lihat Panduan Google Cloud Console & Keluar

Pilihan Kamu (1/2/3): 1

-------------------------------------------------------
📝 Masukkan Client ID & Client Secret
(Dapatkan dari Google Cloud Console > Credentials > OAuth 2.0 Client ID)
-------------------------------------------------------
Client ID     : 123456789-abc.apps.googleusercontent.com
Client Secret : GOCS-xxxxxx

✅ Berhasil membuat file credentials di 'client_secret.json'!
```

---

### 3. Otentikasi Satu Kali (Browser OAuth Login)

Setelah credentials dikonfigurasi, aplikasi akan membuka link otentikasi browser:

```text
=======================================================
🔑 PERLU OTENTIKASI YOUTUBE OAUTH2
Buka URL berikut di browser kamu untuk memberikan izin:

https://accounts.google.com/o/oauth2/auth?...

Sedang menunggu callback otomatis di http://localhost:8080/callback...
=======================================================
```

1. Buka link tersebut di browser.
2. Login akun YouTube milikmu dan beri izin akses upload video.
3. Halaman browser akan menampilkan *"Otentikasi Berhasil!"*.
4. Token otentikasi disimpan otomatis ke `token.json`. **Untuk upload berikutnya, kamu tidak perlu login lagi.**

---

### 4. Contoh Command Upload Lengkap

#### A. Upload Video + Custom Thumbnail (Private)
```bash
./youtube-autoposter \
  -file ./video_tutorial.mp4 \
  -thumbnail ./cover_thumb.png \
  -title "Tutorial Golang Clean Architecture" \
  -description "Dalam video ini kita belajar membuat YouTube Auto-Poster." \
  -tags "golang,tutorial,clean-architecture" \
  -privacy private
```

#### B. Upload Video Langsung Publik
```bash
./youtube-autoposter \
  -file ./my_video.mp4 \
  -title "Belajar Web Development" \
  -privacy public
```

#### C. Scheduled Publish (Jadwal Tayang Otomatis)
```bash
./youtube-autoposter \
  -file ./scheduled_video.mp4 \
  -thumbnail ./thumb.jpg \
  -title "Video Tayang Besok Jam 3 Sore" \
  -publish-at "2026-08-01T15:00:00Z"
```

---

## ⚙️ Daftar Flag / Argument Binary

| Flag | Deskripsi | Default |
| --- | --- | --- |
| `-file` | Path file video yang mau di-upload | *(Wajib)* |
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
