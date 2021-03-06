package main

import (
	"path/filepath"
	"strings"
)

type Ignorer interface {
	IsIgnored(path string) bool
}

type userIgnorer struct {
	ignored     map[string]bool
	ignoredDirs []string
}

func (ug *userIgnorer) IsIgnored(path string) bool {
	if ug.ignored[path] {
		return true
	}
	for _, dir := range ug.ignoredDirs {
		if strings.HasPrefix(path, dir) {
			return true
		}
	}
	return false
}

type smartIgnorer struct {
	includedHiddenFiles map[string]bool
	ui                  *userIgnorer

	// renameDirs is the set of directories that the user did not ask
	// to be watched, but must be watched in order to check if a child
	// is being renamed back in it.
	renameDirs map[string]bool

	// renameChildren is the set of filepaths that the user asked to
	// be watched and that are the direct children of the directories
	// in renameDirs. That is, they are the direct children of
	// directories that the user did not ask to watch but justrun must
	// in order to catch renames. This map is used to make sure we
	// only send an event when the exact files in renameDir they care
	// about are added.
	renameChildren map[string]bool
}

func (si *smartIgnorer) IsIgnored(path string) bool {
	if si.ui.IsIgnored(path) {
		return true
	}
	baseName := filepath.Base(path)
	if strings.HasPrefix(baseName, ".") && !si.includedHiddenFiles[path] {
		return true
	}
	dirPath := filepath.Dir(path)
	if si.renameDirs[dirPath] && !si.renameChildren[path] {
		return true
	}
	return false
}
