---
name: youtube-autoposter
description: Upload videos, thumbnails, set privacy, and schedule publish on YouTube automatically via Golang CLI binary.
---

# YouTube Auto-Poster Skill Guide

This skill allows AI agents to automatically upload videos, assign custom thumbnails, and schedule video publications to YouTube using the `./youtube-autoposter` binary.

## Binary Location
Absolute Path: `/home/gemari-pc/Documents/Sodikin/Agent/siapdigital/youtube-autoposter/youtube-autoposter`
Relative Path: `./youtube-autoposter/youtube-autoposter`

---

## How AI Agents Should Execute Uploads

Run the binary non-interactively via `run_command` tool using flags and `-json` output format:

```bash
/home/gemari-pc/Documents/Sodikin/Agent/siapdigital/youtube-autoposter/youtube-autoposter \
  -file "/path/to/video.mp4" \
  -thumbnail "/path/to/thumbnail.jpg" \
  -title "Video Title" \
  -description "Detailed description" \
  -tags "golang,ai,tutorial" \
  -privacy "private" \
  -json
```

### JSON Output Format
When `-json` flag is provided, the binary outputs structured JSON:

**Success Response:**
```json
{
  "id": "dQw4w9WgXcQ",
  "watch_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "duration": "12s",
  "has_custom_thumb": true
}
```

**Error Response:**
```json
{
  "status": "error",
  "message": "file video tidak ditemukan"
}
```

---

## Command Flags Reference for AI

- `-file` (string, required): Absolute path to video file (`.mp4`, `.mkv`, `.mov`).
- `-thumbnail` (string, optional): Absolute path to thumbnail image (`.png`, `.jpg`, `.webp`).
- `-title` (string, optional): Video title (defaults to filename).
- `-description` (string, optional): Video description.
- `-tags` (string, optional): Comma-separated list of tags.
- `-privacy` (string, optional): `private` (default), `public`, or `unlisted`.
- `-publish-at` (string, optional): RFC3339 timestamp for scheduled publishing (e.g. `2026-08-01T15:00:00Z`).
- `-json` (boolean, optional): Return clean JSON result for agent parsing.
