package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// COBRA CLI
var (
	link          string
	hd            bool
	fullhd        bool
	twok          bool
	fourk         bool
	audio         bool
	playlistaudio bool
	playlistvideo bool
)

// Playlist videos
type Playlist struct {
	Title   string `json:"title"`
	Entries []struct {
		ID string `json:"id"`
	} `json:"entries"`
}

var rootCmd = &cobra.Command{
	Use: "ytgo is a CLI tool to download YouTube videos in any resolution. It download in 24 or 60 FPS automatically.",
	Run: func(cmd *cobra.Command, args []string) {
		var out bytes.Buffer
		var videoPath string
		var videotwoPath string
		var audioPath string
		var audiotwoPath string
		path := "/usr/local/gotemp/" // I created this path before
		videoPath = path + "video.mpg"
		videotwoPath = path + "videotwo.mp4"
		audioPath = path + "audio.mpg"
		audiotwoPath = path + "audiotwo.mp3"
		playlistFilePath := "/Users/cesar/Downloads/"

		cmdFPS := exec.Command("yt-dlp", "--print", "fps", link)
		cmdFPS.Stdout = &out
		err := cmdFPS.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
		fps := strings.TrimSpace(out.String())

		if playlistaudio {
			cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", link)
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("Download goes wrong:", err)
				return
			}

			var playlist Playlist
			err = json.Unmarshal(output, &playlist)
			if err != nil {
				fmt.Println("JSON problem:", err)
				return
			}
			playlistFilePath := "/Users/cesar/Downloads/" + playlist.Title
			err = os.MkdirAll(playlistFilePath, 0755)
			if err != nil {
				fmt.Println("File creation problem:", err)
				return
			}

			for i := 0; i < len(playlist.Entries); i++ {
				vlink := "https://www.youtube.com/watch?v=" + playlist.Entries[i].ID
				downloadAudio(vlink, audioPath, playlistFilePath)
			}
		}

		if playlistvideo {
			cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", link)
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("Download goes wrong:", err)
				return
			}

			var playlist Playlist
			err = json.Unmarshal(output, &playlist)
			if err != nil {
				fmt.Println("JSON problem:", err)
				return
			}
			playlistFilePath := "/Users/cesar/Downloads/" + playlist.Title
			err = os.MkdirAll(playlistFilePath, 0755)
			if err != nil {
				fmt.Println("File creation problem:", err)
				return
			}

			for i := 0; i < len(playlist.Entries); i++ {
				vlink := "https://www.youtube.com/watch?v=" + playlist.Entries[i].ID
				process(vlink, audioPath, videoPath, audiotwoPath, videotwoPath, "bestvideo", playlistFilePath)
			}
		}

		switch fps {
		case "24.0":
			if hd {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "609", playlistFilePath)
			}
			if fullhd {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "614", playlistFilePath)
			}
			if twok {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "620", playlistFilePath)
			}
			if fourk {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "625", playlistFilePath)
			}
		case "60.0":
			if hd {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "612", playlistFilePath)
			}
			if fullhd {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "617", playlistFilePath)
			}
			if twok {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "623", playlistFilePath)
			}
			if fourk {
				process(link, audioPath, videoPath, audiotwoPath, videotwoPath, "628", playlistFilePath)
			}
		}
		if audio {
			playlistFilePath := "/Users/cesar/Downloads/"
			downloadAudio(link, audioPath, playlistFilePath)
		}
	},
}

// Params for CLI
func init() {
	rootCmd.PersistentFlags().StringVarP(&link, "link", "L", "", "Paste a link from a YouTube video")
	rootCmd.PersistentFlags().BoolVarP(&hd, "720p", "H", false, "Set this resolution for your link")
	rootCmd.PersistentFlags().BoolVarP(&fullhd, "1080p", "F", false, "Set this resolution for your link")
	rootCmd.PersistentFlags().BoolVarP(&twok, "1440p", "2", false, "Set this resolution for your link")
	rootCmd.PersistentFlags().BoolVarP(&fourk, "2160p", "4", false, "Set this resolution for your link")
	rootCmd.PersistentFlags().BoolVarP(&audio, "audio", "A", false, "Download only the audio from your link")
	rootCmd.PersistentFlags().BoolVarP(&playlistaudio, "playlist-audio", "P", false, "Download all videos from playlist, but only audio")
	rootCmd.PersistentFlags().BoolVarP(&playlistvideo, "playlist-video", "Q", false, "Download all videos from playlist")
}

// YOUTUBE DOWNLOAD

// Responsible to download from URL both video and audio
func download(vPath, vURL, vRes string, done chan bool) {
	cmdVideo := exec.Command("yt-dlp", "-f", vRes, "-o", vPath, vURL) // "-o": change the name - need to be save in .mpg because it has some bug that QuickPlayer don't run the video in .mp4
	err := cmdVideo.Run()
	if err != nil {
		fmt.Println("Download goes wrong: ", err)
		return
	}
	done <- true
}

// Responsible to download audio from videos in your playlist
func downloadAudio(link, audioPath, playlistFilePath string) {
	cda := make(chan bool)  // Channel Done Audio
	cdca := make(chan bool) // Channel Done Convert Audio

	fmt.Println("Downloading audio...")
	go download(audioPath, link, "233", cda) // AUDIO
	<-cda
	fmt.Println("Converting audio...")

	cmd := exec.Command("yt-dlp", "--get-title", link)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	title := strings.TrimSpace(string(output))

	outputVideo := playlistFilePath + "/" + title + ".mp3" // Remove \n and extra spaces
	go convert(audioPath, outputVideo, 0, cdca)            // AUDIO
	<-cdca
	delete(audioPath)
	fmt.Println("Download is over!")
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
func merge(vPath, aPath, vURL, playlistFilePath string, done chan bool) {
	cmd := exec.Command("yt-dlp", "--get-title", vURL)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	title := strings.TrimSpace(string(output))
	outputVideo := playlistFilePath + "/" + title + ".mp4" // Remove \n and extra spaces

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

// Entire process to download, convert and merge the video
func process(videoURL, audioPath, videoPath, audiotwoPath, videotwoPath, videoResolution, playlistFilePath string) {
	cda := make(chan bool)  // Channel Done Audio
	cdv := make(chan bool)  // Channel Done Video
	cdm := make(chan bool)  // Channel Done Merge
	cdca := make(chan bool) // Channel Done Convert Audio
	cdcv := make(chan bool) // Channel Done Convert Video

	fmt.Println("Downloading video...")
	go download(audioPath, videoURL, "233", cda)           // AUDIO
	go download(videoPath, videoURL, videoResolution, cdv) // VIDEO
	<-cda
	<-cdv

	fmt.Println("Download has finished, starting conversion...")

	go convert(audioPath, audiotwoPath, 0, cdca) // AUDIO
	go convert(videoPath, videotwoPath, 1, cdcv) // VIDEO
	<-cdca
	<-cdcv

	fmt.Println("Conversion has finished, starting merging...")
	go merge(videotwoPath, audiotwoPath, videoURL, playlistFilePath, cdm)
	<-cdm

	delete(videoPath)
	delete(audioPath)
	delete(videotwoPath)
	delete(audiotwoPath)
	fmt.Println("Download is over!")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
