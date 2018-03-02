package ffmpeg

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const FFMPEG_PATH = "ffmpeg.exe"
const FFPROBE_PATH = "ffprobe.exe"

func GetVideoDuration(videoPath string) (float64, error) {
	result, err := exec.Command(FFPROBE_PATH, `-v`, `error`, `-show_entries`, `format=duration`, `-of`, `default=noprint_wrappers=1:nokey=1`, videoPath).Output()
	if err != nil {
		return 0.0, err
	}

	return strconv.ParseFloat(strings.Trim(string(result), "\n\r"), 64)
}

func ffmpegTimeFromSeconds(seconds int64) string {
	return time.Unix(seconds, 0).UTC().Format(`15:04:05.000000`)
}

func CreateVideoThumbnail(videoPath string, thumbnailPath string, thumbnailOffset int64) error {
	return exec.Command(FFMPEG_PATH, `-i`, videoPath, `-ss`, ffmpegTimeFromSeconds(thumbnailOffset), `-vframes`, `1`, thumbnailPath).Run()
}