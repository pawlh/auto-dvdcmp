package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	rawExtrasTest := prompt("Paste in a DVD Compare extras block")
	extras := ParseRawExtras(strings.Split(rawExtrasTest, "\n"))

	workingDirectory := prompt("Enter the working directory")
	workingDirectory = strings.ReplaceAll(workingDirectory, "\\", "")
	mkvFiles := ScanMkvFiles(workingDirectory)

	alreadyUsedExtras := make(map[string]Extra)

	assignments := make(map[string]string)
	for _, mkvFile := range mkvFiles {
		matches := FindMatches(extras, mkvFile.Duration)
		if len(matches) == 0 {
			hours := int(mkvFile.Duration.Hours())
			minutes := int(mkvFile.Duration.Minutes()) % 60
			seconds := int(mkvFile.Duration.Seconds()) % 60
			fmt.Printf("%s --> No match. Duration: %d:%d:%d\n", mkvFile.Path, hours, minutes, seconds)
		} else if len(matches) == 1 {
			if _, ok := alreadyUsedExtras[matches[0].Title]; ok {
				fmt.Printf("%s --> %s is already used\n", mkvFile.Path, matches[0].Title)
				delete(assignments, alreadyUsedExtras[matches[0].Title].Title)
			}
			fmt.Printf("%s --> %s.mkv\n", mkvFile.Path, matches[0].Title)
			alreadyUsedExtras[matches[0].Title] = matches[0]
			assignments[mkvFile.Path] = filepath.Join(workingDirectory, matches[0].Title+".mkv")
		} else {
			fmt.Printf("%s --> Multiple matches\n", mkvFile.Path)
			for i, match := range matches {
				fmt.Printf("  - %d: %s\n", i, match.Title)
			}
		}
	}

	for oldPath, newPath := range assignments {
		err := os.Rename(filepath.Join(workingDirectory, oldPath), newPath)
		if err != nil {
			log.Println(err)
		}
	}
}

func prompt(promptTest string) string {
	fmt.Printf("%s:\n", promptTest)
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}
		text += scanner.Text() + "\n"
	}
	return strings.TrimSuffix(text, "\n")
}
