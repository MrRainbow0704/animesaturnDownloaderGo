#!/usr/bin/env sh

VERSION=`cat ./version.txt`
go mod download && go mod verify
go build -o ./bin/animesaturn-downloader -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=${VERSION}'" ./cmd/animesaturn-downloader
cd ./cmd/animesaturn-downloader-gui
wails build -tags webkit2_41 -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=${VERSION}'"
cd ../..
cp ./build/bin/animesaturn-downloader-gui ./bin/animesaturn-downloader-gui
