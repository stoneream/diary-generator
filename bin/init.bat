@echo off
setlocal

set TARGET=%1

cd /d %~dp0

if not exist diary-generator.exe (
    echo "diary-generator.exe" is not found.
    exit /b 1
)

diary-generator.exe init --base-directory ..\%TARGET% --template-path .\template\%TARGET%.md
