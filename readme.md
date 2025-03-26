# ytgo
A CLI tool to download YouTube videos in any resolution. It download in 24 or 60 FPS automatically.

# Flags
```bash
./ytgo --help
```	
		
Shorthand|Name|Usage
:--:|:---------:|:-----:|
-F | --1080p | Set 1080p resolution for your link
-2 | --1440p | Set 1440p resolution for your link
-4 | --2160p | Set 2160p resolution for your link
-D | --720p | Set 720p resolution for your link
-A | --audio | Download only the audio from your link
-h | --help | help for ytgo
-L | --link string | Paste your link between quotation marks


# Technologies
This project is made 100% in Go, yt-dlp to download videos and FFMPEG to merge and convert video. CobraCLI to made the terminal interface.

```bash
pip install yt-dlr

go get -u github.com/spf13/cobra@latest
```
https://www.ffmpeg.org/download.html


## Get ID codes
```go
	cmdVideo := exec.Command("yt-dlp", "-F", videoURL)
	output, err := cmdVideo.Output()
	if err != nil {
		fmt.Println("Download goes wrong: ", err)
		return
	}
	fmt.Println(string(output))
```


# ID codes from yt-dlp
ID  |Resolution |  FPS  | Type
:--:|:---------:|:-----:|:----:
233 |    x      |   x   | Audio
609 | 1280x720  | 30fps | Video
612 | 1280x720  | 50/60fps | Video
614 | 1920x1080 | 30fps | Video
617 | 1920x1080 | 50/60fps | Video
620 | 2560x1440 | 30fps | Video
623 | 2560x1440 | 60fps | Video
625 | 3840x2160 | 30fps | Video
628 | 3840x2160 | 60fps | Video


# Notes
Using Go Routines, time increase just a few seconds. For downloading a 1080p 30FPS video with 5m18s, it takes 1m24s using Go Routines, and not using it takes 1m28s. For a audio in the same video, drop for 15.9s to 15.7s. I use a MacBook M2 with 8GB to run this simple test. I'll let this code with Go Routines just to good practice.


# Next steps:
- [X] Modularize
- [X] Implementing GoRoutines
- [X] Set resolution from command line too
- [X] Use CLI (with Cobra)
- [X] Download only audio
- [ ] Download playlist
- [ ] Progress bar
- [ ] Host in my main PC
- [ ] API to download from my main PC and send to another