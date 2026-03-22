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

type Circle struct {
	circle       *canvas.Circle
	button       *widget.Button
	x, y, radius float32
	currentStage int
	removed      bool
	stop         chan struct{} // signal to stop growth
}

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
	banner := canvas.NewRectangle(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	banner.Resize(fyne.NewSize(800, 50))
	banner.Move(fyne.NewPos(0, 0))
	content.Add(banner)

	// Game state
	var (
		score       int
		health      int
		gameOver    bool
		newGameBtn  *widget.Button
		activeCircs []*Circle
		gameActive  bool // new flag to stop old goroutines
	)

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
		ft.TextSize = 20
		ft.TextStyle.Bold = true
		ft.Move(fyne.NewPos(x-10, y-30))
		fyne.Do(func() { content.Add(ft); ft.Refresh() })

		go func() {
			for i := 0; i < 20; i++ {
				time.Sleep(25 * time.Millisecond)
				posY := y - 30 - float32(i)
				fyne.Do(func() {
					ft.Move(fyne.NewPos(x-10, posY))
					ft.Refresh()
				})
			}
			fyne.Do(func() { content.Remove(ft) })
		}()
	}

	// Start a new game
	startGame := func() {
		fyne.Do(func() {
			// Stop old circles
			gameActive = false
			for _, c := range activeCircs {
				if !c.removed {
					content.Remove(c.circle)
					content.Remove(c.button)
					close(c.stop)
					c.removed = true
				}
			}
			activeCircs = []*Circle{}
			gameActive = true

			score = 0
			health = 30
			gameOver = false
			scoreText.Text = fmt.Sprintf("Score: %d", score)
			scoreText.Refresh()
			healthText.Text = fmt.Sprintf("Health: %d", health)
			healthText.Refresh()

			if newGameBtn != nil {
				content.Remove(newGameBtn)
				newGameBtn = nil
			}
		})
	}

	// Spawn circles
	go func() {
		rand.Seed(time.Now().UnixNano())
		for range time.Tick(time.Second) {
			if !gameActive || gameOver {
				continue
			}

			num := 1 + rand.Intn(3)
			for i := 0; i < num; i++ {
				if !gameActive {
					break
				}

				radius := float32(20 + rand.Intn(20))
				x := radius + float32(rand.Intn(800-int(2*radius)))
				y := 50 + radius + float32(rand.Intn(600-50-int(2*radius)))

				c := &Circle{
					circle: canvas.NewCircle(color.RGBA{
						R: uint8(50 + rand.Intn(206)),
						G: uint8(50 + rand.Intn(206)),
						B: uint8(50 + rand.Intn(206)),
						A: 255,
					}),
					x:      x,
					y:      y,
					radius: radius,
					stop:   make(chan struct{}),
				}

				c.circle.Resize(fyne.NewSize(radius*2, radius*2))
				c.circle.Move(fyne.NewPos(x-radius, y-radius))

				c.button = widget.NewButton("", nil)
				c.button.Resize(fyne.NewSize(radius*2, radius*2))
				c.button.Move(fyne.NewPos(x-radius, y-radius))
				c.button.Importance = widget.LowImportance

				stagePoints := []int{3, 2, 1, 0, 0, 0, 0}
				stageHealth := []int{0, 0, 0, 0, -1, -2, -3}

				// Handle click
				c.button.OnTapped = func() {
					if c.removed {
						return
					}
					points := 0
					if c.currentStage < len(stagePoints) {
						points = stagePoints[c.currentStage]
					}
					if points > 0 {
						score += points
						showFloatingText(c.x, c.y, fmt.Sprintf("+%d", points), color.RGBA{R: 212, G: 175, B: 55, A: 255})
					}
					fyne.Do(func() {
						content.Remove(c.circle)
						content.Remove(c.button)
						scoreText.Text = fmt.Sprintf("Score: %d", score)
						scoreText.Refresh()
					})
					c.removed = true
					close(c.stop)
				}

				activeCircs = append(activeCircs, c)

				fyne.Do(func() {
					content.Add(c.circle)
					content.Add(c.button)
					c.circle.Refresh()
					content.Refresh()
				})

				// Growth goroutine
				go func(c *Circle) {
					stages := []float32{c.radius, c.radius * 1.3, c.radius * 1.6, c.radius * 1.9, c.radius * 2.2, c.radius * 2.5, c.radius * 2.8}
					for i, s := range stages {
						select {
						case <-c.stop:
							return
						case <-time.After(time.Second):
						}

						c.currentStage = i
						fyne.Do(func() {
							if c.removed || !gameActive {
								return
							}
							c.circle.Resize(fyne.NewSize(s*2, s*2))
							c.circle.Move(fyne.NewPos(c.x-s, c.y-s))
							c.button.Resize(fyne.NewSize(s*2, s*2))
							c.button.Move(fyne.NewPos(c.x-s, c.y-s))
							c.circle.Refresh()
							c.button.Refresh()

							if i >= 4 && !gameOver {
								hpChange := stageHealth[i]
								if hpChange != 0 {
									health += hpChange
									if health < 0 {
										health = 0
									}
									showFloatingText(c.x, c.y, fmt.Sprintf("%d", hpChange), color.RGBA{R: 255, G: 50, B: 50, A: 255})
									healthText.Text = fmt.Sprintf("Health: %d", health)
									healthText.Refresh()
								}

								if health == 0 && !gameOver {
									gameOver = true
									gameActive = false
									gameOverText := canvas.NewText("GAME OVER", color.White)
									gameOverText.TextSize = 48
									gameOverText.TextStyle.Bold = true
									gameOverText.Move(fyne.NewPos(250, 250))
									content.Add(gameOverText)

									newGameBtn = widget.NewButton("New Game", func() {
										fyne.Do(func() {
											content.Remove(gameOverText)
											startGame()
										})
									})
									newGameBtn.Resize(fyne.NewSize(200, 50))
									newGameBtn.Move(fyne.NewPos(300, 320))
									content.Add(newGameBtn)
									content.Refresh()
								}
							}
						})
					}

					fyne.Do(func() {
						if !c.removed {
							content.Remove(c.circle)
							content.Remove(c.button)
							c.removed = true
						}
					})
				}(c)
			}
		}
	}()

	startGame()
	w.ShowAndRun()
}
