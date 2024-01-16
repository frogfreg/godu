#!/bin/bash
rm -rf ./build

os_arch=("linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64")

for target in "${os_arch[@]}"; do
    IFS='/' read -r -a parts <<< "$target"
    os="${parts[0]}"
    arch="${parts[1]}"

    if [ "$os" == "linux" ]; then
        GOOS="$os" GOARCH="$arch" go build -o "./build/godu_${os}_${arch}"
        gzip "./build/godu_${os}_${arch}"
    fi

    if [ "$os" == "darwin" ]; then
        GOOS="$os" GOARCH="$arch" go build -o "./build/godu_macos_${arch}"
        gzip "./build/godu_macos_${arch}"
    fi

    if [ "$os" == "windows" ]; then
        GOOS="$os" GOARCH="$arch" go build -o "./build/godu_${os}_${arch}.exe"
    fi

done