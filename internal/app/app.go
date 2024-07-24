package app

import (
	"db_connector/ui"
	"fyne.io/fyne/v2/app"
	_ "log"
)

func Run() error {
	a := app.New()
	w := ui.NewMainWindow(a)
	w.ShowAndRun()

	return nil
}
