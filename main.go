package main

import (
	"fmt"
	"konverty/processor"

	"github.com/hellflame/argparse"
)

func main() {
	// args := os.Args[1:]

	parser := argparse.NewParser("konverty", "Simple video converter utilizing FFmpeg", nil)
	input := parser.String("i", "input", nil)
	input2 := parser.String("i2", "input2", nil)

	if e := parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		return
	}

	// --------------------------
	// Bitrate / Codec
	// --------------------------

	// fmt.Printf("File input: %s\n", *input)

	// fmt.Printf("bitrate: %dk\ncodec: %s", processor.Get_bitrate(*input), processor.Get_codec(*input))

	// --------------------------
	// PSNR
	// --------------------------

	fmt.Printf("Orig file: %s\nProcessed file: %s\n", *input, *input2)

	processor.Get_PSNR(*input, *input2)
}
