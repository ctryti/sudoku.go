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

type Board struct {
	Width     int
	Height    int
	Dimension int
	Boxes     []Region
	Rows      []Region
	Columns   []Region
	Squares   []Square
}

func initializeBoard(dimension, width, height int) (board *Board) {
	board = new(Board)
	board.Dimension = dimension
	board.Width = width
	board.Height = height

	board.Squares = make([]Square, dimension*dimension)
	board.Rows = make([]Region, dimension)
	board.Boxes = make([]Region, dimension)
	board.Columns = make([]Region, dimension)

	for i := 0; i < dimension; i++ {
		for j := 0; j < dimension; j++ {
			// link the squares next-pointer
			if j != dimension-1 {
				board.Squares[i*dimension+j].Next = &board.Squares[i*dimension+j+1]
			} else if i != dimension-1 {
				board.Squares[i*dimension+j].Next = &board.Squares[(i+1)*dimension]
			}
			board.Squares[i*dimension+j].Row = &board.Rows[i]
			board.Squares[j*dimension+i].Column = &board.Columns[i]
		}
	}

	colNum := 0
	rowNum := 0
	for boxX := 0; boxX < width; boxX++ {
		for boxY := 0; boxY < height; boxY++ {
			for i := 0; i < width; i++ {
				if colNum+width > dimension {
					colNum = 0
					rowNum += width
				}
				for j := 0; j < height; j++ {
					board.Squares[(i+rowNum*dimension)+j+colNum].Box = &board.Boxes[boxX*width+boxY]
				}
			}
			colNum += width
		}
	}

	return board
}

func createBoard(fileName string) (board *Board, e error) {

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
	dimension := translate(scanner.Text())
	scanner.Scan()
	width := translate(scanner.Text())
	scanner.Scan()
	height := translate(scanner.Text())

	board = initializeBoard(dimension, width, height)

	// Read the rest of the file
	for i := 0; i < board.Dimension; i++ {
		scanner.Scan()
		line := scanner.Text()
		for j := 0; j < board.Dimension; j++ {
			value := translate(string(line[j]))
			board.Squares[i*board.Dimension+j].Value = value
			if value == 0 {
				board.Squares[i*board.Dimension+j].Locked = false
			} else {
				board.Squares[i*board.Dimension+j].Locked = true
			}
		}
	}

	// check if translate or the scanner failed
	if err != nil || scanner.Err() != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", board)
	return board, err
}

func main() {
	createBoard(os.Args[1])
}
