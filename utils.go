package main

import (
	"path/filepath"
	"strings"

	"github.com/labstack/gommon/color"
)

// Error log error messages.
func logError(str string) {
	color.Println(color.Red("[✖] ") + color.Yellow(str))
}

// Success log success messages.
func logSuccess(str string) {
	color.Println(color.Green("[✔] ") + color.Yellow(str))
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
