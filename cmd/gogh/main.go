package main

import (
	"gogh/internal/gogh"
	"gogh/internal/view"
	"log"
)

func main() {
	_gogh, err := gogh.New()
	if err != nil {
		log.Fatalln(err)
	}
	_view := view.New(
		_gogh,
	)
	_view.Run()
}
