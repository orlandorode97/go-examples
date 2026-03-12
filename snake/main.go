package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	Width  = 40
	Height = 20
)

type Point struct{ X, Y int }

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Snake struct {
	Body []Point
	Dir  Direction
}

type Game struct {
	Snake Snake
	Food  Point
	Score int
	Over  bool
}

func NewGame() *Game {
	start := Point{Width / 2, Height / 2}
	g := &Game{
		Snake: Snake{
			Body: []Point{start, {start.X - 1, start.Y}, {start.X - 2, start.Y}},
			Dir:  Right,
		},
	}
	g.spawnFood()
	return g
}

func (g *Game) spawnFood() {
	for {
		p := Point{rand.Intn(Width-2) + 1, rand.Intn(Height-2) + 1}
		for _, b := range g.Snake.Body {
			if b == p {
				goto retry
			}
		}
		g.Food = p
		return
	retry:
	}
}

func (g *Game) Update(dir Direction) {
	if g.Over {
		return
	}
	if (dir == Up && g.Snake.Dir == Down) ||
		(dir == Down && g.Snake.Dir == Up) ||
		(dir == Left && g.Snake.Dir == Right) ||
		(dir == Right && g.Snake.Dir == Left) {
		dir = g.Snake.Dir
	}
	g.Snake.Dir = dir

	head := g.Snake.Body[0]
	var newHead Point
	switch g.Snake.Dir {
	case Up:
		newHead = Point{head.X, head.Y - 1}
	case Down:
		newHead = Point{head.X, head.Y + 1}
	case Left:
		newHead = Point{head.X - 1, head.Y}
	case Right:
		newHead = Point{head.X + 1, head.Y}
	}

	if newHead.X <= 0 || newHead.X >= Width-1 || newHead.Y <= 0 || newHead.Y >= Height-1 {
		g.Over = true
		return
	}
	for _, b := range g.Snake.Body {
		if b == newHead {
			g.Over = true
			return
		}
	}
	if newHead == g.Food {
		g.Snake.Body = append([]Point{newHead}, g.Snake.Body...)
		g.Score += 10
		g.spawnFood()
	} else {
		g.Snake.Body = append([]Point{newHead}, g.Snake.Body[:len(g.Snake.Body)-1]...)
	}
}

// writeln writes s followed by \r\n into buf
func writeln(buf *bytes.Buffer, s string) {
	buf.WriteString(s)
	buf.WriteString("\r\n")
}

func (g *Game) Render() {
	grid := make([][]byte, Height)
	for y := range grid {
		grid[y] = make([]byte, Width)
		for x := range grid[y] {
			switch {
			case y == 0 || y == Height-1:
				grid[y][x] = '-'
			case x == 0 || x == Width-1:
				grid[y][x] = '|'
			default:
				grid[y][x] = ' '
			}
		}
	}
	grid[0][0] = '+'
	grid[0][Width-1] = '+'
	grid[Height-1][0] = '+'
	grid[Height-1][Width-1] = '+'
	grid[g.Food.Y][g.Food.X] = '*'
	for i, b := range g.Snake.Body {
		if i == 0 {
			grid[b.Y][b.X] = 'O'
		} else {
			grid[b.Y][b.X] = 'o'
		}
	}

	var buf bytes.Buffer
	buf.WriteString("\033[H") // move cursor home, no clear

	writeln(&buf, fmt.Sprintf("  SNAKE  |  Score: %-4d  |  Length: %-4d", g.Score, len(g.Snake.Body)))
	writeln(&buf, "")
	for _, row := range grid {
		writeln(&buf, string(row))
	}
	writeln(&buf, "")
	writeln(&buf, "  WASD / Arrow Keys: Move  |  Q: Quit")

	if g.Over {
		writeln(&buf, "")
		writeln(&buf, "  +------------------+")
		writeln(&buf, "  |   GAME  OVER!    |")
		writeln(&buf, fmt.Sprintf("  |  Score: %-8d |", g.Score))
		writeln(&buf, "  |  Press R to retry|")
		writeln(&buf, "  +------------------+")
	} else {
		writeln(&buf, "")
		writeln(&buf, "                        ")
		writeln(&buf, "                        ")
		writeln(&buf, "                        ")
		writeln(&buf, "                        ")
		writeln(&buf, "                        ")
	}

	os.Stdout.Write(buf.Bytes())
}

var tty *os.File
var savedStty string

func sttyCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("stty", args...)
	cmd.Stdin = tty
	return cmd
}

func enableRawMode() {
	var err error
	tty, err = os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	out, err := sttyCmd("-g").Output()
	if err != nil {
		panic(err)
	}
	savedStty = strings.TrimSpace(string(out))
	if err := sttyCmd("raw", "-echo", "-isig").Run(); err != nil {
		panic(err)
	}
}

func disableRawMode() {
	sttyCmd(savedStty).Run()
	tty.Close()
}

func readKey() byte {
	buf := make([]byte, 1)
	tty.Read(buf)
	if buf[0] != 27 {
		return buf[0]
	}
	tty.Read(buf)
	if buf[0] != '[' {
		return 0xFF
	}
	tty.Read(buf)
	switch buf[0] {
	case 'A':
		return 201
	case 'B':
		return 202
	case 'C':
		return 203
	case 'D':
		return 204
	}
	return 0xFF
}

func main() {
	enableRawMode()
	defer disableRawMode()

	fmt.Print("\033[?25l\033[2J") // hide cursor, clear screen once
	defer fmt.Print("\033[?25h")

	game := NewGame()
	inputCh := make(chan byte, 16)

	go func() {
		for {
			inputCh <- readKey()
		}
	}()

	ticker := time.NewTicker(120 * time.Millisecond)
	defer ticker.Stop()

	dir := Right
	game.Render()

	for {
		select {
		case b := <-inputCh:
			switch b {
			case 'q', 'Q', 3:
				fmt.Print("\033[2J\033[H\033[?25h")
				return
			case 'r', 'R':
				if game.Over {
					game = NewGame()
					dir = Right
				}
			case 201:
				dir = Up
			case 202:
				dir = Down
			case 203:
				dir = Right
			case 204:
				dir = Left
			case 'w', 'W':
				dir = Up
			case 's', 'S':
				dir = Down
			case 'a', 'A':
				dir = Left
			case 'd', 'D':
				dir = Right
			}
		case <-ticker.C:
			if !game.Over {
				game.Update(dir)
			}
			game.Render()
		}
	}
}
