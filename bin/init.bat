@echo off
setlocal

set TARGET=%1

cd /d %~dp0

java -jar .\diary-generator.jar init --base-directory-path ..\%TARGET% --template-path .\template\%TARGET%.md
