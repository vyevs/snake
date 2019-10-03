package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

const (
	gridRows = 50
	gridCols = 50

	cellHeightPixels = 20
	cellWidthPixels  = 20

	windowWidth  = gridCols * cellWidthPixels
	windowHeight = gridRows * cellHeightPixels
)

func run() {
	window, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title: "Bounce",
		Bounds: pixel.Rect{
			Min: pixel.Vec{
				X: 0,
				Y: 0,
			},
			Max: pixel.Vec{
				X: windowWidth,
				Y: windowHeight,
			},
		},
		VSync: true,
		//Resizable: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	snake := snake{
		pieces:   []vec2{randPos(gridCols, gridRows)},
		velocity: randVelocity(),
		max: vec2{
			x: gridCols,
			y: gridRows,
		},
		foodPos: randPos(gridCols, gridRows),
	}

	var paused bool

	ticker := time.Tick(50 * time.Millisecond)

	for !window.Closed() {
		if window.JustPressed(pixelgl.KeySpace) {
			paused = !paused
		}

		if !paused {
			if window.JustPressed(pixelgl.KeyLeft) && snake.velocity.x != 1 {
				snake.velocity.y = 0
				snake.velocity.x = -1
				fmt.Println("LEFT")
			} else if window.JustPressed(pixelgl.KeyRight) && snake.velocity.x != -1 {
				snake.velocity.y = 0
				snake.velocity.x = 1
				fmt.Println("RIGHT")
			} else if window.JustPressed(pixelgl.KeyUp) && snake.velocity.y != -1 {
				snake.velocity.x = 0
				snake.velocity.y = 1
				fmt.Println("UP")
			} else if window.JustPressed(pixelgl.KeyDown) && snake.velocity.y != 1 {
				snake.velocity.x = 0
				snake.velocity.y = -1
				fmt.Println("DOWN")
			}

			snake = snake.move()

			if snake.dead {
				fmt.Println("You died")
				os.Exit(1)
			}

		}

		snake.draw(window)

		window.Update()

		<-ticker
	}
}

type snake struct {
	pieces []vec2

	velocity vec2

	max vec2

	dead bool

	foodPos vec2
}

type vec2 struct {
	x int
	y int
}

func (s snake) move() snake {
	newX := s.pieces[len(s.pieces)-1].x + s.velocity.x
	if newX >= s.max.x {
		newX = newX - s.max.x
	} else if newX < 0 {
		newX = s.max.x + newX
	}

	newY := s.pieces[len(s.pieces)-1].y + s.velocity.y
	if newY >= s.max.y {
		newY = newY - s.max.y
	} else if newY < 0 {
		newY = s.max.y + newY
	}

	newPos := vec2{x: newX, y: newY}

	foodPos := s.foodPos
	var ate bool
	if newPos == s.foodPos {
		foodPos = randPos(s.max.x, s.max.y)
		ate = true
	}

	if !ate {
		s.pieces = s.pieces[1:]
	}

	var dead bool
	for _, piece := range s.pieces {
		if newPos == piece {
			fmt.Printf("newPos %+v, piece %+v\n", newPos, piece)
			dead = true
		}
	}

	s.pieces = append(s.pieces, newPos)

	return snake{
		pieces:   s.pieces,
		velocity: s.velocity,
		max:      s.max,
		dead:     dead,
		foodPos:  foodPos,
	}
}

func (s snake) draw(target *pixelgl.Window) {
	pixels := make([]color.RGBA, windowWidth*windowHeight)

	picData := pixel.PictureData{
		Pix:    pixels,
		Stride: windowWidth,
		Rect: pixel.Rect{
			Min: pixel.Vec{
				X: 0,
				Y: 0,
			},
			Max: pixel.Vec{
				X: windowWidth,
				Y: windowHeight,
			},
		},
	}

	for _, piece := range s.pieces {
		pieceIdx := piece.y*windowWidth*cellHeightPixels + piece.x*cellWidthPixels

		for i := 0; i < cellWidthPixels; i++ {
			piecePixelIdx := pieceIdx + i

			for j := 0; j < cellHeightPixels; j++ {
				piecePixelIdx := piecePixelIdx + j*windowWidth

				pixels[piecePixelIdx] = colornames.Black
			}
		}
	}

	foodIdx := s.foodPos.y*windowWidth*cellHeightPixels + s.foodPos.x*cellWidthPixels

	for i := 0; i < cellWidthPixels; i++ {
		foodPixelIdx := foodIdx + i

		for j := 0; j < cellHeightPixels; j++ {
			foodPixelIdx := foodPixelIdx + j*windowWidth

			pixels[foodPixelIdx] = colornames.Black
		}
	}

	sprite := pixel.NewSprite(&picData, picData.Bounds())

	target.Clear(colornames.White)

	sprite.Draw(target, pixel.IM.Moved(target.Bounds().Center()))
}

func randPos(maxX, maxY int) vec2 {
	return vec2{
		x: rand.Intn(maxX),
		y: rand.Intn(maxY),
	}
}

var possibleVelocities = []vec2{
	vec2{
		x: -1,
		y: 0,
	},
	vec2{
		x: 0,
		y: -1,
	},
	vec2{
		x: 1,
		y: 0,
	},
	vec2{
		x: 0,
		y: 1,
	},
}

func randVelocity() vec2 {
	return possibleVelocities[rand.Intn(len(possibleVelocities))]
}
