package main

import (
	"konverty/processor"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	argParser := NewParser()

	_ = argParser

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	processor.ConvertBatch(*argParser.input, *argParser.output, logger, *argParser.overwrite)
}
