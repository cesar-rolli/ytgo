# ytgo
App to download YouTube videos in any resolution with CLI.


# Technologies
This project is made 100% in Go, yt-dlp to download videos and FFMPEG to merge and convert video.


# ID codes from yt-dlp
### Audio only: 
**251**


### Video only:
- **609:** 1280x720   30 fps
- **612:** 1280x720   60 fps

.
- **614:** 1920x1080  30 fps
- **617:** 1920x1080  60 fps

.
- **620:** 2560x1440  30 fps
- **623:** 2560x1440  60 fps

.
- **625:** 3840x2160  30 fps
- **628:** 3840x2160  60 fps


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