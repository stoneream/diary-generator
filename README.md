# diary-generator

## install

```bash
go install github.com/stoneream/diary-generator/v2@latest
```

## usage

```bash
# ベースディレクトリの作成
mkdir diary

# テンプレートファイルの作成
touch template.md

# init

diary-generator init

# 2025/01/01 16:00:00 initialized successfully at diary_2025-02-20.md

# archive

diary-generator archive --target-ym 2023-01
```

## download

[Releases · stoneream/diary-generator](https://github.com/stoneream/diary-generator/releases)
