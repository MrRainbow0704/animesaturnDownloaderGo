package main

import (
	"github.com/MrRainbow0704/animesaturnDownloaderGo/frontend"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	app := &App{}
	opts := &options.App{
		Title:  "Animesaturn Downlaoder",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: frontend.Assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 22, B: 25, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
		Logger: &log.AppLogger{},
	}

	if err := wails.Run(opts); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
