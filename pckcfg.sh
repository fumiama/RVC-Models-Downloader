#!/bin/sh
# Should be run by go generate. DO NOT run it directly.

files=($(find $@ | grep -v .DS_Store | sort))

rm -rf cfg.zip
zip -9 cfg.zip "${files[@]}"
