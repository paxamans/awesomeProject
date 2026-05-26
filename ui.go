package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// cachedApps holds the most recently fetched autostart application list.
// It is refreshed only on startup and after mutations (add/delete/rename).
var cachedApps []string

// appTable is the main table widget showing autostart entries.
var appTable *widget.Table

// refreshApps re-fetches the autostart list and refreshes the table.
func refreshApps() {
	cachedApps = getAutostartApps()
	if appTable != nil {
		appTable.Refresh()
	}
}

// buildUI constructs the full application UI and returns the root container.
func buildUI(w fyne.Window) fyne.CanvasObject {
	// Initial load of the app list.
	cachedApps = getAutostartApps()

	nameLabel := widget.NewLabel("My Autostart Apps:")

	appTable = widget.NewTable(
		// Size callback: uses cached list length.
		func() (int, int) {
			return len(cachedApps), 1
		},
		// Create callback: template cell.
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), container.NewHBox())
		},
		// Update callback: populate each cell from cache.
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			if tci.Row >= len(cachedApps) {
				return
			}
			c := co.(*fyne.Container)
			appEntry := cachedApps[tci.Row]
			appName := strings.Split(appEntry, ": ")[0]

			c.Objects = []fyne.CanvasObject{
				widget.NewLabel(appEntry),
				container.NewHBox(
					// Delete button.
					widget.NewButtonWithIcon("", resourceDeleteIconPng, func() {
						err := DeleteAppFromAutostart(appName)
						if err != nil {
							dialog.ShowError(err, w)
							return
						}
						refreshApps()
					}),
					// Rename button — uses NewForm + Entry instead of
					// the non-existent NewEntryDialog.
					widget.NewButtonWithIcon("", resourceEditIconPng, func() {
						nameEntry := widget.NewEntry()
						nameEntry.SetPlaceHolder("New name")
						formItems := []*widget.FormItem{
							widget.NewFormItem("Name", nameEntry),
						}
						dlg := dialog.NewForm("Rename App", "Rename", "Cancel", formItems, func(ok bool) {
							if !ok || nameEntry.Text == "" {
								return
							}
							err := RenameAppInAutostart(appName, nameEntry.Text)
							if err != nil {
								dialog.ShowError(err, w)
								return
							}
							refreshApps()
						}, w)
						dlg.Resize(fyne.NewSize(400, 200))
						dlg.Show()
					}),
				),
			}
		},
	)

	addButton := widget.NewButton("Add App", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}

			appPath := reader.URI().String()

			err = AddAppToAutostart(appPath)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			refreshApps()
		}, w)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".exe"}))
		fileDialog.Show()
	})

	refreshButton := widget.NewButton("Refresh", func() {
		refreshApps()
	})

	return container.NewBorder(
		nameLabel,
		container.NewHBox(addButton, refreshButton),
		nil,
		nil,
		container.NewMax(appTable),
	)
}
