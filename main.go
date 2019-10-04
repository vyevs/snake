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

var (
	bgColor     = colornames.White
	entityColor = colornames.Black
)

func run() {
	window, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title: "Snake",
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

	fmt.Printf("rows: %d, cols: %d\n", gridRows, gridCols)
	fmt.Printf("cellWidthPixels: %d, cellHeightPixels: %d\n", cellWidthPixels, cellHeightPixels)
	fmt.Printf("windowWidth: %d, windowHeight: %d\n", windowWidth, windowHeight)

	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	snake := snake{
		pieces: []vec2{
			randomPosition(gridCols, gridRows),
		},
		direction: directions[rand.Intn(len(directions))], // random initial direction
		max: vec2{
			x: gridCols,
			y: gridRows,
		},
		foodPos: randomPosition(gridCols, gridRows),
	}

	// draw an initial tail based on the initial direction, always behind the direction we are headed
	switch snake.direction {
	case DirectionUp:
		fmt.Println("INITIAL UP")
		for i := 0; i < 10; i++ {
			lastPiece := snake.pieces[0]

			newPiece := vec2{x: lastPiece.x, y: lastPiece.y - 1}
			if newPiece.y < 0 {
				newPiece.y = gridRows - 1
			}

			snake.pieces = append([]vec2{newPiece}, snake.pieces...)
		}
	case DirectionDown:
		fmt.Println("INITIAL DOWN")
		for i := 0; i < 10; i++ {
			lastPiece := snake.pieces[0]

			newPiece := vec2{x: lastPiece.x, y: lastPiece.y + 1}
			if newPiece.y > gridRows {
				newPiece.y = 0
			}

			snake.pieces = append([]vec2{newPiece}, snake.pieces...)
		}
	case DirectionLeft:
		fmt.Println("INITIAL LEFT")
		for i := 0; i < 10; i++ {
			lastPiece := snake.pieces[0]

			newPiece := vec2{x: lastPiece.x + 1, y: lastPiece.y}
			if newPiece.x > gridCols {
				newPiece.x = 0
			}

			snake.pieces = append([]vec2{newPiece}, snake.pieces...)
		}
	case DirectionRight:
		fmt.Println("INITIAL RIGHT")
		for i := 0; i < 10; i++ {
			lastPiece := snake.pieces[0]

			newPiece := vec2{x: lastPiece.x - 1, y: lastPiece.y}
			if newPiece.x < 0 {
				newPiece.x = gridCols - 1
			}

			snake.pieces = append([]vec2{newPiece}, snake.pieces...)
		}
	}

	var paused bool

	for !window.Closed() {
		if window.JustPressed(pixelgl.KeySpace) {
			paused = !paused
		}

		if !paused {
			if window.JustPressed(pixelgl.KeyLeft) && snake.direction.x != 1 {
				snake.direction.y = 0
				snake.direction.x = -1
			} else if window.JustPressed(pixelgl.KeyRight) && snake.direction.x != -1 {
				snake.direction.y = 0
				snake.direction.x = 1
			} else if window.JustPressed(pixelgl.KeyUp) && snake.direction.y != -1 {
				snake.direction.x = 0
				snake.direction.y = 1
			} else if window.JustPressed(pixelgl.KeyDown) && snake.direction.y != 1 {
				snake.direction.x = 0
				snake.direction.y = -1
			}

			snake = snake.move()

			if snake.dead {
				fmt.Println("You died")
				os.Exit(1)
			}

		}

		snake.draw(window)

		window.Update()
	}
}

type snake struct {
	pieces []vec2

	direction vec2

	max vec2

	dead bool

	foodPos vec2
}

type vec2 struct {
	x int
	y int
}

func (s snake) move() snake {
	newX := s.pieces[len(s.pieces)-1].x + s.direction.x
	if newX >= s.max.x {
		newX = newX - s.max.x
	} else if newX < 0 {
		newX = s.max.x + newX
	}

	newY := s.pieces[len(s.pieces)-1].y + s.direction.y
	if newY >= s.max.y {
		newY = newY - s.max.y
	} else if newY < 0 {
		newY = s.max.y + newY
	}

	newPos := vec2{x: newX, y: newY}

	foodPos := s.foodPos
	var ate bool
	if newPos == s.foodPos {
		foodPos = randomPosition(s.max.x, s.max.y)
		ate = true
	}

	if !ate {
		s.pieces = s.pieces[1:]
	}

	var dead bool
	for _, piece := range s.pieces {
		if newPos == piece {
			dead = true
		}
	}

	s.pieces = append(s.pieces, newPos)

	return snake{
		pieces:    s.pieces,
		direction: s.direction,
		max:       s.max,
		dead:      dead,
		foodPos:   foodPos,
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

				if piecePixelIdx < 0 {
					fmt.Printf("piece: (%d, %d), idx: %d\n", piece.x, piece.y, piecePixelIdx)
				}

				pixels[piecePixelIdx] = entityColor
			}
		}
	}

	foodIdx := s.foodPos.y*windowWidth*cellHeightPixels + s.foodPos.x*cellWidthPixels

	for i := 0; i < cellWidthPixels; i++ {
		foodPixelIdx := foodIdx + i

		for j := 0; j < cellHeightPixels; j++ {
			foodPixelIdx := foodPixelIdx + j*windowWidth

			pixels[foodPixelIdx] = entityColor
		}
	}

	sprite := pixel.NewSprite(&picData, picData.Bounds())

	target.Clear(bgColor)

	sprite.Draw(target, pixel.IM.Moved(target.Bounds().Center()))
}

func randomPosition(maxX, maxY int) vec2 {
	return vec2{
		x: rand.Intn(maxX),
		y: rand.Intn(maxY),
	}
}

var (
	DirectionUp = vec2{
		x: 0,
		y: 1,
	}

	DirectionDown = vec2{
		x: 0,
		y: -1,
	}

	DirectionLeft = vec2{
		x: -1,
		y: 0,
	}

	DirectionRight = vec2{
		x: 1,
		y: 0,
	}
)

var directions = []vec2{
	DirectionUp,
	DirectionDown,
	DirectionLeft,
	DirectionRight,
}
