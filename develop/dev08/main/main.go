package main

import (
	"log"
	"unix-shell-utility/pkg/executor"
)

func main() {
	e, err := executor.New()
	if err != nil {
		log.Fatal(err)
	}

	e.Run()
}
