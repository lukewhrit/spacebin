// +build mage

package main

import "github.com/magefile/mage/sh"

// Build generates a binary of the project
func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	return sh.Run("go", "build", "--ldflags", "-s -w", "-tags", "sqlite", "./")
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

	return sh.RunV("./curiosity")
}
