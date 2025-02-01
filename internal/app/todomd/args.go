package todomd

import (
	"path/filepath"
	"slices"
)

var args struct {
	Filenames []string `arg:"positional"`
}

func cleanFilenames(filenames []string) (result []string) {
	for _, filename := range filenames {
		result = append(result, filepath.Clean(filename))
	}
	slices.Sort(result)
	return slices.Compact(result)
}
