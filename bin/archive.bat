@echo off
setlocal

set CONFIG_FILE=%1
set TARGET_YM=%2

cd /d %~dp0

if not exist diary-generator.exe (
    echo "diary-generator.exe" is not found.
    exit /b 1
)

diary-generator.exe --config %CONFIG_FILE% archive --target-ym %TARGET_YM%
