package main

import (
	"flag"
	"os"
	"strings"
)

// IsExists Check if a file exists
func IsExists(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func GetNameAndEpisode(data string) (string, string) {
	// Watch Dark Gathering - 12 Online
	strip := strings.TrimSuffix(strings.TrimPrefix(data, "Watch "), " Online")
	split := strings.Split(strip, "- ")
	name := strings.TrimSpace(split[0])
	episode := strings.TrimSpace(split[1])

	return name, episode
}

func GetFlags() (*string, *string) {
	quality := flag.String("q", "1080p", "Select anime quality.")
	name := flag.String("n", "", "Anime name")
	flag.Parse()

	return name, quality
}
