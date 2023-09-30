package main

import (
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
	split := strings.Split(strip, "-")
	name := strings.TrimSpace(split[0])
	episode := strings.TrimSpace(split[1])

	return name, episode
}
