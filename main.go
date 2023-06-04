package main

import (
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var appTable *widget.Table

func main() {
	a := app.New()
	w := a.NewWindow("Autostart App Manager")

	w.Resize(fyne.NewSize(1300, 600))

	nameLabel := widget.NewLabel("My Autostart Apps")
	appTable = widget.NewTable(
		func() (int, int) {
			return len(getAutostartApps()), 1
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel(""), container.NewHBox())
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			c := co.(*fyne.Container)
			c.Objects = []fyne.CanvasObject{
				widget.NewLabel(getAutostartApps()[tci.Row]),
				container.NewHBox(
					widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
						err := DeleteAppFromAutostart(getAutostartApps()[tci.Row])
						if err != nil {
							dialog.ShowError(err, w)
							return
						}
						appTable.Refresh()
					}),
					widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
						entryDialog := dialog.NewEntryDialog("New name", "Enter new name for the app", func(newName string) {
							if newName != "" {
								err := RenameAppInAutostart(getAutostartApps()[tci.Row], newName)
								if err != nil {
									dialog.ShowError(err, w)
									return
								}
								appTable.Refresh()
							}
						}, w)
						entryDialog.Show()
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
				dialog.ShowInformation("No file selected", "No file selected.", w)
				return
			}

			appPath := reader.URI().String()

			err = AddAppToAutostart("Added app with Yocker", appPath)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			appTable.Refresh()
		}, w)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".exe"}))
		fileDialog.Show()
	})

	content := container.NewBorder(
		nameLabel,
		addButton,
		nil,
		nil,
		container.NewMax(appTable),
	)

	w.SetContent(content)
	w.ShowAndRun()
}

// AddAppToAutostart adds the given application to autostart on Windows
func AddAppToAutostart(appName, appPath string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.SetStringValue(appName, appPath)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAppFromAutostart deletes the given application from autostart on Windows
func DeleteAppFromAutostart(appName string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.DeleteValue(appName)
	if err != nil {
		return err
	}

	return nil
}

// RenameAppInAutostart renames the given application in autostart on Windows
func RenameAppInAutostart(oldAppName, newAppName string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	value, _, err := k.GetStringValue(oldAppName)
	if err != nil {
		return err
	}

	err = k.DeleteValue(oldAppName)
	if err != nil {
		return err
	}

	err = k.SetStringValue(newAppName, value)
	if err != nil {
		return err
	}

	return nil
}

// getAutostartApps returns a list of applications in the autostart list
func getAutostartApps() []string {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer k.Close()

	names, err := k.ReadValueNames(0)
	if err != nil {
		log.Println(err)
		return nil
	}

	apps := make([]string, len(names))
	for i, name := range names {
		value, _, err := k.GetStringValue(name)
		if err != nil {
			log.Println(err)
			continue
		}
		apps[i] = strings.Join([]string{name, value}, ": ")
	}

	return apps
}
