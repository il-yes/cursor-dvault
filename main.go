package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	vaults_storage "vault-app/internal/vault/infrastructure/storage"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS
var avatarStore *vaults_storage.AvatarStore
 
func main() {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("❌ Error loading .env file:", errEnv)
	}
	privateKey := os.Getenv("STELLAR_PRIVATE_KEY")
	if privateKey == "" {
		fmt.Println("❌ STELLAR_PRIVATE_KEY is empty")
	}
    rootPath := "./vault" // your vault root
    avatarStore = vaults_storage.NewAvatarStore(rootPath)

	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "ANKHORA",
		Width:  1000,
		Height: 590,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		// Mac: &mac.Options{
		// 	OnUrlOpen: app.OnOpenURL,
		// },
		// SingleInstanceLock: &options.SingleInstanceLock{
		// 	OnSecondInstanceLaunch: func(data options.SecondInstanceData) {
		// 		for _, arg := range data.Args {
		// 			if strings.HasPrefix(arg, "ankhora://") {
		// 				app.OnOpenURL(arg)
		// 			}
		// 		}
		// 	},
		// },

		OnShutdown: func(ctx context.Context) {
			app.Logger.Info("🛑 App shutting down, flushing sessions...")
			app.FlushAllSessions()

			if app.cancel != nil {
				app.cancel()
			}
			app.Logger.Info("👋 Shutdown complete")
		},
		Bind: []interface{}{
			app,
			avatarStore,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}

}
