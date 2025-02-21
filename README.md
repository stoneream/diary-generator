# diary-generator

## install

```bash
TAG=x.y.z
go install github.com/stoneream/diary-generator@$TAG
```

## usage

```bash
# ベースディレクトリの作成
mkdir diary

# テンプレートファイルの作成
touch template.md

# 初期化

diary-generator init

# アーカイブ

diary-generator archive --target-ym 2023-01

# サマリ

diary-generator summary

cd archive
diary-generator summary --target-prefix diary_
```

## download

[Releases · stoneream/diary-generator](https://github.com/stoneream/diary-generator/releases)
