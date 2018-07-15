package main

import (
	"fmt"

	"github.com/uyhcire/hexit/src"
)

func main() {
	for i := 0; i < 1000; i++ {
		fmt.Printf("Played %d games\n", i)
		outputFilename := fmt.Sprintf("%d", i)
		hexit.GenerateTrainingGame(outputFilename)
	}
}
