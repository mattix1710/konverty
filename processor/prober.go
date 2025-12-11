package processor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func Get_bitrate(file string) int {
	probe_args := []string{
		"ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=bit_rate",
		"-of", "default=noprint_wrappers=1:nokey=1",
		file,
	}

	cmd := exec.Command(probe_args[0], probe_args[1:]...)

	var stderr strings.Builder
	cmd.Stderr = &stderr

	results, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running ffprobe: %v\n", err)
		if stderr.Len() > 0 {
			fmt.Printf("ffprobe stderr: %s\n", stderr.String())
		}
		return -1
	}
	bitrate := strings.TrimSpace(string(results))

	bitrate_int, _ := strconv.Atoi(bitrate)
	return bitrate_int / 1000
}

func Get_codec(file string) string {
	probe_args := []string{
		"ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=codec_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		file,
	}

	cmd := exec.Command(probe_args[0], probe_args[1:]...)

	var stderr strings.Builder
	cmd.Stderr = &stderr

	results, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running ffprobe: %v\n", err)
		if stderr.Len() > 0 {
			fmt.Printf("ffprobe stderr: %s\n", stderr.String())
		}
		return "-1"
	}
	codec := strings.TrimSpace(string(results))
	return codec
}

type PSNR struct {
	average float64
	min     float64
	max     float64
}

func Get_PSNR(orig string, processed string) PSNR {
	probe_args := []string{
		"ffmpeg",
		"-i", orig,
		"-i", processed,
		"-lavfi", "psnr",
		"-f", "null",
		"-",
	}

	cmd := exec.Command(probe_args[0], probe_args[1:]...)

	var stderr strings.Builder
	cmd.Stderr = &stderr

	results, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running ffmpeg: %v\n", err)
		if stderr.Len() > 0 {
			fmt.Printf("ffmpeg stderr: %s\n", stderr.String())
		}
		return PSNR{}
	}

	fmt.Println(results)

	return PSNR{}
}
