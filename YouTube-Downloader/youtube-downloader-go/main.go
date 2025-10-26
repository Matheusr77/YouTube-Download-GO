package main

import (
	"context"
	"embed"

	"youtube-downloader-go/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	app := backend.NewApp()

	err := wails.Run(&options.App{
		Title:  "YouTube Downloader",
		Width:  600,
		Height: 400,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			app.Startup(ctx)
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Erro ao iniciar o app:", err.Error())
	}
}
