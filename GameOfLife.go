// TODO LIST:
// * separate into 2 files, GameOfLife.go and js.go (and, maybe, termbox?)
// * clean code
// * git repo
// * remove/hide global variables
// * add termbox output support
// * create Makefile to build js-canvas and termbox

package main
import "time"
import "math/rand"
import "github.com/gopherjs/gopherjs/js"

const (
	ROWS = 100
	COLS = 170
	GENERATIONS = 100
	CELL_SIZE = 12
	TICK_LENGTH = 100 // milliseconds
	INITIAL_ALIVE_THRESHOLD = 35 // percentage
)

var g_rows, g_cols int

// GLOBAL VARIABLES
var Generation int
var colors []string = []string{"red", "blue", "brown", "yellow", "green", "lightyellow", "orange", "cyan", "pink", "lime"}
var color string
var Width, Height int

/*
type Canvas struct {
	canvas *js.Object
	ctx *js.Object
}

func (c *Canvas) width() int {
	return c.canvas.Get("width")
}

func (c *Canvas) height() int {
	return c.canvas.Get("height")
}
*/


// ======================================================================
// cell
// ======================================================================
type cell struct {
	alive bool
	generation int
}

// ======================================================================
// Board
// ======================================================================
type Board struct {
	rows, cols int
	board [][][]cell
	current, temp int
	generation int
}

func NewBoard(r, c int) *Board {
	b := &Board{}
	b.rows = r
	b.cols = c

	b.board = make([][][]cell, 2, 2)
	b.current = 0
	b.temp = 1
	b.board[b.current] = make([][]cell, b.rows, b.rows)
	b.board[b.temp] = make([][]cell, b.rows, b.rows)

	for y := 0; y < b.rows; y++ {
		b.board[b.current][y] = make([]cell, b.cols, b.cols)
		b.board[b.temp][y] = make([]cell, b.cols, b.cols)

	}

	return b
}

func (b *Board) RandomlySeedBoard() *Board {
	for r:=0; r<b.rows; r++ {
		for c:=0; c<b.cols; c++ {
			if rand.Intn(100) < INITIAL_ALIVE_THRESHOLD {
				b.board[b.current][r][c].alive = true
			} else {
				b.board[b.current][r][c].alive = false
			}
		}
	}

	return b
}

func (b *Board) cellNeighbours(r, c int) (neighbours int) {
	isValid := func(r, c int) bool {
		return (r >= 0 && r < b.rows) && (c >= 0 && c < b.cols)
	}

	// c-1
	if isValid(r-1, c-1) && b.board[b.current][r-1][c-1].alive {
		neighbours++
	}
	if isValid(r, c-1) && b.board[b.current][r][c-1].alive {
		neighbours++
	}
	if isValid(r+1, c-1) && b.board[b.current][r+1][c-1].alive {
		neighbours++
	}

	// c
	if isValid(r-1, c) && b.board[b.current][r-1][c].alive {
		neighbours++
	}
	if isValid(r+1, c) && b.board[b.current][r+1][c].alive {
		neighbours++
	}

	// c+1
	if isValid(r-1, c+1) && b.board[b.current][r-1][c+1].alive {
		neighbours++
	}
	if isValid(r, c+1) && b.board[b.current][r][c+1].alive {
		neighbours++
	}
	if isValid(r+1, c+1) && b.board[b.current][r+1][c+1].alive {
		neighbours++
	}

	return
}

// step - make a new generation. Kill dead cells, make alive ones.
func (b *Board) step() {
	for r:=0; r<b.rows; r++ {
		for c:=0; c<b.cols; c++ {
			n := b.cellNeighbours(r, c)
			if n == 3 {
				b.board[b.temp][r][c] = b.board[b.current][r][c]
				if b.board[b.temp][r][c].alive == false {
					b.board[b.temp][r][c].alive = true
					b.board[b.temp][r][c].generation = b.generation
				}
			} else if n == 2 {
				b.board[b.temp][r][c] = b.board[b.current][r][c]
			} else {
				b.board[b.temp][r][c].alive = false
			}
		}
	}

	b.current, b.temp = b.temp, b.current // amazing swap $ -)
}

func (b *Board) draw(ctx *js.Object) {
	for r:=0; r<b.rows; r++ {
		for c:=0; c<b.cols; c++ {
			ctx.Set("strokeStyle", "rgba(242, 198, 65, 0.1)")
			ctx.Call("strokeRect", c*CELL_SIZE, r*CELL_SIZE, CELL_SIZE, CELL_SIZE)
			if b.board[b.current][r][c].alive {
				ctx.Set("fillStyle", color)
			} else {
				ctx.Set("fillStyle", "rgb(38, 38, 38)")
			}
			ctx.Call("fillRect", c*CELL_SIZE, r*CELL_SIZE, CELL_SIZE, CELL_SIZE)
		}
	}
}

// ======================================================================
// web part
// ======================================================================
func createCanvas() *js.Object {
	document := js.Global.Get("document")

	width := js.Global.Get("innerWidth").Int()
	height := js.Global.Get("innerHeight").Int()
	// ugly, just to fully eliminate scrool areas
	width -= CELL_SIZE/2
	height -= CELL_SIZE/2

	g_rows = height / CELL_SIZE
	g_cols = width / CELL_SIZE


	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", CELL_SIZE*g_cols)
	canvas.Set("height", CELL_SIZE*g_rows)

	body := document.Get("body")
	body.Get("style").Set("margin", "0px")
	body.Get("style").Set("padding", "0px")
	document.Get("body").Call("appendChild", canvas)

	println("width:", width, ", g_cols:", g_cols)
	println("height:", height, ", g_rows:", g_rows)

	// SETUP EVENTS
	body.Call("addEventListener", "keydown", func() {
		println("keydown")
	})

	body.Call("addEventListener", "keyup", func() {
		println("keyup")
	})

	body.Call("addEventListener", "keypress", func() {
		js.Global.Call("clearInterval", func() {
		})
		println("keypress")
	})

	canvas.Call("addEventListener", "click", func() {
		println("click")
		color = colors[rand.Intn(len(colors))]
	})

	return canvas
}

/*
func createTextOutput() *js.Object {
	document := js.Global.Get("document")
	p := document.Call("createElement", "p")
	document.Get("body").Call("appendChild", p)
	return p
}
*/

// ======================================================================
// main
// ======================================================================

var g_board Board

func main() {
	rand.Seed(time.Now().UnixNano())

	color = colors[rand.Intn(len(colors))]

	canvas := createCanvas()
	ctx := canvas.Call("getContext", "2d")

	g_board = *NewBoard(g_rows, g_cols)
	g_board.RandomlySeedBoard()
	//g_board = *NewRandomBoard(g_rows, g_cols)

	println("Canvas initalized.")

	//println("js.Global.Call(\"Math.random\"):", js.Global.Call("Math.random"))
	//println("js.Global.Call(\"Math.random\"):", js.Global.Call("Math.random"))
	//println(js.Global.Get("document").Call("Math.random"))
	//println(js.Global.Call("random"))


	js.Global.Call("setInterval", func(){
		Generation++
		g_board.draw(ctx)
		g_board.step()
	}, TICK_LENGTH)
}
