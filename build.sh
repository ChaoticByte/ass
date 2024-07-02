#!/usr/bin/env bash

ver="$(git rev-parse --short HEAD)"

if ! git diff-index --quiet HEAD
then
    ver="${ver}-dev"
fi

go build -ldflags="-X 'main.Version=${ver}'" -o ass_v${ver}
