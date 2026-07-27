---
name: youtube-autoposter
description: Comprehensive YouTube Auto-Poster CLI for AI agents. Inspect accounts, channels, video files, and upload videos with custom thumbnails non-interactively using structured JSON outputs.
---

# 🤖 YouTube Auto-Poster Agent Skill Guide

This skill allows AI agents to inspect, manage, and upload videos to YouTube autonomously using the `./youtube-autoposter` binary.

---

## 📍 Binary Path
- Absolute Path: `/home/gemari-pc/Documents/Sodikin/Agent/siapdigital/youtube-autoposter/youtube-autoposter`
- Relative Path: `./youtube-autoposter/youtube-autoposter`

---

## 🛠️ AI Agent Inspection Commands (All return JSON)

### 1. List Saved Account Profiles
Scan for all saved account tokens (`token.json`, `token_*.json`):
```bash
./youtube-autoposter -list-profiles
```
**JSON Output:**
```json
[
  { "name": "Default Account", "token_file": "token.json" },
  { "name": "Akun Dua", "token_file": "token_akun_dua.json" }
]
```

### 2. List Connected YouTube Channels for Account
Get channel details (Title, ID, Subscriber count, Video count) for an account:
```bash
./youtube-autoposter -profile "akun_dua" -list-channels
```
**JSON Output:**
```json
[
  {
    "id": "UCw83mhlRxJw_U82G_Hmo-mg",
    "title": "Siap Digital",
    "subscriber_count": 0,
    "video_count": 3
  }
]
```

### 3. Scan Local Directory for Video Files
Recursively scan directory for video files (`.mp4`, `.mkv`, `.mov`, etc.):
```bash
./youtube-autoposter -list-videos
```
**JSON Output:**
```json
[
  {
    "path": "/home/.../video.mp4",
    "rel_path": "./video.mp4",
    "size_bytes": 26633830,
    "size_formatted": "25.40 MB"
  }
]
```

---

## 🚀 AI Agent Video Upload Command

Run non-interactively with `-json` flag:

```bash
./youtube-autoposter \
  -profile "akun_dua" \
  -file "./video.mp4" \
  -thumbnail "./cover.png" \
  -title "Automated Upload Title" \
  -description "Uploaded by AI Agent" \
  -tags "ai,golang,tutorial" \
  -privacy "private" \
  -json
```

**JSON Output:**
```json
{
  "id": "dQw4w9WgXcQ",
  "watch_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "duration": "15s",
  "has_custom_thumb": true
}
```
