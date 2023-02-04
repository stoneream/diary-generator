@echo off
setlocal

set TARGET=%1
set STARTSWITH=%2

cd /d %~dp0

java -jar .\diary-generator.jar archive --base-directory-path ..\%TARGET% --starts-with %STARTSWITH%
