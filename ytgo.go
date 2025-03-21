package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Responsible to download from URL both video and audio
func download(vPath, vURL, vRes string, done chan bool) {
	cmdVideo := exec.Command("yt-dlp", "-f", vRes, "-o", vPath, vURL)
	err := cmdVideo.Run()
	if err != nil {
		fmt.Println("Download goes wrong: ", err)
		return
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

// Merge audio and video together and get video's name
func merge(vPath, aPath, vURL string, done chan bool) {
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
		fmt.Println("Merge goes wrong: ", err)
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
	cda := make(chan bool)  // Channel Done Audio
	cdv := make(chan bool)  // Channel Done Video
	cdm := make(chan bool)  // Channel Done Merge
	cdcv := make(chan bool) // Channel Done Convert Video
	cdca := make(chan bool) // Channel Done Convert Audio

	var videoURL string
	var videoResolution string
	var videoPath string
	var videotwoPath string
	var audioPath string
	var audiotwoPath string
	path := "/usr/local/gotemp/" // I created this path before
	videoPath = path + "video.mpg"
	videotwoPath = path + "videotwo.mp4"
	audioPath = path + "audio.mpg"
	audiotwoPath = path + "audiotwo.mp3"

	fmt.Println("Downloading video...")

	videoURL = os.Args[1]
	videoResolution = os.Args[2]

	go download(audioPath, videoURL, "251", cda)           // AUDIO
	go download(videoPath, videoURL, videoResolution, cdv) // VIDEO
	<-cda
	<-cdv

	fmt.Println("Download has finished, starting conversion...")

	go convert(audioPath, audiotwoPath, 0, cdca) // AUDIO
	go convert(videoPath, videotwoPath, 1, cdcv) // VIDEO
	<-cdca
	<-cdcv

	fmt.Println("Conversion has finished, starting merging...")
	go merge(videotwoPath, audiotwoPath, videoURL, cdm)
	<-cdm

	delete(videoPath)
	delete(audioPath)
	delete(videotwoPath)
	delete(audiotwoPath)
	fmt.Println("Download is over!")
}
