package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Touch Demo")

	btn := widget.NewButton("Touch Me!", func() {
		println("Button touched!")
	})

	w.SetContent(btn)
	w.ShowAndRun()
}
