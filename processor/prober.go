package processor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func GetBitrate(file string) int {
	probeArgs := []string{
		"ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=bit_rate",
		"-of", "default=noprint_wrappers=1:nokey=1",
		file,
	}

	cmd := exec.Command(probeArgs[0], probeArgs[1:]...)

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

	bitrateInt, _ := strconv.Atoi(bitrate)
	return bitrateInt / 1000
}

func GetCodec(file string) string {
	probeArgs := []string{
		"ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=codec_name",
		"-of", "default=noprint_wrappers=1:nokey=1",
		file,
	}

	cmd := exec.Command(probeArgs[0], probeArgs[1:]...)

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

func GetTotalFrames(file string) int {
	probeArgs := []string{
		"ffprobe",
		"-v", "error", // INFO: Hide "info" output (version info, etc.)
		"-select_streams", "v:0", // INFO: Select only the first video stream
		"-count_packets",                          // INFO: Count the number of packets per stream and report it in the corresponding stream section
		"-show_entries", "stream=nb_read_packets", // INFO: Show only the entry for "nb_read_packets"
		"-of", "csv=p=0", // INFO: Set the output formatting. In this case it hides the descriptions and only shows the value.
		file,
	}

	cmd := exec.Command(probeArgs[0], probeArgs[1:]...)

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

	framesCountStr := strings.TrimSpace(string(results))
	framesCountInt, _ := strconv.Atoi(framesCountStr)
	return framesCountInt
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
