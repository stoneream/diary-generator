# diary-generator

## development

```bash
mise trust
mise install
mise run deps
mise run build
mise run test
mise run vet
```

## install

```bash
go install github.com/stoneream/diary-generator/v2@latest
```

## usage

```bash
# ベースディレクトリの作成
mkdir diary
cd diary

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
