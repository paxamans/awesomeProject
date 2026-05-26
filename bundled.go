package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed saves/awesome_logo.png
var awesomeLogoBytes []byte

//go:embed saves/delete_icon.png
var deleteIconBytes []byte

//go:embed saves/edit_icon.png
var editIconBytes []byte

var (
	resourceAwesomeLogoPng = fyne.NewStaticResource("awesome_logo.png", awesomeLogoBytes)
	resourceDeleteIconPng  = fyne.NewStaticResource("delete_icon.png", deleteIconBytes)
	resourceEditIconPng    = fyne.NewStaticResource("edit_icon.png", editIconBytes)
)
