package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type MkvFile struct {
	Path     string
	Extra    Extra
	Duration time.Duration
}

func ScanMkvFiles(workingDirectory string) []MkvFile {
	files, err := os.ReadDir(workingDirectory)
	if err != nil {
		panic(err)
	}

	var mkvFiles []MkvFile
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".mkv") {
			continue
		}

		duration := getVideoDuration(filepath.Join(workingDirectory, file.Name()))

		mkvFiles = append(mkvFiles, MkvFile{
			Path:     file.Name(),
			Duration: duration,
		})
	}

	return mkvFiles
}

func getVideoDuration(path string) time.Duration {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		panic("ffmpeg is not installed")
	}

	cmd := exec.Command("ffmpeg", "-i", path)
	output, _ := cmd.CombinedOutput()
	//if err != nil {
	//	panic(err)
	//}

	durationRegEx := regexp.MustCompile(`Duration: (\d{2}):(\d{2}):(\d{2})`)
	matches := durationRegEx.FindStringSubmatch(string(output))
	if len(matches) != 4 {
		panic("Could not find duration for " + path)
	}

	hours, _ := time.ParseDuration(matches[1] + "h")
	minutes, _ := time.ParseDuration(matches[2] + "m")
	seconds, _ := time.ParseDuration(matches[3] + "s")

	return hours + minutes + seconds
}
