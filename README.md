# 🎬 YouTube Auto-Poster (Golang)

A lightweight, powerful YouTube Auto-Poster CLI written in Golang using YouTube Data API v3 and OAuth2 authorization.

---

## 📁 Struktur Project

Tool ini ditempatkan di subfolder terisolasi dengan `go.mod` sendiri:
```text
youtube-autoposter/
├── auth.go             # OAuth2 Flow & Token Manager (Local Callback Server)
├── uploader.go         # Core Upload Engine (YouTube Data API v3)
├── main.go             # CLI Interface & Argument Parser
├── go.mod / go.sum     # Go Module Dependencies
├── .gitignore          # Proteksi client_secret.json & token.json dari Git
└── README.md           # Dokumentasi Penggunaan
```

---

## 🔑 Langkah Persiapan (Google Cloud Console)

Sebelum menjalankan tool ini, kamu perlu membuat OAuth 2.0 Client ID dari Google Cloud Console:

1. Buka [Google Cloud Console](https://console.cloud.google.com/).
2. Buat Project baru (misal: `YouTube Auto Poster`).
3. Buka menu **APIs & Services** > **Library**, cari **YouTube Data API v3**, lalu klik **Enable**.
4. Buka menu **OAuth consent screen**:
   - Pilih **External** lalu klik *Create*.
   - Isi Nama Aplikasi dan Email Dukungan.
   - Tambahkan Scope: `.../auth/youtube.upload` dan `.../auth/youtube.readonly`.
   - Pada bagian **Test users**, tambahkan alamat email Google / YouTube milikmu.
5. Buka menu **Credentials** > **Create Credentials** > **OAuth client ID**:
   - Application type: **Desktop app** (atau Web application dengan Authorized redirect URI: `http://localhost:8080/callback`).
   - Simpan dan download file JSON credentials-nya.
6. Rename file JSON tersebut menjadi `client_secret.json` dan letakkan di dalam folder `youtube-autoposter/`.

---

## 🚀 Cara Penggunaan

### 1. Direct Build / Run dengan Go

```bash
cd youtube-autoposter

# Build binary
go build -o youtube-autoposter .

# Jalankan bantuan CLI
./youtube-autoposter
```

### 2. Contoh Command Upload

#### A. Upload Video + Custom Thumbnail
```bash
./youtube-autoposter \
  -file ./video_tutorial.mp4 \
  -thumbnail ./cover_thumbnail.png \
  -title "Tutorial Golang Lanjutan" \
  -description "Dalam video ini kita belajar membuat YouTube Auto-Poster." \
  -tags "golang,tutorial,coding" \
  -privacy private
```

#### B. Upload Video Publik
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

## ⚙️ Daftar Flag / Argument

| Flag | Description | Default |
| --- | --- | --- |
| `-file` | Path ke file video (contoh: `./video.mp4`) | *(Wajib)* |
| `-thumbnail` | Path ke file gambar thumbnail (contoh: `./thumb.png` / `./thumb.jpg`) | `""` |
| `-title` | Judul video YouTube | Nama file tanpa ekstensi |
| `-description` | Deskripsi lengkap video | `""` |
| `-tags` | Tags video dipisahkan koma | `""` |
| `-category` | ID Kategori (22=People & Blogs, 28=Sci&Tech, 27=Edu) | `22` |
| `-privacy` | Status Privasi (`public`, `private`, `unlisted`) | `private` |
| `-publish-at` | Waktu tayang otomatis (Format RFC3339, misal `2026-08-01T15:00:00Z`) | `""` |
| `-secret` | Path file OAuth JSON | `client_secret.json` |
| `-token` | Path file token yang disimpan | `token.json` |

---

## 🔒 Catatan Keamanan
File `client_secret.json` dan `token.json` **secara otomatis di-ignore oleh `.gitignore`**, sehingga aman dan tidak akan masuk ke git commit repository utama!
