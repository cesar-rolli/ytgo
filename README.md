# ytgo
App to download YouTube videos in any resolution with CLI.


# Technologies
This project is made 100% in Go, yt-dlp to download videos and FFMPEG to merge and convert video. CobraCLI to made the terminal interface.

```bash
pip install yt-dlr

go get -u github.com/spf13/cobra@latest

https://www.ffmpeg.org/download.html
```


# ID codes from yt-dlp
ID  |Resolution |  FPS  | Type
:--:|:---------:|:-----:|:----:
251 |    x      |   x   | Audio
609 | 1280x720  | 30fps | Video
612 | 1280x720  | 60fps | Video
614 | 1920x1080 | 30fps | Video
617 | 1920x1080 | 60fps | Video
620 | 2560x1440 | 30fps | Video
623 | 2560x1440 | 60fps | Video
625 | 3840x2160 | 30fps | Video
628 | 3840x2160 | 60fps | Video


# Flags
- hd
- fullhd
- 2k
- 4k
- audio


# Notes
For a 8 minutes long 1080p30 video, the entire process takes 2m45s and using goroutines it takes 2m34s, using a M2 8GB RAM MacBook Air.


# Next steps:
- [X] Modularize
- [X] Implementing GoRoutines
- [X] Set resolution from command line too
- [ ] Use CLI (with Cobra)
- [ ] Download only audio
- [ ] Download playlist
- [ ] Host in my main PC
- [ ] API to download from my main PC and send to another