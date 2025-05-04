#!/usr/bin/env sh

go mod download && go mod verify
go build -o ./bin/animesaturn-downloader ./cmd/animesaturn-downloader
cd ./cmd/animesaturn-downloader-gui
wails build -tags webkit2_41
cd ../..
cp ./build/bin/animnesaturndownloader ./bin/animesaturn-downloader-gui
