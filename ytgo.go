package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Responsible to download from URL both video and audio
func download(vPath, vURL string, n int, done chan bool) {
	switch n {
	case 0: // Download audio
		cmdVideo := exec.Command("yt-dlp", "-f", "251", "-o", vPath, vURL)
		err := cmdVideo.Run()
		if err != nil {
			fmt.Println("Download goes wrong: ", err)
			return
		}
	case 1: // Download video
		cmdVideo := exec.Command("yt-dlp", "-f", "614", "-o", vPath, vURL)
		err := cmdVideo.Run()
		if err != nil {
			fmt.Println("Download goes wrong: ", err)
			return
		}
	}

	done <- true
}

// Convert .mpg final video to .mp4 and with video's real name
func convert(vPath, v2Path string, n int, done chan bool) {
	switch n {
	case 0: // Convert audio
		cmdConvert := exec.Command("ffmpeg", "-i", vPath, "-vn", "-ar", "44100", "-ac", "2", "-b:a", "192k", v2Path)
		err := cmdConvert.Run()
		if err != nil {
			fmt.Println("Audio conversion goes wrong: ", err)
			return
		}
	case 1: // Convert video
		cmdConvert := exec.Command("ffmpeg", "-i", vPath, "-c:v", "libx264", "-preset", "fast", "-crf", "23", "-c:a", "aac", "-b:a", "192k", v2Path)
		err := cmdConvert.Run()
		if err != nil {
			fmt.Println("Video conversion goes wrong: ", err)
			return
		}
	}

	done <- true
}

// Overlap audio and video together and get video's name
func overlap(vPath, aPath, vURL string, done chan bool) {
	cmdFilename := exec.Command("yt-dlp", "-f", "best", "-o", "%(title)s.%(ext)s", "--print", "filename", vURL)

	var out bytes.Buffer
	cmdFilename.Stdout = &out

	err := cmdFilename.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	outputVideo := "/Users/cesar/Downloads/" + strings.TrimSpace(out.String()) // Remove \n and extra spaces

	cmdCombine := exec.Command("ffmpeg", "-i", vPath, "-i", aPath, "-c:v", "copy", "-c:a", "aac", outputVideo)

	err = cmdCombine.Run()
	if err != nil {
		fmt.Println("Overlap goes wrong: ", err)
		return
	}

	done <- true
}

// Delete all temporary files
func delete(file string) {
	err := os.Remove(file)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	cda := make(chan bool)
	cdv := make(chan bool)
	cdo := make(chan bool)
	cdcv := make(chan bool)
	cdca := make(chan bool)

	var videoURL string
	var videoPath string
	var videotwoPath string
	var audioPath string
	var audiotwoPath string
	path := "/usr/local/gotemp/" // I created this path before
	videoPath = path + "video.mpg"
	videotwoPath = path + "videotwo.mp4"
	audioPath = path + "audio.mpg"
	audiotwoPath = path + "audiotwo.mp3"

	// dirPath := "/usr/local/gotemp"
	fmt.Println("Downloading video...")

	videoURL = os.Args[1]
	// videoRes := os.Args[2]
	// var resolution string

	go download(audioPath, videoURL, 0, cda)
	go download(videoPath, videoURL, 1, cdv)
	<-cda
	<-cdv

	fmt.Println("Download has finished, starting conversion...")

	go convert(audioPath, audiotwoPath, 0, cdca)
	go convert(videoPath, videotwoPath, 1, cdcv)
	<-cdca
	<-cdcv

	fmt.Println("Conversion has finished, starting overlaping...")
	go overlap(videotwoPath, audiotwoPath, videoURL, cdo)
	<-cdo

	delete(videoPath)
	delete(audioPath)
	delete(videotwoPath)
	delete(audiotwoPath)
	fmt.Println("Download is over!")
}

// ID codes
// 609: 1280x720   30 fps
// 612: 1280x720   60 fps

// 614: 1920x1080  30 fps
// 617: 1920x1080  60 fps

// 620: 2560x1440  30 fps
// 623: 2560x1440  60 fps

// 625: 3840x2160  30 fps
// 628: 3840x2160  60 fps
