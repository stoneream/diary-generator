@echo off
setlocal

set TARGET=%1
set STARTSWITH=%2

cd /d %~dp0

if not exist diary-generator.exe (
    echo "diary-generator.exe" is not found.
    exit /b 1
)

diary-generator.exe archive --base-directory-path ..\%TARGET% --starts-with %STARTSWITH%
