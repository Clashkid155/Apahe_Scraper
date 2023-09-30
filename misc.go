package main

import (
	"os"
)

// IsExists Check if a file exists
func IsExists(file string) bool {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
