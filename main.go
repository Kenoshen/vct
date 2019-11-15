package main

import (
	"bufio"
	"os"
)

func main() {
	var reader = bufio.NewReader(os.Stdin)
	data, _ := LoadData()

	current := &Current{
		Data:   data,
		Reader: reader,
	}

	if data == nil {
		Initialize(current)
	}

	for true {
		Menu(current)
	}
}
