// +build mage

package main

import "github.com/magefile/mage/sh"

func Build() error {
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	return sh.Run("go", "build", "--ldflags", "-s -w", "-tags", "sqlite", "./")
}

func Format() error {
	return sh.Run("go", "fmt", "./...")
}
