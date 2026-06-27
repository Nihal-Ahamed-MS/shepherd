package utils

// TODO: load the data from .gitignore
var SkipDirs = map[string]bool{
	".git": true, "node_modules": true, ".next": true, "build": true,
	"vendor": true, "dist": true, ".cache": true, "bin": true,
}
