package main

import (
	"fmt"
	"konverty/processor"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func main() {
	argParser := NewParser()

	_ = argParser

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	// TODO: path existing verification

	// processor.IsVideoFile(*argParser.input)

	// processor.BatchConvert(*argParser.input, *argParser.output, logger, *argParser.overwrite)

	// processor.GetPSNR(*argParser.input, *argParser.input2)

	// Process whole folder
	// files, _ := os.ReadDir(*argParser.input)
	// for _, file := range files {
	// 	origFilePath := filepath.Join(*argParser.input, file.Name())
	// 	procFilePath := filepath.Join(*argParser.input2, fmt.Sprintf("%s.mp4", strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))))
	// 	processor.CheckFileConversionMatch(origFilePath, procFilePath, logger)
	// }

	// processor.CheckFileConversionMatch(*argParser.input, *argParser.input2, logger)

	// processor.Convert(*argParser.input, *argParser.output, logger, false)
	// processor.QualityAssessment(*argParser.input, *argParser.output, logger)

	// TEMP: testing 23 fast vs 30 medium quality for "quality-tests/in"
	// Process files in a folder
	files, _ := os.ReadDir(*argParser.input)
	for _, file := range files {
		if !strings.Contains(file.Name(), "Piknik") {
			continue
		}
		origFilePath := filepath.Join(*argParser.input, file.Name())

		procFilePath := filepath.Join(*argParser.input2, fmt.Sprintf("%s_26_fast.mp4", strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))))
		processor.Convert(origFilePath, procFilePath, logger, "26", "fast", false)
		processor.QualityAssessment(origFilePath, procFilePath, logger)

		procFilePath = filepath.Join(*argParser.input2, fmt.Sprintf("%s_28_fast.mp4", strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))))
		processor.Convert(origFilePath, procFilePath, logger, "28", "fast", false)
		processor.QualityAssessment(origFilePath, procFilePath, logger)
	}
}
