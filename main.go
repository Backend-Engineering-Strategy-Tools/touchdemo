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
	w.SetFullScreen(true)

	content := container.NewWithoutLayout()
	w.SetContent(content)

	score := 0
	scoreText := canvas.NewText("Score: 0", color.White)
	scoreText.TextSize = 24
	scoreText.Move(fyne.NewPos(20, 20))
	content.Add(scoreText)

	rand.Seed(time.Now().UnixNano())

	// Spawn 1-5 circles every second
	go func() {
		for range time.Tick(time.Second) {
			numCircles := 1 + rand.Intn(5)
			for i := 0; i < numCircles; i++ {
				// Capture coordinates and radius locally
				x := float32(rand.Intn(800))
				y := float32(rand.Intn(600))
				radius := float32(20 + rand.Intn(15))

				circle := canvas.NewCircle(color.RGBA{
					R: uint8(rand.Intn(256)),
					G: uint8(rand.Intn(256)),
					B: uint8(rand.Intn(256)),
					A: 255,
				})
				circle.Resize(fyne.NewSize(radius*2, radius*2))
				circle.Move(fyne.NewPos(x-radius, y-radius))

				// Wrap circle in a transparent button to get Tapped behavior
				btn := widget.NewButton("", func(cx, cy, r float32, c *canvas.Circle, b *widget.Button) func() {
					return func() {
						score += 3
						content.Remove(c)
						content.Remove(b)
						scoreText.Text = fmt.Sprintf("Score: %d", score)
						scoreText.Refresh()
					}
				}(x, y, radius, circle, nil)) // btn passed later

				btn.Resize(fyne.NewSize(radius*2, radius*2))
				btn.Move(fyne.NewPos(x-radius, y-radius))
				btn.Importance = widget.LowImportance // hide button visuals

				content.Add(circle)
				content.Add(btn)
				canvas.Refresh(circle)

				// Animate growth then remove
				go func(c *canvas.Circle, b *widget.Button, r, px, py float32) {
					time.Sleep(time.Second)
					r *= 2
					c.Resize(fyne.NewSize(r*2, r*2))
					c.Move(fyne.NewPos(px-r, py-r))
					b.Resize(fyne.NewSize(r*2, r*2))
					b.Move(fyne.NewPos(px-r, py-r))
					canvas.Refresh(c)
					time.Sleep(time.Second)
					content.Remove(c)
					content.Remove(b)
				}(circle, btn, radius, x, y)
			}
		}
	}()

	w.ShowAndRun()
}
