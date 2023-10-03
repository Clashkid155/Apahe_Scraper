package config

import (
	"github.com/adrg/xdg"
	"path/filepath"
)

var (
	UserDocument = filepath.Join(xdg.UserDirs.Documents, "Apahe")
	FirstRun     = filepath.Join(UserDocument, ".first_run")
)
