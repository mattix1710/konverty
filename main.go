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
	output := parser.String("o", "output", nil)

	_ = input2

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

	// fmt.Printf("Orig file: %s\nProcessed file: %s\n", *input, *input2)

	// processor.GetPSNR(*input, *input2)

	// --------------------------
	// Converter
	// --------------------------

	fmt.Printf("Orig file: %s\nOut file: %s\n", *input, *output)
	processor.Convert(*input, *output)

}
