@echo off
setlocal

set REPO=sodikinnaa/youtube-autoposter
set BINARY_NAME=youtube-autoposter.exe
set SKILL_FILE=SKILL.md
set TARGET=youtube-autoposter-windows-amd64.exe

set DOWNLOAD_URL=https://github.com/%REPO%/releases/latest/download/%TARGET%
set RAW_SKILL_URL=https://raw.githubusercontent.com/%REPO%/main/SKILL.md

if exist "%BINARY_NAME%" (
    echo 🔄 Meng-update binary '%BINARY_NAME%' yang sudah ada ke versi terbaru...
) else (
    echo 🚀 Mendownload YouTube Auto-Poster untuk Windows CMD...
)

curl -sSL "%DOWNLOAD_URL%" -o "%BINARY_NAME%.tmp"
if errorlevel 1 (
    echo ⚠️ Download rilis spesifik gagal. Mengambil fallback...
    curl -sSL "https://github.com/%REPO%/releases/latest/download/youtube-autoposter-windows-amd64.exe" -o "%BINARY_NAME%.tmp"
)

move /Y "%BINARY_NAME%.tmp" "%BINARY_NAME%" >nul

echo 📄 Mendownload file panduan AI Agent (SKILL.md)...
curl -sSL "%RAW_SKILL_URL%" -o "%SKILL_FILE%"

if not exist "sample_video.mp4" (
    echo 📹 Mendownload sample video publik ^(sample_video.mp4^) untuk uji coba otomatis...
    curl -sSL "https://raw.githubusercontent.com/intel-iot-devkit/sample-videos/master/person-bicycle-car-detection.mp4" -o "sample_video.mp4"
)

echo.
echo =======================================================
echo ✅ SUKSES INSTALL / UPDATE YOUTUBE AUTO-POSTER (WINDOWS CMD)!
echo =======================================================
echo 📁 Binary Executable : .\%BINARY_NAME%
echo 🤖 Panduan AI Agent  : .\%SKILL_FILE%
echo 🌐 Web UI Dashboard  : .\%BINARY_NAME% -web
echo =======================================================
echo 💡 Fitur Port Otomatis: Jika port sibuk, server akan memilih port bebas/acak secara otomatis.
echo.
echo Jalankan Web Studio: .\%BINARY_NAME% -web
echo Jalankan Interactive CLI: .\%BINARY_NAME%

if "%1"=="-run" (
    echo.
    echo 🚀 Menjalankan Web UI Studio secara otomatis...
    .\%BINARY_NAME% -web
) else if "%1"=="--run" (
    echo.
    echo 🚀 Menjalankan Web UI Studio secara otomatis...
    .\%BINARY_NAME% -web
) else if "%1"=="-web" (
    echo.
    echo 🚀 Menjalankan Web UI Studio secara otomatis...
    .\%BINARY_NAME% -web
)
