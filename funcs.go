package main

import (
	"fmt"

	"github.com/hellflame/argparse"
)

type Parser struct {
	parser    *argparse.Parser
	argMode   *string
	input     *string
	input2    *string
	output    *string
	overwrite *bool
}

func NewParser() Parser {
	p := Parser{}
	p.parser = argparse.NewParser("konverty", "Simple video converter utilizing FFmpeg", nil)

	p.argMode = p.parser.String("m", "mode", nil)
	p.input = p.parser.String("i", "input", nil)
	p.input2 = p.parser.String("i2", "input2", nil)
	p.output = p.parser.String("o", "output", nil)
	p.overwrite = p.parser.Flag("x", "overwrite", nil)

	if e := p.parser.Parse(nil); e != nil {
		fmt.Println(e.Error())
		// TEMP commented
		// os.Exit(-1)
	}
	return p
}

// # Mode handler
//
// # Detect and handle different operating modes
//
// - convert -
func ModeMenu(mode string) {
	// "convert"
}
