# diary-generator

## usage

```bash
# init

diary-generator --config ./diary/cofig init

# archive

diary-generator --config ./diary/cofig init archive --target-ym 2023-01
```

## config

`config.yaml` を作成し、以下のように設定する。

```yaml
name: diary
baseDirectory: base-directory-full-path-here
templateFile: template-file-full-path-here
enabledArchiveSummary: true
```

## Windowsのタスクスケジューラーに追加する例

**TODO UPDATE**

`F:\Dropbox\memo` 以下に当リポジトリの `bin` ディレクトリをコピーし、その中に `diary-generator.exe` がある前提の例。

```powershell
schtasks /create /tn "Init Diary" /tr "F:\Dropbox\memo\bin\init.bat diary" /sc daily /st 07:00
```

## download

[Releases · stoneream/diary-generator](https://github.com/stoneream/diary-generator/releases)
