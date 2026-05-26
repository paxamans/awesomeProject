package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Awesome Autostart Manager")
	w.SetIcon(resourceAwesomeLogoPng)
	w.Resize(fyne.NewSize(700, 600))

	w.SetContent(buildUI(w))
	w.ShowAndRun()
}
