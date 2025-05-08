// Package tests_init Description: This file is used to set the working directory to the project root
package tests_init

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	root, err := findProjectRoot()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(root)
	if err != nil {
		panic(err)
	}
}

func findProjectRoot() (string, error) {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
