# diary-generator

## usage

```bash
# init

java -jar bin/diary-generator.jar init --base-directory ./diary --template-path bin/template/diary.md 

# archive

java -jar bin/diary-generator.jar archive --base-directory ./diary --starts-with 2023-01
```

## Windowsのタスクスケジューラーに追加する例

`F:\Dropbox\memo` 以下に `.jar` と当リポジトリの `bin` を配置していることを前提とする。  

```powershell
schtasks /create /tn "Init Diary" /tr "F:\Dropbox\memo\bin\init.bat diary" /sc daily /st 07:00
```

## download

[Releases · stoneream/diary-generator](https://github.com/stoneream/diary-generator/releases)

## memo

テンプレートファイルとディレクトリの名前は合わせたほうが扱いやすい。

