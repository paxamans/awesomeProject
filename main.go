package main

import (
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var appTable *widget.Table

func main() {
	a := app.New()
	w := a.NewWindow("Awesome")
	w.SetIcon(loadIcon("saves/awesome_logo.png"))

	w.Resize(fyne.NewSize(700, 600))

	nameLabel := widget.NewLabel("My Autostart Apps:")
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
					widget.NewButtonWithIcon("", loadIcon("saves/delete_icon.png"), func() {
						appName := strings.Split(getAutostartApps()[tci.Row], ": ")[0]
						err := DeleteAppFromAutostart(appName)
						if err != nil {
							dialog.ShowError(err, w)
							return
						}
						appTable.Refresh()
					}),
					widget.NewButtonWithIcon("", loadIcon("saves/edit_icon.png"), func() {
						oldAppName := strings.Split(getAutostartApps()[tci.Row], ": ")[0]
						entryDialog := dialog.NewEntryDialog("New name", "Enter new name for the app", func(newName string) {
							if newName != "" {
								err := RenameAppInAutostart(oldAppName, newName)
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

			err = AddAppToAutostart(appPath)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			appTable.Refresh()
		}, w)
		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".exe"}))
		fileDialog.Show()
	})

	refreshButton := widget.NewButton("Refresh", func() {
		appTable.Refresh()
	})
	content := container.NewBorder(
		nameLabel,
		container.NewHBox(addButton, refreshButton),
		nil,
		nil,
		container.NewMax(appTable),
	)

	w.SetContent(content)
	w.ShowAndRun()
}
func loadIcon(filename string) fyne.Resource {
	icon, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return fyne.NewStaticResource(filename, icon)
}

// AddAppToAutostart adds the given application to autostart on Windows
func AddAppToAutostart(appPath string) error {
	// Initialize the COM library
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// Create the WshShell object
	wshShell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer wshShell.Release()

	// Convert to IDispatch
	wshShellDisp, err := wshShell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshShellDisp.Release()

	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	// Get the path to the Startup folder
	startupPath := filepath.Join(currentUser.HomeDir, "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup")

	// Get the name of the application from the path to the executable file
	appName := strings.TrimSuffix(filepath.Base(appPath), filepath.Ext(appPath))

	// Create the shortcut
	shortcutPath := filepath.Join(startupPath, appName+".lnk")
	wshShortcut, err := oleutil.CallMethod(wshShellDisp, "CreateShortcut", shortcutPath)
	if err != nil {
		return err
	}
	defer wshShortcut.Clear()

	// Set the target path of the shortcut
	_, err = oleutil.PutProperty(wshShortcut.ToIDispatch(), "TargetPath", appPath)
	if err != nil {
		return err
	}

	// Save the shortcut
	_, err = oleutil.CallMethod(wshShortcut.ToIDispatch(), "Save")
	if err != nil {
		return err
	}

	return nil
}

// DeleteAppFromAutostart deletes the given application from autostart on Windows
func DeleteAppFromAutostart(appName string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	// Get the path of the app
	appPath, _, err := k.GetStringValue(appName)
	if err != nil {
		return err
	}

	// Check if the file exists
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		return fmt.Errorf("the file does not exist: %s", appPath)
	}

	err = k.DeleteValue(appName)
	if err != nil {
		return err
	}

	return nil
}

// RenameAppInAutostart renames the given application in autostart on Windows
func RenameAppInAutostart(oldAppName, newAppName string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	// Get the path of the old app
	oldAppPath, _, err := k.GetStringValue(oldAppName)
	if err != nil {
		return err
	}

	// Check if the file exists
	if _, err := os.Stat(oldAppPath); os.IsNotExist(err) {
		return fmt.Errorf("the file does not exist: %s", oldAppPath)
	}

	err = k.DeleteValue(oldAppName)
	if err != nil {
		return err
	}

	err = k.SetStringValue(newAppName, oldAppPath)
	if err != nil {
		return err
	}

	return nil
}

// getAutostartApps returns a list of applications in the autostart list
func getAutostartApps() []string {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		log.Println(err)
		return nil
	}

	// Get the path to the Startup folder
	startupPath := filepath.Join(currentUser.HomeDir, "AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup")

	// Read the files in the Startup folder
	files, err := ioutil.ReadDir(startupPath)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Create a slice to hold the names of the autostart applications
	apps := make([]string, 0, len(files))

	// Loop over the files and add their names to the slice
	for _, file := range files {
		// Only consider .lnk files (shortcuts)
		if filepath.Ext(file.Name()) == ".lnk" {
			// Remove the .lnk extension from the name
			name := strings.TrimSuffix(file.Name(), ".lnk")
			apps = append(apps, name)
		}
	}

	// Open the Run registry key
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		log.Println(err)
		return apps
	}
	defer k.Close()

	// Get the names of the autostart entries
	names, err := k.ReadValueNames(0)
	if err != nil {
		log.Println(err)
		return apps
	}

	// Append the names of the autostart entries to the apps slice
	apps = append(apps, names...)

	return apps
}
