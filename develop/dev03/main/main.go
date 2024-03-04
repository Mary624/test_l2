package main

import (
	"fmt"
	mainsort "utility-sort/pkg/main_sort"
)

func main() {
	b, err := mainsort.Sort()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !b {
		fmt.Println("is sorted")
	}
}
