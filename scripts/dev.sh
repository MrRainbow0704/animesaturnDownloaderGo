#!/usr/bin/env sh

VERSION=`cat ./version.txt` + "-dev"
cd ./cmd/animesaturn-downloader-gui
wails dev -tags webkit2_41 -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=${VERSION}'"
cd ../..
