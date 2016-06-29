// TODO LIST:
// * remove/hide global variables

package main
import "github.com/gopherjs/gopherjs/js"

const (
	ROWS = 100
	COLS = 170
	GENERATIONS = 100
	CELL_SIZE = 12
	TICK_LENGTH = 200 // milliseconds
	INITIAL_ALIVE_THRESHOLD = 35 // percentage
)

var g_rows, g_cols int

// GLOBAL VARIABLES
var colors []string = []string{"red", "blue", "brown", "yellow", "green", "lightyellow", "orange", "cyan", "pink", "lime"}
var color string
var Width, Height int
var g_math *js.Object
var g_board Board

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
			if random(100) < INITIAL_ALIVE_THRESHOLD {
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
	body.Call("appendChild", canvas)

	// SETUP EVENTS
	canvas.Call("addEventListener", "click", func() {
		color = colors[random(len(colors))]
	})

	return canvas
}

// ======================================================================
// auxiliary functions
// ======================================================================

func random(n int) int {
	f := g_math.Call("random").Float()
	f *= float64(n)
	return int(f)
}

func init() {
	g_math = js.Global.Get("Math")
}

// ======================================================================
// main
// ======================================================================

func main() {
	color = colors[random(len(colors))]

	canvas := createCanvas()
	ctx := canvas.Call("getContext", "2d")

	g_board = *NewBoard(g_rows, g_cols)
	g_board.RandomlySeedBoard()

	js.Global.Call("setInterval", func(){
		g_board.draw(ctx)
		g_board.step()
	}, TICK_LENGTH)
}
