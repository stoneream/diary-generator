#!/bin/bash

CONFIG_FILE=$1
TARGET_YM=$2

cd $(dirname $0)

if [ ! -e ./diary-generator ]; then
  echo "diary-generator is not found."
  exit 1
fi

diary-generator --config $CONFIG_FILE archive --target-ym $TARGET_YM
