package server

import (
	"../solver"
	"fmt"
	"github.com/go-martini/martini"
)

func NewServer(port string) {
	m := martini.Classic()
	m.Get("/:board", func(params martini.Params) string {
		if params["board"] == "favicon.ico" || params["board"] == "" {
			return "No board provided"
		}
		result := solver.Solver(params["board"])
		fmt.Println(result)

		return result
	})
	m.RunOnAddr(":" + port)
}
