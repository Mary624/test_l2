package main

import (
	"fmt"
	"log"
	"utility-wget/pkg/wget"
)

func main() {
	err := wget.Wget()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("downloaded to working folder")
}
