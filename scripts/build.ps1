$VERSION = Get-Content ./version.txt -Raw
go mod download && go mod verify
go build -o ./bin/animesaturn-downloader.exe  -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=$VERSION'" ./cmd/animesaturn-downloader
cd ./cmd/animesaturn-downloader-gui
wails build -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=$VERSION'"
Remove-Item "./config.json" -Force
Remove-Item "./.cache" -Force -Recurse
cd ../..
Copy-Item "./build/bin/animesaturn-downloader-gui.exe" -Destination "./bin/animesaturn-downloader-gui.exe" -Force
