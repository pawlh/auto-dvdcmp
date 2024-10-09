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
	mkvFiles := ScanMkvFiles(workingDirectory)

	assignments := processMkvFiles(mkvFiles, extras, workingDirectory)
	renameFiles(assignments, workingDirectory)
}

func processMkvFiles(mkvFiles []MkvFile, extras []Extra, workingDirectory string) map[string]string {
	alreadyUsedExtras := make(map[string]Extra)
	assignments := make(map[string]string)

	for _, mkvFile := range mkvFiles {
		matches := FindMatches(extras, mkvFile.Duration)
		handleMatches(mkvFile, matches, alreadyUsedExtras, assignments, workingDirectory)
	}

	return assignments
}

func handleMatches(mkvFile MkvFile, matches []Extra, alreadyUsedExtras map[string]Extra, assignments map[string]string, workingDirectory string) {
	if len(matches) == 0 {
		printNoMatch(mkvFile)
	} else if len(matches) == 1 {
		handleSingleMatch(mkvFile, matches[0], alreadyUsedExtras, assignments, workingDirectory)
	} else {
		printMultipleMatches(mkvFile, matches)
	}
}

func printNoMatch(mkvFile MkvFile) {
	hours := int(mkvFile.Duration.Hours())
	minutes := int(mkvFile.Duration.Minutes()) % 60
	seconds := int(mkvFile.Duration.Seconds()) % 60
	fmt.Printf("%s --> No match. Duration: %d:%d:%d\n", mkvFile.Path, hours, minutes, seconds)
}

func handleSingleMatch(mkvFile MkvFile, match Extra, alreadyUsedExtras map[string]Extra, assignments map[string]string, workingDirectory string) {
	if _, ok := alreadyUsedExtras[match.Title]; ok {
		fmt.Printf("%s --> %s is already used\n", mkvFile.Path, match.Title)
		delete(assignments, alreadyUsedExtras[match.Title].Title)
	}
	fmt.Printf("%s --> %s.mkv\n", mkvFile.Path, match.Title)
	alreadyUsedExtras[match.Title] = match
	assignments[mkvFile.Path] = filepath.Join(workingDirectory, match.Title+".mkv")
}

func printMultipleMatches(mkvFile MkvFile, matches []Extra) {
	fmt.Printf("%s --> Multiple matches\n", mkvFile.Path)
	for i, match := range matches {
		fmt.Printf("  - %d: %s\n", i, match.Title)
	}
}

func renameFiles(assignments map[string]string, workingDirectory string) {
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
