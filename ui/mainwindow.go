package ui

import (
	"database/sql"
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

		statusLabel.SetText("Connection successful")
		dbWindow := NewDatabaseWindow(a, db)
		dbWindow.Show()
		w.Close()
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

func NewDatabaseWindow(a fyne.App, db *sql.DB) fyne.Window {
	w := a.NewWindow("Database Explorer")

	// Entry for SQL query
	queryEntry := widget.NewMultiLineEntry()
	queryEntry.SetPlaceHolder("Enter your SQL query here...")

	// Result label
	resultLabel := widget.NewLabel("")

	executeButton := widget.NewButton("Execute", func() {
		query := queryEntry.Text
		rows, err := database.ExecuteQuery(db, query)
		if err != nil {
			resultLabel.SetText("Error: " + err.Error())
			return
		}
		defer rows.Close()

		resultText := "Query Results:\n"
		columns, err := rows.Columns()
		if err != nil {
			resultLabel.SetText("Error fetching columns: " + err.Error())
			return
		}

		for rows.Next() {
			values := make([]interface{}, len(columns))
			pointers := make([]interface{}, len(values))
			for i := range values {
				pointers[i] = &values[i]
			}

			if err := rows.Scan(pointers...); err != nil {
				resultLabel.SetText("Error: " + err.Error())
				return
			}

			for i, colName := range columns {
				resultText += colName + ": " + fmt.Sprintf("%v", values[i]) + "\n"
			}
			resultText += "-------------------------\n"
		}

		if err := rows.Err(); err != nil {
			resultLabel.SetText("Error: " + err.Error())
			return
		}

		resultLabel.SetText(resultText)
	})

	// Get list of tables
	tables, err := database.GetTables(db)
	if err != nil {
		resultLabel.SetText("Error fetching tables: " + err.Error())
	}
	tableList := widget.NewList(
		func() int {
			return len(tables)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(tables[i])
		},
	)

	// Layout
	w.SetContent(container.NewBorder(
		nil, nil, tableList, executeButton,
		container.NewVBox(queryEntry, resultLabel),
	))

	return w
}
