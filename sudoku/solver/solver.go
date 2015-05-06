package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Square struct {
	Value  int
	Locked bool
	Next   *Square
	Row    *Region
	Column *Region
	Box    *Region
}

type Region struct {
	IsSet []bool
}

func NewRegion(d int) *Region {
	r := new(Region)
	r.IsSet = make([]bool, d)
	return r
}

type Board struct {
	W       int // Width
	H       int // Height
	D       int // Dimension
	Boxes   []*Region
	Rows    []*Region
	Columns []*Region
	Squares []*Square
}

func NewBoard(d, w, h int) *Board {
	b := new(Board)
	b.D = d
	b.W = w
	b.H = h
	b.Boxes = make([]*Region, d)
	b.Rows = make([]*Region, d)
	b.Columns = make([]*Region, d)
	b.Squares = make([]*Square, d*d)

	for i := 0; i < d; i++ {
		for j := 0; j < d; j++ {
			b.Squares[i*d+j] = new(Square)
		}
		b.Rows[i] = NewRegion(d)
		b.Columns[i] = NewRegion(d)
		b.Boxes[i] = NewRegion(d)
	}
	return b
}

func (s *Square) setValue(value int) {
	s.Value = value
	s.Row.IsSet[value-1] = true
	s.Column.IsSet[value-1] = true
	s.Box.IsSet[value-1] = true
}

func (s *Square) resetValue(value int) {
	s.Value = 0
	s.Row.IsSet[value-1] = false
	s.Column.IsSet[value-1] = false
	s.Box.IsSet[value-1] = false
}

func (s *Square) checkValue(value int) bool {
	return true && !s.Row.IsSet[value-1] && !s.Column.IsSet[value-1] && !s.Box.IsSet[value-1]
}

func (s *Square) _solve(board *Board) {
	if !s.Locked {
		for i := 1; i <= board.D; i++ {
			if s.checkValue(i) {
				s.setValue(i)
				if s.Next == nil {
					fmt.Println("Next == nil!")
					board.printBoard()
				} else {
					s.Next._solve(board)
				}
				s.resetValue(i)
			}
		}
	} else if s.Next == nil {
		fmt.Println("Next == nil!")
		board.printBoard()
	} else {
		s.Next._solve(board)
	}

}

func (b *Board) solve() {
	b.Squares[0]._solve(b)
}

func (board *Board) printBoard() {
	for i := 0; i < board.D; i++ {
		for j := 0; j < board.D; j++ {
			value := board.Squares[i*board.D+j].Value
			if value == 0 {
				fmt.Printf(". ")
			} else {
				fmt.Printf("%d ", board.Squares[i*board.D+j].Value)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func initializeBoard(d, w, h int) (board *Board) {

	board = NewBoard(d, w, h)

	for i := 0; i < d; i++ {
		for j := 0; j < d; j++ {
			if i*d+j+1 != d*d {
				board.Squares[i*d+j].Next = board.Squares[i*d+j+1]
			}
			board.Squares[i*d+j].Row = board.Rows[i]
			board.Squares[j*d+i].Column = board.Columns[i]
		}
	}

	colNum := 0
	rowNum := 0
	for boxX := 0; boxX < w; boxX++ {
		for boxY := 0; boxY < h; boxY++ {
			for i := 0; i < w; i++ {
				if colNum+w > d {
					colNum = 0
					rowNum += w
				}
				for j := 0; j < h; j++ {
					board.Squares[((i+rowNum)*d)+j+colNum].Box = board.Boxes[boxX*w+boxY]
				}
			}
			colNum += w
		}
	}
	return board
}

func createBoard(fileName string) (board *Board) {

	var err error

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// Anonymous function that translates '.' to 0, A:Z to 10:37
	translate := func(s string) (i int) {
		if err != nil {
			return
		}
		if s == "." {
			s = "0"
		} else if s >= "A" && s <= "Z" {
			s = "1" + string(rune(s[0])-'A'+'0')
		}
		i, err = strconv.Atoi(s)
		return i
	}

	// Read the first 3 lines of the file
	scanner.Scan()
	d := translate(scanner.Text())
	scanner.Scan()
	w := translate(scanner.Text())
	scanner.Scan()
	h := translate(scanner.Text())

	board = initializeBoard(d, w, h)

	// Read the rest of the file
	for i := 0; i < d; i++ {
		scanner.Scan()
		line := scanner.Text()
		for j := 0; j < d; j++ {
			value := translate(string(line[j]))
			if value == 0 {
				board.Squares[i*d+j].Locked = false
				board.Squares[i*d+j].Value = 0
			} else {
				board.Squares[i*d+j].Locked = true
				board.Squares[i*d+j].setValue(value)
			}
		}
	}

	// check if translate or the scanner failed
	if err != nil || scanner.Err() != nil {
		log.Fatal(err)
	}

	return board
}

func main() {
	board := createBoard(os.Args[1])
	board.printBoard()
	board.solve()
}
