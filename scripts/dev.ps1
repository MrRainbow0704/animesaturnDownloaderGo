$VERSION = (Get-Content ./version.txt -Raw) + "-dev"
cd ./cmd/animesaturn-downloader-gui
go tool wails dev -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=$VERSION'"
cd ../..
