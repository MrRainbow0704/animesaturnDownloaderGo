$VERSION = (Get-Content ./version.txt -Raw) + "-dev"
cd ./cmd/animesaturn-downloader-gui
wails dev -ldflags="-X 'github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version.version=$VERSION'"
cd ../..
