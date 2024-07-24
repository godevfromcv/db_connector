package ui

import (
	_ "database/sql"
	"db_connector/internal/database"
	"fmt"
	"fyne.io/fyne/v2"
	_ "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewMainWindow(a fyne.App) fyne.Window {
	w := a.NewWindow("DB Client")

	hostEntry := widget.NewEntry()
	hostEntry.SetPlaceHolder("Host")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("User")

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Password")

	statusLabel := widget.NewLabel("")

	connectButton := widget.NewButton("Connect", func() {
		host := hostEntry.Text
		port := portEntry.Text
		user := userEntry.Text
		password := passwordEntry.Text

		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=require", host, port, user, password)

		db, err := database.InitDB(connStr)
		if err != nil {
			statusLabel.SetText("Connection failed: " + err.Error())
			return
		}

		defer db.Close()
		statusLabel.SetText("Connection successful")
	})

	form := container.NewVBox(
		widget.NewLabel("Enter Database Connection Details"),
		hostEntry,
		portEntry,
		userEntry,
		passwordEntry,
		connectButton,
		statusLabel,
	)

	w.SetContent(form)
	return w
}
