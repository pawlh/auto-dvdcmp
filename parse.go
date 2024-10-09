package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type Extra struct {
	Title    string
	Duration time.Duration
}

func ParseRawExtras(rawExtras []string) []Extra {
	extras := make([]Extra, 0)
	for _, rawExtra := range rawExtras {
		isExtra, extra := parseExtra(rawExtra)
		if !isExtra {
			continue
		}
		extras = append(extras, extra)
	}
	return extras
}

func parseExtra(rawExtra string) (bool, Extra) {
	validExtra, rawTitle, rawDuration := splitExtra(rawExtra)
	if !validExtra || !isValidDurationFormat(rawDuration) {
		return false, Extra{}
	}

	title := cleanupTitle(rawTitle)

	duration, err := parseDuration(rawDuration)
	if err != nil {
		log.Println(err)
		return false, Extra{}
	}

	return true, Extra{
		Title:    title,
		Duration: duration,
	}
}

func cleanupTitle(title string) string {
	title = strings.TrimLeft(title, "- ")
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, ":", "-")
	return title
}

func isValidDurationFormat(rawDuration string) bool {
	durationRegex := regexp.MustCompile(`\d+:\d+:\d+|\d+:\d+`)
	if !durationRegex.MatchString(rawDuration) {
		return false
	}
	return true
}

func parseDuration(rawDuration string) (time.Duration, error) {
	rawDuration = strings.TrimPrefix(rawDuration, "(")
	rawDuration = strings.TrimSuffix(rawDuration, ")")
	durationSegments := strings.Split(rawDuration, ":")
	switch len(durationSegments) {
	case 2:
		formattedDuration := fmt.Sprintf("%sm%ss", durationSegments[0], durationSegments[1])
		return time.ParseDuration(formattedDuration)
	case 3:
		formattedDuration := fmt.Sprintf("%sh%sm%ss", durationSegments[0], durationSegments[1], durationSegments[2])
		return time.ParseDuration(formattedDuration)
	}
	return 0, fmt.Errorf("invalid duration: %s", rawDuration)
}

// splitExtra splits a line into a title and a duration
// returns false if the line does not contain a duration
func splitExtra(line string) (bool, string, string) {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return false, "", ""
	}
	title := strings.Join(parts[:len(parts)-1], " ")
	duration := parts[len(parts)-1]
	return true, title, duration
}

// FindMatches returns extras that are a close match to the given duration
func FindMatches(extras []Extra, duration time.Duration) []Extra {
	var matches []Extra

	const tolerance = 750 * time.Millisecond
	for _, extra := range extras {
		if extra.Duration >= duration-tolerance && extra.Duration <= duration+tolerance {
			matches = append(matches, extra)
		}
	}
	return matches
}
