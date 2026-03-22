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

	// Top banner
	health := 30
	score := 0
	gameOver := false

	banner := canvas.NewRectangle(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	banner.Resize(fyne.NewSize(800, 50))
	banner.Move(fyne.NewPos(0, 0))
	content.Add(banner)

	scoreText := canvas.NewText(fmt.Sprintf("Score: %d", score), color.White)
	scoreText.TextSize = 24
	scoreText.Move(fyne.NewPos(20, 10))
	content.Add(scoreText)

	healthText := canvas.NewText(fmt.Sprintf("Health: %d", health), color.White)
	healthText.TextSize = 24
	healthText.Move(fyne.NewPos(680, 10))
	content.Add(healthText)

	// Floating text helper
	showFloatingText := func(x, y float32, text string, col color.Color) {
		ft := canvas.NewText(text, col)
		ft.TextStyle.Bold = true
		ft.TextSize = 20
		ft.Move(fyne.NewPos(x-10, y-30))
		content.Add(ft)
		ft.Refresh()
		go func() {
			for i := 0; i < 20; i++ {
				time.Sleep(25 * time.Millisecond)
				ft.Move(fyne.NewPos(x-10, y-30-float32(i)))
				ft.Refresh()
			}
			fyne.Do(func() { content.Remove(ft) })
		}()
	}

	// New Game function
	startGame := func() {
		score = 0
		health = 30
		gameOver = false
		scoreText.Text = fmt.Sprintf("Score: %d", score)
		scoreText.Refresh()
		healthText.Text = fmt.Sprintf("Health: %d", health)
		healthText.Refresh()
	}

	go func() {
		for range time.Tick(time.Second) {
			if gameOver {
				continue
			}

			num := 1 + rand.Intn(3)
			for i := 0; i < num; i++ {
				radius := float32(20 + rand.Intn(20))
				x := float32(rand.Intn(800))
				y := float32(rand.Intn(600-50) + 50) // avoid top banner

				circle := canvas.NewCircle(color.RGBA{
					R: uint8(50 + rand.Intn(206)),
					G: uint8(50 + rand.Intn(206)),
					B: uint8(50 + rand.Intn(206)),
					A: 255,
				})
				circle.Resize(fyne.NewSize(radius*2, radius*2))
				circle.Move(fyne.NewPos(x-radius, y-radius))

				btn := widget.NewButton("", nil)
				btn.Resize(fyne.NewSize(radius*2, radius*2))
				btn.Move(fyne.NewPos(x-radius, y-radius))
				btn.Importance = widget.LowImportance

				stagePoints := []int{3, 2, 1, 0, 0, 0, 0}    // points per stage
				stageHealth := []int{0, 0, 0, 0, -1, -2, -3} // health per stage

				currentStage := 0

				btn.OnTapped = func() {
					if currentStage < len(stagePoints) {
						points := stagePoints[currentStage]
						if points > 0 {
							score += points
							showFloatingText(x, y, fmt.Sprintf("+%d", points), color.RGBA{R: 212, G: 175, B: 55, A: 255}) // gold
						}
					}
					fyne.Do(func() {
						content.Remove(circle)
						content.Remove(btn)
						scoreText.Text = fmt.Sprintf("Score: %d", score)
						scoreText.Refresh()
					})
					currentStage = len(stagePoints)
				}

				fyne.Do(func() {
					content.Add(circle)
					content.Add(btn)
					circle.Refresh()
					content.Refresh()
				})

				go func(c *canvas.Circle, b *widget.Button, px, py, r float32) {
					stages := []float32{r, r * 1.3, r * 1.6, r * 1.9, r * 2.2, r * 2.5, r * 2.8}
					for i, s := range stages {
						time.Sleep(time.Second)
						currentStage = i
						fyne.Do(func() {
							c.Resize(fyne.NewSize(s*2, s*2))
							c.Move(fyne.NewPos(px-s, py-s))
							b.Resize(fyne.NewSize(s*2, s*2))
							b.Move(fyne.NewPos(px-s, py-s))
							c.Refresh()
							content.Refresh()
						})

						// Health penalty stages
						if i >= 4 && currentStage == i {
							fyne.Do(func() {
								if currentStage < len(stagePoints) {
									hpChange := stageHealth[i]
									health += hpChange
									if health < 0 {
										health = 0
									}
									if hpChange < 0 {
										showFloatingText(px, py, fmt.Sprintf("%d", hpChange), color.RGBA{R: 255, G: 50, B: 50, A: 255}) // red
									}
									healthText.Text = fmt.Sprintf("Health: %d", health)
									healthText.Refresh()
									if health == 0 {
										gameOver = true
										// Show game over banner and new game button
										go func() {
											fyne.Do(func() {
												gameOverText := canvas.NewText("GAME OVER", color.White)
												gameOverText.TextSize = 48
												gameOverText.TextStyle.Bold = true
												gameOverText.Move(fyne.NewPos(250, 250))
												content.Add(gameOverText)
												content.Refresh()

												newGameBtn := widget.NewButton("New Game", func() {})

												newGameBtn.OnTapped = func() {
													content.Remove(gameOverText)
													content.Remove(newGameBtn)
													startGame()
												}

												newGameBtn.Resize(fyne.NewSize(200, 50))
												newGameBtn.Move(fyne.NewPos(300, 320))
												content.Add(newGameBtn)
												content.Refresh()
											})
										}()
									}
								}
							})
						}
					}

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
