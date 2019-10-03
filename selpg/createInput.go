package main

import (
	"fmt"
	"os"
)

func create() {
	file, err := os.OpenFile("test.txt", os.O_RDWR, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to open file !")
		os.Exit(1)
	}
	for i := 1; i <= 200; i++ {
		fmt.Fprintf(file, "%d", i)
		if i%10 == 0 {
			fmt.Fprintf(file, "\f\n")
		} else {
			fmt.Fprintf(file, "\n")
		}
	}
}
