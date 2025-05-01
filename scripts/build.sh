#!/usr/bin/env sh

go mod download && go mod verify
go build -o ./bin/animesaturn-downloader.exe ./cmd/animesaturn-downloader
cd ./cmd/animesaturn-downloader-gui
wails build
cd ../..
cp "./build/bin/animesaturndownloader" "./bin/animesaturn-downlaoder-gui"
