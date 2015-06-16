package server

import (
	"fmt"
	"github.com/go-martini/martini"
	"no/ctryti/sudoku/solver"
)

func NewServer(port string) {
	m := martini.Classic()
	m.Get("/:board", func(params martini.Params) string {
		result := solver.Solver(params["board"])
		fmt.Println(result)
		return result
	})
	m.RunOnAddr(":" + port)
}
