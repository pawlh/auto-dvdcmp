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

// ScanMkvFiles scans the working directory for .mkv files and returns a slice of MkvFile structs
func ScanMkvFiles(workingDirectory string) []MkvFile {
	files, err := os.ReadDir(workingDirectory)
	if err != nil {
		panic(err)
	}

	var mkvFiles []MkvFile
	for _, file := range files {
		if isMkvFile(file) {
			mkvFiles = append(mkvFiles, createMkvFile(workingDirectory, file))
		}
	}

	return mkvFiles
}

func isMkvFile(file os.DirEntry) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), ".mkv")
}

func createMkvFile(workingDirectory string, file os.DirEntry) MkvFile {
	duration := getVideoDuration(filepath.Join(workingDirectory, file.Name()))
	return MkvFile{
		Path:     file.Name(),
		Duration: duration,
	}
}

func getVideoDuration(path string) time.Duration {
	checkFfmpegInstalled()

	cmd := exec.Command("ffmpeg", "-i", path)
	output, _ := cmd.CombinedOutput()

	return parseFfmpegDuration(output, path)
}

func checkFfmpegInstalled() {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		panic("ffmpeg is not installed")
	}
}

func parseFfmpegDuration(output []byte, path string) time.Duration {
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
