package processor

import (
	"bufio"
	"fmt"
	"io"
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

// Helper function to process a line and update progress/PSNR values
func processLine_old(line string, progressRegex, psnrRegex *regexp.Regexp, totalFrames int, avgPSNR, minPSNR, maxPSNR *float64, foundPSNR *bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// TODO: if line does not start with...

	// Show progress (frame count, fps, time, speed, and percentage if we have total frames)
	if matches := progressRegex.FindStringSubmatch(line); matches != nil {
		frameStr := strings.TrimSpace(matches[1])
		fpsStr := matches[2]
		timeStr := matches[3]
		speedStr := matches[4]

		if totalFrames > 0 {
			frameNum, _ := strconv.Atoi(frameStr)
			percentage := float64(frameNum) / float64(totalFrames) * 100
			fmt.Printf("\rProcessing... Frame: %s/%d (%.1f%%) | Time: %s | FPS: %s | Speed: %sx",
				frameStr, totalFrames, percentage, timeStr, fpsStr, speedStr)
		} else {
			fmt.Printf("\rProcessing... Frame: %s | Time: %s | FPS: %s | Speed: %sx",
				frameStr, timeStr, fpsStr, speedStr)
		}
	}

	fmt.Printf("\nINFO: parsing line:%s\n", line)

	// Parse final PSNR summary (appears at the end)
	if psnrMatches := psnrRegex.FindStringSubmatch(line); psnrMatches != nil {
		fmt.Println("INFO: parsing PSNR")
		fmt.Printf("INFO: parsing line:%s\n", line)
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

func processLine(line string, progressRegex, psnrRegex *regexp.Regexp, totalFrames int, avgPSNR, minPSNR, maxPSNR *float64, foundPSNR *bool, progressBar *progressbar.ProgressBar) {
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

func Get_PSNR(orig string, processed string) PSNR {
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
					processLine(lineBuf.String(), progressRegex, psnrRegex, totalFrames, &avgPSNR, &minPSNR, &maxPSNR, &foundPSNR, processingProgress)
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
				processLine(line, progressRegex, psnrRegex, totalFrames, &avgPSNR, &minPSNR, &maxPSNR, &foundPSNR, processingProgress)
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
