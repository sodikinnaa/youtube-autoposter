# 🎬 YouTube Auto-Poster CLI (Golang - Clean Architecture)

Aplikasi CLI otomatisasi upload video YouTube berbasis Golang dengan arsitektur **Clean Architecture**. Dirancang sangat ramah pengguna (**User POV Interactive Wizard**), mendukung upload video, thumbnail kustom, penjadwalan tayang (*scheduled publish*), otentikasi otomatis Google OAuth2, dan output JSON untuk integrasi AI Agent.

---

## 🧭 ALUR PENUH PENGGUNAAN (POV USER AWAM / DEVELOPER)

Aplikasi ini menggunakan alur **3 Langkah Mudah**:

```text
[1. Jalankan Binary] ➡️ [2. Verifikasi Credentials & Channel] ➡️ [3. Isi Detail & Upload Video]
```

---

## 📌 LANGKAH 1: MENJALANKAN BINARY

### ⚡ Opsi Fast Install / Download Binary via `curl`:

#### Linux (AMD64):
```bash
curl -sSL https://raw.githubusercontent.com/sodikinnaa/youtube-autoposter/main/install.sh | bash
```

atau download binary rilis secara langsung:

```bash
curl -sSL https://github.com/sodikinnaa/youtube-autoposter/releases/latest/download/youtube-autoposter-linux-amd64 -o youtube-autoposter && chmod +x youtube-autoposter
./youtube-autoposter
```

#### Build dari Source Code:
```bash
cd youtube-autoposter
go build -o youtube-autoposter .
./youtube-autoposter
```

---

## 📌 LANGKAH 2: SETUP KREDENSIAL & OTENTIKASI (SEKALI DI AWAL)

### A. Jika `client_secret.json` Belum Ada (Wizard Otomatis)
Jika aplikasi baru pertama kali dijalankan dan belum ada file kredensial Google, aplikasi akan **otomatis membuka Wizard Interaktif**:

```text
⚠️ File credentials 'client_secret.json' tidak ditemukan. Membuka wizard penyiapan...

=======================================================
🧙 WIZARD SETUP GOOGLE OAUTH CREDENTIALS
=======================================================
Pilih cara penyiapan credentials:
  [1] Input Client ID & Client Secret secara manual
  [2] Masukkan path file JSON credentials yang sudah didownload
  [3] Lihat Panduan Google Cloud Console & Keluar

Pilihan Kamu (1/2/3): 1
```

1. **Pilih `1`**: Cukup paste `Client ID` dan `Client Secret` milikmu. Aplikasi akan **membuat file `client_secret.json` secara otomatis**.
2. **Pilih `2`**: Jika kamu sudah mendownload file JSON dari Google Cloud Console (misal `~/Downloads/client_secret.json`), cukup ketik lokasinya dan aplikasi akan menyalinnya.

#### 💡 Cara Singkat Mendapatkan Client ID & Client Secret Google:
1. Buka [Google Cloud Console](https://console.cloud.google.com/).
2. Buat Project baru > Buka **APIs & Services** > **Library** > Cari & aktifkan **YouTube Data API v3**.
3. Di menu **OAuth consent screen**: Pilih **External**, isi nama app, dan tambahkan email kamu di bagian **Test users**.
4. Di menu **Credentials** > **Create Credentials** > **OAuth client ID**: Pilih **Desktop app**.
5. Salin *Client ID* & *Client Secret*-nya.

---

### B. Otentikasi Akun Google di Browser (Sekali Saja)
Setelah kredensial ada, aplikasi akan memberikan link otentikasi browser:

```text
=======================================================
🔑 PERLU OTENTIKASI YOUTUBE OAUTH2
Buka URL berikut di browser kamu untuk memberikan izin:

https://accounts.google.com/o/oauth2/auth?...

Sedang menunggu callback otomatis di http://localhost:8080/callback...
=======================================================
```

1. Klik / buka link di atas pada browser.
2. Login akun YouTube milikmu dan klik **Izinkan / Allow**.
3. Browser akan menampilkan **"Otentikasi Berhasil!"**.
4. Token disimpan otomatis di `token.json`. **Upload berikutnya tidak perlu login browser lagi!**

---

### C. Verifikasi Channel YouTube Terhubung
Sebelum menanyakan file video, aplikasi **selalu menampilkan info channel yang terhubung**:

```text
=======================================================
📺 CHANNEL YOUTUBE TERHUBUNG SUKSES
=======================================================
👤 Nama Channel : Coding Bersama Sodiq
🆔 Channel ID   : UC1234567890abcdef
👥 Subscribers  : 1,250
📹 Total Video  : 42
=======================================================
Lanjutkan upload video ke channel di atas? (Y/n): y
```

Ketik `y` atau tekan Enter untuk melanjutkan.

---

## 📌 LANGKAH 3: INTERACTIVE UPLOAD WIZARD (POV USER)

Setelah channel terkonfirmasi, kamu akan dipandu mengisi informasi video satu per satu:

### 1. File Video
```text
🎥 Path File Video (contoh: ./my_video.mp4): ./tutorial.mp4
```
*Aplikasi akan otomatis mengecek keberadaan file. Jika file tidak ada, kamu akan diminta memasukkan path yang benar.*

### 2. Custom Thumbnail (Opsional)
```text
🖼️ Path Custom Thumbnail (opsional, contoh: ./thumb.jpg, tekan Enter jika tidak ada): ./cover.png
```
*Mendukung file `.png`, `.jpg`, `.jpeg`, atau `.webp`. Tekan Enter jika ingin menggunakan thumbnail default YouTube.*

### 3. Judul Video
```text
📝 Judul Video [default: tutorial]: Belajar Golang Clean Architecture
```
*Tekan Enter langsung jika ingin memakai nama file sebagai judul.*

### 4. Deskripsi Video (Opsional)
```text
📄 Deskripsi Video (opsional, tekan Enter jika kosong): Panduan lengkap membuat CLI YouTube Auto-Poster dengan Clean Architecture.
```

### 5. Tags Video (Opsional)
```text
🏷️ Tags Video (pisahkan dengan koma, contoh: coding,golang,tutorial): golang,clean-architecture,cli
```

### 6. Status Privasi
```text
🔒 Status Privasi Video:
   [1] Private  (Hanya kamu yang bisa lihat - Default)
   [2] Public   (Bisa ditonton siapa saja)
   [3] Unlisted (Hanya yang punya link yang bisa lihat)
Pilihan Status (1/2/3) [default: 1]: 1
```

### 7. Jadwal Tayang Otomatis / Scheduled Publish (Opsional)
```text
⏱️ Jadwalkan Tayang Otomatis? (y/N): y
Masukkan tanggal & waktu (format: YYYY-MM-DD HH:MM, contoh: 2026-08-01 15:00): 2026-08-01 15:00
✅ Tayang Otomatis dijadwalkan pada: 01 Aug 2026 15:00 WIB
```

---

### 📋 RINGKASAN & KONFIRMASI AKHIR

Sebelum proses upload berjalan, aplikasi akan menampilkan ringkasan utuh:

```text
=======================================================
📋 RINGKASAN KONFIGURASI UPLOAD
=======================================================
📺 Channel Target: Coding Bersama Sodiq (UC1234567890abcdef)
🎥 File Video    : ./tutorial.mp4
🖼️ Thumbnail     : ./cover.png
📝 Judul         : Belajar Golang Clean Architecture
📄 Deskripsi     : Panduan lengkap membuat CLI YouTube Auto-Poster...
🏷️ Tags          : golang, clean-architecture, cli
🔒 Privasi       : private
⏱️ Jadwal        : 2026-08-01T15:00:00+07:00
=======================================================
Apakah konfigurasi di atas sudah sesuai dan siap upload? (Y/n): y
```

Ketik `y` / tekan Enter, lalu proses upload akan berjalan:

```text
🚀 Memulai proses upload video ke YouTube...

=======================================================
🎉 UPLOAD VIDEO BERHASIL DIPROSES!
📺 Video ID    : dQw4w9WgXcQ
🔗 Link Watch  : https://www.youtube.com/watch?v=dQw4w9WgXcQ
⏱️ Durasi Upload: 12s
🖼️ Thumbnail   : Custom Thumbnail Berhasil Dipasang
=======================================================
```

---

## 🔌 INTEGRASI MODEL CONTEXT PROTOCOL (MCP) SERVER

Binary `./youtube-autoposter` ini dapat dijalankan langsung sebagai **MCP Server (Model Context Protocol)** berbasis `stdio` JSON-RPC 2.0.

### 1. Cara Menjalankan MCP Server
```bash
./youtube-autoposter -mcp
```

### 2. Cara Menambahkan ke MCP Client / AI Agent Configuration

Tambahkan konfigurasi berikut ke file konfigurasi MCP (seperti `claude_desktop_config.json`, `mcp.json`, atau `agy.json`):

```json
{
  "mcpServers": {
    "youtube-autoposter": {
      "command": "/home/gemari-pc/Documents/Sodikin/Agent/siapdigital/youtube-autoposter/youtube-autoposter",
      "args": ["-mcp"]
    }
  }
}
```

### 🛠️ Tool MCP yang Disediakan:
1. `list_profiles`: Mengembalikan daftar profile akun YouTube yang tersimpan.
2. `list_channels`: Mengembalikan daftar channel YouTube yang terhubung pada profile.
3. `list_videos`: Pindai folder secara rekursif untuk mendeteksi file video (`.mp4`, `.mkv`, dll).
4. `upload_video`: Meng-upload video ke YouTube dengan thumbnail, deskripsi, tags, dan privasi.

---

## 🤖 UNTUK AI AGENT & AUTOMATION SCRIPT (NON-INTERAKTIF)

Jika kamu ingin menjalankan via **Scripting / Cronjob / AI Agent (Google Antigravity / Gemini / Claude)**, gunakan flag CLI langsung:

```bash
./youtube-autoposter \
  -file "./tutorial.mp4" \
  -thumbnail "./cover.png" \
  -title "Belajar Golang Clean Architecture" \
  -description "Deskripsi otomatis" \
  -tags "golang,tutorial" \
  -privacy "private" \
  -json
```

### Format Respon JSON (`-json`):
```json
{
  "id": "dQw4w9WgXcQ",
  "watch_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "duration": "12s",
  "has_custom_thumb": true
}
```

---

## ⚙️ TABEL FLAG CLI COMPLETE

| Flag | Fungsi | Default Value | Contoh Penggunaan |
| --- | --- | --- | --- |
| `-i` | Paksa buka Mode Interaktif Wizard | `false` | `./youtube-autoposter -i` |
| `-file` | Path file video | *(Wajib jika non-interaktif)* | `-file "./video.mp4"` |
| `-thumbnail` | Path file gambar thumbnail | `""` | `-thumbnail "./cover.png"` |
| `-title` | Judul video YouTube | Nama file tanpa ekstensi | `-title "Judul Keren"` |
| `-description` | Deskripsi lengkap video | `""` | `-description "Deskripsi"` |
| `-tags` | Tags dipisahkan koma | `""` | `-tags "golang,coding"` |
| `-category` | Category ID YouTube | `22` (People & Blogs) | `-category "28"` (Sci&Tech) |
| `-privacy` | Privasi (`public`, `private`, `unlisted`) | `private` | `-privacy "public"` |
| `-publish-at` | Jadwal tayang (RFC3339) | `""` | `-publish-at "2026-08-01T15:00:00Z"` |
| `-json` | Output JSON terstruktur untuk AI Agent | `false` | `-json` |
| `-secret` | Path file OAuth client JSON | `client_secret.json` | `-secret "./my_secret.json"` |
| `-token` | Path simpanan token OAuth | `token.json` | `-token "./my_token.json"` |

---

## 📁 STRUKTUR CLEAN ARCHITECTURE PROJECT

```text
youtube-autoposter/
├── internal/
│   ├── cli/                  # 🖥️ Interactive User POV CLI Interface & Prompts
│   │   └── interactive.go
│   ├── domain/               # 💎 Domain Models, Entities, & Interfaces
│   │   └── video.go
│   ├── infrastructure/       # 🔌 Implementasi External (OAuth2 & YouTube Data API)
│   │   ├── oauth/
│   │   │   ├── authenticator.go
│   │   │   └── wizard.go     # Interactive Credential Wizard
│   │   └── youtube/
│   │       └── client.go     # YouTube Service Client & Thumbnail Uploader
│   └── usecase/              # ⚙️ Application Business Logic
│       ├── get_channel_info.go # UseCase Verifikasi Channel YouTube
│       └── upload_video.go   # UseCase Orkestrasi Upload Video
├── main.go                   # 🚀 Entry Point CLI & Dependency Injection
├── SKILL.md                  # 🤖 AI Agent Skill Manifest
├── .gitignore                # 🔒 Proteksi Kredensial & Token dari Git
├── README.md                 # 📖 Panduan Lengkap
└── youtube-autoposter       # Executable Binary Go
```

---

## 🔒 CATATAN KEAMANAN (SECURITY)
File `client_secret.json` dan `token.json` **secara otomatis di-ignore oleh `.gitignore`**. File kredensial dan token rahasia milikmu **tidak akan pernah** masuk ke komit Git ataupun ter-push ke repository publik/private.
