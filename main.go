package main

import (
	"ppl-study-planner/internal/tui"
)

func main() {
	if err := tui.Run(); err != nil {
		panic(err)
	}
}
