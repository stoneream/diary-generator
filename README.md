# diary-generator

## usage

```bash
# init

diary-generator init --base-directory ./diary --template-path bin/template/diary.md 

# archive

diary-generator archive --base-directory ./diary --starts-with 2023-01
```

## Windowsのタスクスケジューラーに追加する例

`F:\Dropbox\memo` 以下に当リポジトリの `bin` ディレクトリをコピーし、その中に `diary-generator.exe` がある前提の例。

```powershell
schtasks /create /tn "Init Diary" /tr "F:\Dropbox\memo\bin\init.bat diary" /sc daily /st 07:00
```

## download

[Releases · stoneream/diary-generator](https://github.com/stoneream/diary-generator/releases)
