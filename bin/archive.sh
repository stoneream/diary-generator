#!/bin/bash

TARGET=$1
STARTSWITH=$2

cd $(dirname $0)

if [ ! -e ./diary-generator ]; then
  echo "diary-generator is not found."
  exit 1
fi

diary-generator archive --base-directory ../$TARGET --starts-with $STARTSWITH
