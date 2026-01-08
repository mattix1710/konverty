package processor

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/fatih/color"
)

func checkPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// If file exists - return true
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		// If file does NOT exist - return false
		return false
	} else {
		// If it was not possible to specify file existence - display error and return false
		fmt.Println(err)
		return false
	}
}

var videoExt = []string{"mp4", "avi", "mpg", "mpeg", "mkv", "flv", "wmv", "rmvb", "mov"}

func isVideo(fileName string) bool {
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) >= 2 {
		if slices.Contains(videoExt, fileNameParts[len(fileNameParts)-1]) {
			return true
		}
	}
	return false
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func iterDirAndCopy(dir string, copyDir string, logger *log.Logger, ifOverwriteAll bool) {

	logger.Debugf("Creating copy directory with path: %s", copyDir)
	os.MkdirAll(copyDir, 0750)

	logger.Debugf("Start iteration over dir %s", dir)

	files, _ := os.ReadDir(dir)
	for _, subPath := range files {
		origSubPath := filepath.Join(dir, subPath.Name())
		copySubPath := filepath.Join(copyDir, subPath.Name())

		// When directory found
		if subPath.IsDir() {
			// Iterate inside it
			iterDirAndCopy(filepath.Join(dir, subPath.Name()), copySubPath, logger, ifOverwriteAll)
			continue
		}
		// If regular file found
		// If video found - Convert it
		if isVideo(subPath.Name()) {
			logger.Infof("Converting video \"%s\" and saving to directory \"%s\"", subPath.Name(), copyDir)
			copySubPathBaseName := filepath.Join(copyDir, fileNameWithoutExtension(subPath.Name()))
			Convert(origSubPath, fmt.Sprintf("%s.mp4", copySubPathBaseName), logger, ifOverwriteAll)
		} else {
			// Copy other files to its appropriate locations
			logger.Infof("Copying file \"%s\" to directory \"%s\"", subPath.Name(), copyDir)
			copyFile(origSubPath, copySubPath)
		}
	}
}

func iterDir(dir string, copyDir string, logger *log.Logger) {
	logger.Debugf("Creating copy directory with path: %s", copyDir)
	os.MkdirAll(copyDir, 0750)

	logger.Debugf("Start iteration over dir %s", dir)
	files, _ := os.ReadDir(dir)
	fmt.Println("LP\tFILENAME")
	for ii, subPath := range files {
		it := ii + 1
		fmt.Printf("%d\t%s\n", it, subPath.Name())
		origSubPath := filepath.Join(dir, subPath.Name())
		copySubPath := filepath.Join(copyDir, subPath.Name())
		// When directory found
		if subPath.IsDir() {
			// Iterate inside it
			fmt.Printf("========================\n=== Inside dir %s\n", origSubPath)
			iterDir(filepath.Join(dir, subPath.Name()), copySubPath, logger)
			fmt.Printf("=== Exiting dir %s\n========================\n", origSubPath)
			continue
		}
	}
}

func logFileSizeCheck(fileSizeIn, fileSizeOut float64) {
	if fileSizeIn > fileSizeOut {
		clr := color.New(color.FgGreen, color.Bold).SprintFunc()
		sizeDiff := int64((1 - (fileSizeOut / fileSizeIn)) * 100)
		fmt.Printf("File size shrank about %s%%! %s\n", clr(fmt.Sprintf("%d", sizeDiff)), clr(fmt.Sprintf("[%.2fMB vs %.2fMB]", fileSizeIn, fileSizeOut)))
	} else if fileSizeIn < fileSizeOut {
		clr := color.New(color.FgRed, color.Bold).SprintFunc()
		sizeDiff := int64(((fileSizeOut / fileSizeIn) - 1) * 100)
		fmt.Printf("File size rose about %s%%! %s\n", clr(fmt.Sprintf("%d", sizeDiff)), clr(fmt.Sprintf("[%.2fMB vs %.2fMB]", fileSizeIn, fileSizeOut)))
	} else {
		fmt.Println("File size did not change!")
	}
}
