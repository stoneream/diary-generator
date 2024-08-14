#!/bin/bash

TARGET=$1

cd $(dirname $0)

if [ ! -e ./diary-generator ]; then
  echo "diary-generator is not found."
  exit 1
fi

./diary-generator init --base-directory ../$TARGET --template-file ./template/$TARGET.md
