package main

import (
	"fmt"
	"utility-grep/pkg/grep"
)

func main() {
	res, err := grep.Grep()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(res)
}
