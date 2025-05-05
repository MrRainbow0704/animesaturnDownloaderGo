$VERSION = Get-Content ./version.txt -Raw
cd ./cmd/animesaturn-downloader-gui
wails dev -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=$VERSION'"
cd ../..
