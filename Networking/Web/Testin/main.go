package main

import (
	"fmt"
	"log"
)

func main() {
	data, err := load_file(css, "style", "")
	if err != nil {
		log.Fatal("File Load Error: ", err)
	}
	fmt.Println(string(data))
}
