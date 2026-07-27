# 🎬 YouTube Auto-Poster CLI (Golang)

Tool CLI ringan & cepat untuk mengunggah video ke YouTube otomatis menggunakan Golang. Dibangun dengan pendekatan **Clean Architecture**, mendukung upload video, custom thumbnail, penjadwalan tayang (*scheduled publish*), multi-akun, mode interaktif wizard, hingga integrasi native **MCP (Model Context Protocol)** dan **AI Agent**.

---

## ⚡ INSTALASI CEPAT (1-LINE COMMAND)

Kamu bisa langsung mendownload binary siap pakai via `curl`:

```bash
curl -sSL https://raw.githubusercontent.com/sodikinnaa/youtube-autoposter/main/install.sh | bash
```

Perintah di atas akan mendeteksi OS & arsitektur komputermu, mengunduh binary rilis terbaru, dan menyertakan panduan AI Agent (`SKILL.md`). Jika binary sudah ada, installer akan otomatis melakukan update ke versi terbaru.

---

## 🚀 CARA PENGGUNAAN

### 1. Mode Interaktif (Wizard Termudah) ⭐ Recommended
Cukup jalankan binary tanpa argumen apa pun:

```bash
./youtube-autoposter
```

Kamu akan dipandu selangkah demi selangkah melalui wizard interaktif di terminal:
1. **Pilih Akun**: Pilih dari akun yang tersimpan atau tambah akun baru.
2. **Koneksi Channel**: Aplikasi mengecek & menampilkan detail channel YouTube yang terhubung.
3. **Pilih Video**: Otomatis mendeteksi file video (`.mp4`, `.mkv`, `.mov`) di folder project & sub-folder, atau masukkan path manual.
4. **Thumbnail & Metadata**: Masukkan thumbnail kustom, judul, deskripsi, tags, privasi (*public / private / unlisted*), serta jadwal tayang.
5. **Konfirmasi**: Tinjau ringkasan sebelum proses upload berjalan.

---

### 2. Mode Scripting & Otomatisasi (Flag Mode)
Untuk kebutuhan cronjob, shell script, atau otomasi:

```bash
./youtube-autoposter \
  -file "./video.mp4" \
  -thumbnail "./cover.jpg" \
  -title "Tutorial Golang Lanjutan" \
  -description "Deskripsi video otomatis" \
  -tags "golang,tutorial" \
  -privacy private
```

---

### 3. Integrasi MCP Server (Model Context Protocol)
Aplikasi ini dapat langsung dijadikan **MCP Server** berbasis `stdio` JSON-RPC 2.0 untuk AI Client (seperti Antigravity, Claude Desktop, Cursor, atau VS Code MCP Extension):

```bash
./youtube-autoposter -mcp
```

Contoh konfigurasi di MCP Client (`mcp.json` / `claude_desktop_config.json`):
```json
{
  "mcpServers": {
    "youtube-autoposter": {
      "command": "/path/to/youtube-autoposter",
      "args": ["-mcp"]
    }
  }
}
```

---

### 4. Perintah Khusus AI Agent (`-json`)
Untuk AI Agent yang ingin menginspeksi lingkungan atau mendapatkan respon terstruktur:

- **List Profiles**: `./youtube-autoposter -list-profiles`
- **List Channels**: `./youtube-autoposter -profile "akun_dua" -list-channels`
- **Scan Videos**: `./youtube-autoposter -list-videos`
- **Upload with JSON Output**: `./youtube-autoposter -file "./video.mp4" -json`

---

## ⚙️ TABEL FLAG LENGKAP

| Flag | Fungsi | Default |
| --- | --- | --- |
| `-i` | Buka Mode Interaktif Wizard | `false` |
| `-file` | Path file video yang di-upload | *(Wajib jika non-interaktif)* |
| `-thumbnail` | Path file gambar thumbnail | `""` |
| `-title` | Judul video | Nama file tanpa ekstensi |
| `-description` | Deskripsi video | `""` |
| `-tags` | Tags dipisahkan koma | `""` |
| `-category` | ID Kategori (22=People & Blogs, 28=Sci&Tech) | `22` |
| `-privacy` | Privasi (`public`, `private`, `unlisted`) | `private` |
| `-publish-at` | Waktu tayang otomatis (format RFC3339) | `""` |
| `-profile` | Alias profile akun yang digunakan | `""` |
| `-mcp` | Jalankan sebagai MCP Server (stdio) | `false` |
| `-json` | Output respon JSON terstruktur | `false` |
| `-secret` | Path file OAuth client JSON | `client_secret.json` |
| `-token` | Path simpanan token OAuth | `token.json` |

---

## 🔒 CATATAN KEAMANAN (SECURITY NOTE)

Keamanan kredensial dan akunmu adalah prioritas utama:

- **Strict `.gitignore`**: Seluruh file kredensial Google OAuth (`client_secret*.json`) dan token otentikasi (`token*.json`) secara ketat di-ignore oleh `.gitignore`.
- **Lokal & Terisolasi**: Token akses disimpan secara lokal di mesin komputermu dan **tidak akan pernah** ikut ter-commit atau ter-push ke repository Git.
- **Bebas Hardcoded Secrets**: Tidak ada kode rahasia atau API key yang ditanam di dalam kode sumber (*source code*).

---

## 📄 LISENSI

Proyek ini dilisensikan di bawah **[MIT License](LICENSE)**. Kamu bebas menggunakan, mengubah, dan mendistribusikan ulang kode ini sesuai ketentuan lisensi MIT.
