// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/sh"
)

// Build generates a binary of the project
func Build() error {
	if os.Getenv("NO_SQLITE") == "1" {
		if err := sh.Run("go", "mod", "download"); err != nil {
			return err
		}

		return sh.Run("go", "build", "--ldflags", "-s -w", "-o", "bin/curiosity", "./")
	} else {
		if err := sh.Run("go", "mod", "download"); err != nil {
			return err
		}

		return sh.Run("go", "build", "--ldflags", "-s -w", "-tags", "sqlite", "-o", "bin/curiosity", "./")
	}
}

// Format lints and fixes all files in the directory
func Format() error {
	return sh.Run("go", "fmt", "./...")
}

// Run builds a binary and executes it
func Run() error {
	err := Build()

	if err != nil {
		return err
	}

	return sh.RunV("./bin/curiosity")
}
