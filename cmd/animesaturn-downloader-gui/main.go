package main

import (
	"context"
	"flag"

	"github.com/MrRainbow0704/animesaturnDownloaderGo/frontend"
	log "github.com/MrRainbow0704/animesaturnDownloaderGo/internal/logger"
	"github.com/MrRainbow0704/animesaturnDownloaderGo/internal/version"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var appLogger = &log.AppLogger{}

func main() {
	var localConfig bool
	flag.BoolVar(&localConfig, "local-config", false, "force the program to use local config")
	var localCache bool
	flag.BoolVar(&localCache, "local-cache", false, "force the program to use local cache")

	app := &App{}
	opts := &options.App{
		Title:            "Animesaturn Downlaoder " + version.Get(),
		Width:            1280,
		Height:           720,
		AssetServer:      &assetserver.Options{Assets: frontend.Assets},
		BackgroundColour: &options.RGBA{R: 18, G: 22, B: 25, A: 1},
		OnStartup: func(ctx context.Context) {
			flag.Parse()
			app.startup(ctx, localConfig, localCache)
		},
		Bind:    []any{app},
		Logger:  appLogger,
		Windows: &windows.Options{IsZoomControlEnabled: false},
		Mac:     &mac.Options{DisableZoom: true},
	}

	if err := wails.Run(opts); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
