package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Touch Game")
	w.Resize(fyne.NewSize(800, 600))

	content := container.NewWithoutLayout()
	w.SetContent(content)

	// Dark background
	bg := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	bg.Resize(fyne.NewSize(800, 600))
	content.Add(bg)

	// Score
	score := 0
	scoreText := canvas.NewText("Score: 0", color.White)
	scoreText.TextSize = 24
	scoreText.Move(fyne.NewPos(20, 20))
	content.Add(scoreText)

	go func() {
		for range time.Tick(time.Second) {
			num := 1 + rand.Intn(3)
			for i := 0; i < num; i++ {
				radius := float32(20 + rand.Intn(30))
				x := float32(rand.Intn(800))
				y := float32(rand.Intn(600))

				// Circle
				circle := canvas.NewCircle(color.RGBA{
					R: uint8(50 + rand.Intn(206)),
					G: uint8(50 + rand.Intn(206)),
					B: uint8(50 + rand.Intn(206)),
					A: 255,
				})
				circle.Resize(fyne.NewSize(radius*2, radius*2))
				circle.Move(fyne.NewPos(x-radius, y-radius))

				// Transparent button for tap
				btn := widget.NewButton("", nil)
				btn.Resize(fyne.NewSize(radius*2, radius*2))
				btn.Move(fyne.NewPos(x-radius, y-radius))
				btn.Importance = widget.LowImportance

				// Tap behavior
				btn.OnTapped = func() {
					score += 1
					fyne.Do(func() {
						content.Remove(circle)
						content.Remove(btn)
						scoreText.Text = fmt.Sprintf("Score: %d", score)
						scoreText.Refresh()
					})
				}

				// Add to container on UI thread
				fyne.Do(func() {
					content.Add(circle)
					content.Add(btn)
					circle.Refresh()
					content.Refresh()
				})

				// Animate growth
				go func(c *canvas.Circle, b *widget.Button, px, py, r float32) {
					time.Sleep(time.Second)
					fyne.Do(func() {
						r *= 2
						c.Resize(fyne.NewSize(r*2, r*2))
						c.Move(fyne.NewPos(px-r, py-r))
						b.Resize(fyne.NewSize(r*2, r*2))
						b.Move(fyne.NewPos(px-r, py-r))
						c.Refresh()
						content.Refresh()
					})
					time.Sleep(time.Second)
					fyne.Do(func() {
						content.Remove(c)
						content.Remove(b)
					})
				}(circle, btn, x, y, radius)
			}
		}
	}()

	w.ShowAndRun()
}
