package utils

import (
	"path/filepath"
	"shepherd/core/types"
)

func DetectLanguage(path string) types.Language {
	switch filepath.Ext(path) {
	case ".ts", ".tsx":
		return types.LangTypeScript
	case ".js", ".jsx":
		return types.LangJavaScript
	default:
		return types.LangUnknown
	}
}
