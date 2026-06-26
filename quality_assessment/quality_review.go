package main

import (
	"fmt"
	"log"
	"os/exec"
)

type Frame struct {
	N        int
	Ssim_y   float32
	Ssim_u   float32
	Ssim_v   float32
	Ssim_avg float32
}

func check(e error) {
	if e != nil {
		// panic(e)
		log.Fatal(e)
	}
}

const PATH_METRICS = "./metrics"

func main() {

	// var frames []Frame

	// // reading file contents
	// dat, err := os.ReadFile("./metrics/file_test-23-x265-fast_ssim.json")
	// check(err)
	// // fmt.Print(string(dat))

	// err = json.Unmarshal(dat, &frames)
	// check(err)
	// for it, el := range frames {
	// 	if it < 10 {
	// 		fmt.Printf("%+v\n", el)
	// 	}
	// }

	// fmt.Printf("%+v\n", frames[:10])

	// dir_list, err := os.ReadDir(PATH_METRICS)
	// check(err)

	// for _, el := range dir_list {
	// 	fmt.Println(el.Name())
	// }

	// cmd := exec.Command("echo", "-n", `{"Name": Bob, "Age": 32}`)
	// cmd := exec.Command("ffmpeg")
	// cmd := exec.Command("cmd", "/c", "echo", "Hello BB", "\n\tfrom", "GG")

	// exiftool.exe "<location>/img.jpg"
	// cmd := exec.Command("cmd", "/c", "exiftool", `<location>/img.jpg`)

	audio_in := `01 - Polyphia - Playing God.m4a`
	const AUDIO_ROOT = `./examples/audio/`
	cmd := exec.Command("cmd", "/c", "ffmpeg.exe", "-i", AUDIO_ROOT+audio_in, "-b:a", "320k", AUDIO_ROOT+"Playing God.mp3")

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))

}
