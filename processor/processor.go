package processor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type PSNR struct {
	Average float64
	Min     float64
	Max     float64
}

func processLinePSNR(line string, progressRegex, psnrRegex *regexp.Regexp, totalFrames int, avgPSNR, minPSNR, maxPSNR *float64, foundPSNR *bool, progressBar *progressbar.ProgressBar) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	if matches := progressRegex.FindStringSubmatch(line); matches != nil {
		frameStr := strings.TrimSpace(matches[1])
		framesCountInt, err := strconv.Atoi(frameStr)
		check(err)
		fpsStr := matches[2]
		timeStr := matches[3]
		speedStr := matches[4]
		_ = fpsStr
		_ = timeStr
		_ = speedStr

		progressBar.Set(framesCountInt)
	}

	// Parse final PSNR summary (appears at the end)
	if psnrMatches := psnrRegex.FindStringSubmatch(line); psnrMatches != nil {
		var parseErr error
		avgStr := psnrMatches[1]
		minStr := psnrMatches[2]
		maxStr := psnrMatches[3]

		*avgPSNR, parseErr = strconv.ParseFloat(avgStr, 64)
		if parseErr == nil {
			*minPSNR, parseErr = strconv.ParseFloat(minStr, 64)
			if parseErr == nil {
				*maxPSNR, parseErr = strconv.ParseFloat(maxStr, 64)
				if parseErr == nil {
					*foundPSNR = true
				}
			}
		}
	}
}

func GetPSNR(orig string, processed string) PSNR {
	probeArgs := []string{
		"ffmpeg",
		"-i", orig,
		"-i", processed,
		"-v", "info", // INFO: further details on: https://ffmpeg.org/ffmpeg.html#toc-Generic-options
		"-stats",
		"-stats_period", "0.1", // Report stats every 0.1 seconds for more frequent updates
		"-lavfi", "psnr",
		"-f", "null",
		"-",
	}

	// Assuming *orig* and *processed* are the same underlying file - retrieve total frames
	totalFrames := GetTotalFrames(orig)

	// Run command
	cmd := exec.Command(probeArgs[0], probeArgs[1:]...)

	// Get stderr pipe for real-time streaming
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		return PSNR{}
	}

	// Parse PSNR values from stderr output
	var avgPSNR, minPSNR, maxPSNR float64
	var foundPSNR bool
	var wg sync.WaitGroup

	// Regex patterns for parsing
	// Format: frame=  195 fps=189 q=-0.0 size=N/A time=00:00:06.50 bitrate=N/A speed=6.31x
	progressRegex := regexp.MustCompile(`frame=\s*(\d+)\s+fps=([\d.]+).*?time=(\d{2}:\d{2}:\d{2}\.\d{2}).*?speed=\s*([\d.]+)x`)
	// Format: [Parsed_psnr_0 @ ...] PSNR y:46.675758 u:48.947649 v:48.843571 average:47.297480 min:44.863792 max:51.908735
	psnrRegex := regexp.MustCompile(`PSNR.*?average:([\d.]+).*?min:([\d.]+).*?max:([\d.]+)`)

	// Read from stderr in a goroutine for real-time processing
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderr)
		var lineBuf strings.Builder

		processingProgress := progressbar.NewOptions(int(totalFrames),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetDescription("[yellow]Processing PSNR[reset]"),
			progressbar.OptionShowCount(),
			progressbar.OptionShowElapsedTimeOnFinish(),
			progressbar.OptionSetTheme(progressbar.ThemeUnicode))

		for {
			b, err := reader.ReadByte()
			if err != nil {
				if err != io.EOF {
					// Error reading, break
				}
				// Process remaining buffer
				if lineBuf.Len() > 0 {
					processLinePSNR(lineBuf.String(), progressRegex, psnrRegex, totalFrames, &avgPSNR, &minPSNR, &maxPSNR, &foundPSNR, processingProgress)
				}
				processingProgress.Finish()
				break
			}

			if b != '\r' && b != '\n' {
				// if there was no NL or CR character - just append the data
				lineBuf.WriteByte(b)
				continue
			}

			// Processing line when CR/NL character detected
			line := lineBuf.String()
			if len(line) > 0 {
				// If line does not start with "frame=" or "[Parsed_psnr" - reset and continue
				if !strings.HasPrefix(line, "frame=") && !strings.HasPrefix(line, "[Parsed_psnr") {
					lineBuf.Reset()
					continue
				}
				processLinePSNR(line, progressRegex, psnrRegex, totalFrames, &avgPSNR, &minPSNR, &maxPSNR, &foundPSNR, processingProgress)
				lineBuf.Reset()
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting ffmpeg: %v\n", err)
		return PSNR{}
	}

	// Wait for stderr reading to complete
	wg.Wait()

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		// FFmpeg may return exit code 1 when using null output, but that's often normal
		// Only report error if we didn't get PSNR values
		if !foundPSNR {
			fmt.Printf("\nError running ffmpeg: %v\n", err)
			return PSNR{}
		}
	}

	fmt.Println() // New line after progress output

	if !foundPSNR {
		fmt.Println("Warning: PSNR values not found in output")
		return PSNR{}
	}

	result := PSNR{
		Average: avgPSNR,
		Min:     minPSNR,
		Max:     maxPSNR,
	}

	fmt.Printf("\nPSNR Statistics:\n")
	fmt.Printf("  Average: %.6f dB\n", result.Average)
	fmt.Printf("  Min:     %.6f dB\n", result.Min)
	fmt.Printf("  Max:     %.6f dB\n", result.Max)

	return result
}

func processLineConverter(line string, progressRegex, resultsRegex *regexp.Regexp, totalFrames int, outElapsedTime, outAvgQP *float64, outFramesGenerated *int, ifFinalOutput *bool, progressBar *progressbar.ProgressBar) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	if matches := progressRegex.FindStringSubmatch(line); matches != nil {
		frameStr := strings.TrimSpace(matches[1])
		framesCountInt, err := strconv.Atoi(frameStr)
		check(err)
		fpsStr := matches[2]
		timeStr := matches[3]
		speedStr := matches[4]
		_ = fpsStr
		_ = timeStr
		_ = speedStr

		progressBar.Set(framesCountInt)
	}

	// Parse final convertion summary (appears at the end)
	if convertMatches := resultsRegex.FindStringSubmatch(line); convertMatches != nil {
		var parseErr error
		framesGeneratedStr := convertMatches[1]
		elapsedTimeStr := convertMatches[2]
		avgQPStr := convertMatches[3]

		*outFramesGenerated, parseErr = strconv.Atoi(framesGeneratedStr)
		if parseErr == nil {
			*outElapsedTime, parseErr = strconv.ParseFloat(elapsedTimeStr, 64)
			if parseErr == nil {
				*outAvgQP, parseErr = strconv.ParseFloat(avgQPStr, 64)
				if parseErr == nil {
					*ifFinalOutput = true
				}
			}
		}
	}
}

func Convert(orig string, out_file string) bool {
	probeArgs := []string{
		"ffmpeg",
		"-i", orig,
		"-c:v", "libx265",
		"-crf", "23",
		"-preset", "fast",
		"-c:a", "aac",
		"-b:a", "192k",
		out_file,
	}

	// check whether out_file exists
	if checkPath(out_file) {
		// if exists - ask user if overwrite the file (file will be removed before executing ffmpeg cmd)
		fmt.Printf("File \"%s\" already exists. Overwrite? [y/N]: ", out_file)
		var answer string
		fmt.Scan(&answer)
		if strings.TrimSpace(answer) != "y" && strings.TrimSpace(answer) != "Y" {
			fmt.Printf("Conversion of file \"%s\" aborted due to output file existence.\n", orig)
			return false
		}
		// Removing file
		fmt.Println("DEBUG: Removing file in output location...")
		err := os.Remove(out_file)
		check(err)
		fmt.Println("DEBUG: file successfully removed!")
	}

	// Retrieve total frames
	totalFrames := GetTotalFrames(orig)

	cmd := exec.Command(probeArgs[0], probeArgs[1:]...)

	// Get stderr pipe for real-time streaming (stderr is regular pipe for ffmpeg output format)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
		return false
	}

	var outAvgQP, outElapsedTime float64
	var outFramesGenerated int
	var ifFinalOutput bool
	var wg sync.WaitGroup

	// Regex patterns for parsing
	// Format: frame= 1139 fps= 22 q=28.4 size=   10496KiB time=00:00:38.90 bitrate=2210.4kbits/s speed=0.754x elapsed=0:00:51.57
	progressRegex := regexp.MustCompile(`frame=\s*(\d+)\s+fps=\s*([\d.]+).*?time=(\d{2}:\d{2}:\d{2}\.\d{2}).*?speed=\s*([\d.]+)x.*?`)

	// Parsing final QP results
	// Format: encoded 1068 frames in 17.02s (62.76 fps), 2918.42 kb/s, Avg QP:28.45
	resultsRegex := regexp.MustCompile(`encoded\s*([\d]+)\s*frames\s*in\s*([\d.]+)s.*?Avg QP:([\d.]+)`)

	// Read from stderr in a goroutine for real-time processing
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stderr)
		var lineBuf strings.Builder

		processingProgress := progressbar.NewOptions(int(totalFrames),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetDescription("[yellow]Converting video[reset]"),
			progressbar.OptionShowCount(),
			progressbar.OptionShowElapsedTimeOnFinish(),
			progressbar.OptionSetTheme(progressbar.ThemeUnicode),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]█[reset]",
				SaucerPadding: "[green]░[reset]",
				BarStart:      "[",
				BarEnd:        "]",
			}))

		for {
			b, err := reader.ReadByte()
			if err != nil {
				if err != io.EOF {
					// Error reading, break
				}
				// Process remaining buffer
				if lineBuf.Len() > 0 {
					processLineConverter(lineBuf.String(), progressRegex, resultsRegex, totalFrames, &outElapsedTime, &outAvgQP, &outFramesGenerated, &ifFinalOutput, processingProgress)
				}
				processingProgress.Finish()
				break
			}

			if b != '\r' && b != '\n' {
				// if there was no NL or CR character - just append the data
				lineBuf.WriteByte(b)
				continue
			}

			// Processing line when CR/NL character detected
			line := lineBuf.String()
			if len(line) > 0 {
				// If line does not start with "frame=" or "encoded" - reset and continue
				if !strings.HasPrefix(line, "frame=") && !strings.HasPrefix(line, "encoded") {
					lineBuf.Reset()
					continue
				}
				processLineConverter(line, progressRegex, resultsRegex, totalFrames, &outElapsedTime, &outAvgQP, &outFramesGenerated, &ifFinalOutput, processingProgress)
				lineBuf.Reset()
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting ffmpeg: %v\n", err)
		return false
	}

	// Wait for stderr reading to complete
	wg.Wait()

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		// FFmpeg may return exit code 1 when using null output, but that's often normal
		// Only report error if we didn't get final QP output (some other error occurred that wasn't covered yet)
		// TODO: error catcher
		return false
	}

	fmt.Println() // New line after progress output

	// Print out some final data (??)
	fmt.Printf("Converting finished with total frames processed %d with average QP %f in a total of %fs :)\n", outFramesGenerated, outAvgQP, outElapsedTime)

	return true
}
