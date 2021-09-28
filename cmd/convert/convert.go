package main

import (
	"fmt"
	"os"

	"github.com/ksimuk/course-csv-converter/internal/convert"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage: course-csv-converter <input> <output>")
		os.Exit(1)
	}
	src := args[0]
	dst := args[1]
	convert.Convert(src, dst)
}
