package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ksimuk/course-csv-converter/internal/convert"
)

const SUFFIX = "_out"

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: course-csv-converter <folder>")
		os.Exit(1)
	}
	src := args[0]
	files, err := ioutil.ReadDir(src)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		file := f.Name()
		if strings.HasSuffix(file, SUFFIX+".csv") {
			continue
		}
		ext := filepath.Ext(file)
		if ext == ".csv" {
			name := strings.TrimSuffix(file, ext)
			outName := name + SUFFIX + ext
			fmt.Printf("Converting %s result %s\n", file, outName)
			convert.Convert(filepath.Join(src, file), filepath.Join(src, outName))
		}
	}
}
