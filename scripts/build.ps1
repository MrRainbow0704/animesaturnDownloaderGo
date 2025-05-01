go mod download && go mod verify
go build -o ./bin/animesaturn-downloader.exe ./cmd/animesaturn-downloader
cd ./cmd/animesaturn-downloader-gui
wails build
cd ../..
Copy-Item "./build/bin/animnesaturndownloader.exe" -Destination "./bin/animesaturn-downloader-gui.exe" -Force