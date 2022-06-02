package main

import (
	"gobang/asset"
	"gobang/window"
	"os"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/flopp/go-findfont"
)

func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyhl.ttc") {
			os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func main() {
	defer os.Unsetenv("FYNE_FONT")
	app := app.NewWithID("gobang_ai")
	app.SetIcon(theme.FyneLogo())

	wnd := window.DialogWindow(
		"五子棋人机大战",
		asset.ResourceBoardPng,
		asset.ResourceWhitechessPng,
		asset.ResourceBlackchessPng,
		app,
	)
	wnd.ShowAndRun()
}
