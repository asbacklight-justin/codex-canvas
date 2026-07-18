package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/jackmordaunt/icns"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: iconpack input.png output.icns")
		os.Exit(2)
	}
	input, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer input.Close()
	img, _, err := image.Decode(input)
	if err != nil {
		panic(err)
	}
	output, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer output.Close()
	if err := icns.Encode(output, img); err != nil {
		panic(err)
	}
}
