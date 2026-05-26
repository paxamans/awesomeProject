package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"golang.org/x/sys/windows/registry"
)

// startupFolderPath returns the current user's Windows Startup folder path.
func startupFolderPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(currentUser.HomeDir, "AppData", "Roaming",
		"Microsoft", "Windows", "Start Menu", "Programs", "Startup"), nil
}

// AddAppToAutostart adds the given application to autostart by creating
// a shortcut (.lnk) in the Windows Startup folder.
// appPath must be a native filesystem path (not a file:// URI).
func AddAppToAutostart(appPath string) error {
	// Strip file:// URI scheme if present (Fyne file dialogs return URIs).
	appPath = strings.TrimPrefix(appPath, "file://")
	// On Windows the URI may look like file:///C:/... — trim the leading slash
	// that precedes the drive letter.
	if len(appPath) >= 3 && appPath[0] == '/' && appPath[2] == ':' {
		appPath = appPath[1:]
	}

	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	wshShell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("create WScript.Shell: %w", err)
	}
	defer wshShell.Release()

	wshShellDisp, err := wshShell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("QueryInterface: %w", err)
	}
	defer wshShellDisp.Release()

	startupDir, err := startupFolderPath()
	if err != nil {
		return fmt.Errorf("startup folder: %w", err)
	}

	appName := strings.TrimSuffix(filepath.Base(appPath), filepath.Ext(appPath))
	shortcutPath := filepath.Join(startupDir, appName+".lnk")

	wshShortcut, err := oleutil.CallMethod(wshShellDisp, "CreateShortcut", shortcutPath)
	if err != nil {
		return fmt.Errorf("CreateShortcut: %w", err)
	}
	defer wshShortcut.Clear()

	if _, err = oleutil.PutProperty(wshShortcut.ToIDispatch(), "TargetPath", appPath); err != nil {
		return fmt.Errorf("set TargetPath: %w", err)
	}
	if _, err = oleutil.CallMethod(wshShortcut.ToIDispatch(), "Save"); err != nil {
		return fmt.Errorf("save shortcut: %w", err)
	}

	return nil
}

// DeleteAppFromAutostart removes the named application from both the
// registry Run key and the Startup folder. Both removals are attempted
// regardless of individual failures; errors are collected.
func DeleteAppFromAutostart(appName string) error {
	var errs []string

	// Registry removal.
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err == nil {
		if delErr := k.DeleteValue(appName); delErr != nil && delErr != registry.ErrNotExist {
			errs = append(errs, fmt.Sprintf("registry: %v", delErr))
		}
		k.Close()
	}

	// Startup folder removal.
	startupDir, err := startupFolderPath()
	if err != nil {
		errs = append(errs, fmt.Sprintf("startup path: %v", err))
	} else {
		lnkPath := filepath.Join(startupDir, appName+".lnk")
		if rmErr := os.Remove(lnkPath); rmErr != nil && !os.IsNotExist(rmErr) {
			errs = append(errs, fmt.Sprintf("startup folder: %v", rmErr))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("delete autostart: %s", strings.Join(errs, "; "))
	}
	return nil
}

// RenameAppInAutostart renames the autostart entry in both the registry
// and the Startup folder. Both operations are attempted regardless of
// individual failures.
func RenameAppInAutostart(oldAppName, newAppName string) error {
	var errs []string

	// Registry rename.
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err == nil {
		oldAppPath, _, getErr := k.GetStringValue(oldAppName)
		if getErr == nil {
			if delErr := k.DeleteValue(oldAppName); delErr != nil {
				errs = append(errs, fmt.Sprintf("registry delete old: %v", delErr))
			} else if setErr := k.SetStringValue(newAppName, oldAppPath); setErr != nil {
				errs = append(errs, fmt.Sprintf("registry set new: %v", setErr))
			}
		}
		k.Close()
	}

	// Startup folder rename.
	startupDir, err := startupFolderPath()
	if err != nil {
		errs = append(errs, fmt.Sprintf("startup path: %v", err))
	} else {
		oldPath := filepath.Join(startupDir, oldAppName+".lnk")
		newPath := filepath.Join(startupDir, newAppName+".lnk")
		if renErr := os.Rename(oldPath, newPath); renErr != nil && !os.IsNotExist(renErr) {
			errs = append(errs, fmt.Sprintf("startup folder: %v", renErr))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("rename autostart: %s", strings.Join(errs, "; "))
	}
	return nil
}

// getAutostartApps returns a deduplicated list of autostart application
// names by reading both the Startup folder and the registry Run key.
func getAutostartApps() []string {
	seen := make(map[string]struct{})
	var apps []string

	addUnique := func(name string) {
		if _, exists := seen[name]; !exists {
			seen[name] = struct{}{}
			apps = append(apps, name)
		}
	}

	// --- Startup folder ---
	startupDir, err := startupFolderPath()
	if err == nil {
		entries, readErr := os.ReadDir(startupDir)
		if readErr == nil {
			for _, entry := range entries {
				if !entry.IsDir() && filepath.Ext(entry.Name()) == ".lnk" {
					name := strings.TrimSuffix(entry.Name(), ".lnk")
					addUnique(name)
				}
			}
		}
	}

	// --- Registry ---
	k, err := registry.OpenKey(registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err == nil {
		names, _ := k.ReadValueNames(0)
		for _, name := range names {
			addUnique(name)
		}
		k.Close()
	}

	return apps
}
