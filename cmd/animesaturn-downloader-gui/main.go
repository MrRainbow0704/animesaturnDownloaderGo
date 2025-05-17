package main

import (
	"github.com/MrRainbow0704/animesaturnDownloaderGo/frontend"
	_ "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/cache"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var appLogger = &log.AppLogger{}

func main() {
	app := &App{Version: version.Get()}
	opts := &options.App{
		Title:            "Animesaturn Downlaoder " + app.Version,
		Width:            1280,
		Height:           720,
		AssetServer:      &assetserver.Options{Assets: frontend.Assets},
		BackgroundColour: &options.RGBA{R: 18, G: 22, B: 25, A: 1},
		OnStartup:        app.startup,
		Bind:             []any{app},
		Logger:           appLogger,
		DisableResize:    true,
	}

	if err := wails.Run(opts); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
