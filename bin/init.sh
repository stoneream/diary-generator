#!/bin/bash

TARGET=$1

cd $(dirname $0)

java -jar ./diary-generator.jar init --base-directory-path ../$TARGET --template-path ./template/$TARGET.md
