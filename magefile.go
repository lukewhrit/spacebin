// +build mage

/*
 * Copyright 2020 Luke Whrit, Jack Dorland; The Spacebin Authors

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 *  Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	}

	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	return sh.Run("go", "build", "--ldflags", "-s -w", "-tags", "sqlite", "-o", "bin/curiosity", "./")
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
