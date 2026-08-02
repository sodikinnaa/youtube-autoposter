$ErrorActionPreference = "Stop"

$Repo = "sodikinnaa/youtube-autoposter"
$BinaryName = "youtube-autoposter.exe"
$SkillFile = "SKILL.md"
$Target = "youtube-autoposter-windows-amd64.exe"

$DownloadUrl = "https://github.com/$Repo/releases/latest/download/$Target"
$RawSkillUrl = "https://raw.githubusercontent.com/$Repo/main/SKILL.md"

if (Test-Path -Path $BinaryName) {
    Write-Host "🔄 Meng-update binary '$BinaryName' yang sudah ada ke versi terbaru..." -ForegroundColor Cyan
} else {
    Write-Host "🚀 Mendownload YouTube Auto-Poster untuk Windows (PowerShell)..." -ForegroundColor Cyan
}

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile "$BinaryName.tmp"
} catch {
    Write-Host "⚠️ Download rilis spesifik gagal. Mengambil fallback..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri "https://github.com/$Repo/releases/latest/download/youtube-autoposter-windows-amd64.exe" -OutFile "$BinaryName.tmp"
}

Move-Item -Path "$BinaryName.tmp" -Destination $BinaryName -Force

try {
    Write-Host "📄 Mendownload file panduan AI Agent (SKILL.md)..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $RawSkillUrl -OutFile $SkillFile
} catch {
    Write-Host "⚠️ Gagal mendownload SKILL.md, mengabaikan..." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=======================================================" -ForegroundColor Green
Write-Host "✅ SUKSES INSTALL / UPDATE YOUTUBE AUTO-POSTER (WINDOWS)!" -ForegroundColor Green
Write-Host "=======================================================" -ForegroundColor Green
Write-Host "📁 Binary Executable : .\$BinaryName"
Write-Host "🤖 Panduan AI Agent  : .\$SkillFile"
Write-Host "🌐 Web UI Dashboard  : .\$BinaryName -web"
Write-Host "======================================================="
Write-Host "💡 Fitur Port Otomatis: Jika port sibuk, server akan memilih port bebas/acak secara otomatis."
Write-Host ""
Write-Host "Jalankan Web Studio: .\$BinaryName -web"
Write-Host "Jalankan Interactive CLI: .\$BinaryName"

if ($args[0] -eq "-run" -or $args[0] -eq "--run" -or $args[0] -eq "-web") {
    Write-Host ""
    Write-Host "🚀 Menjalankan Web UI Studio secara otomatis..." -ForegroundColor Cyan
    Start-Process -FilePath ".\$BinaryName" -ArgumentList "-web" -NoNewWindow
}
