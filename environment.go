package main

import (
	"fmt"
	"os/exec"
)

func verifyDep(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("could not find `%s` - did you `brew install %s`?", cmd, cmd)
		return false
	}
	return true
}

func verifyDeps() bool {
	if !verifyDep("lame") {
		return false
	}
	if !verifyDep("flac") {
		return false
	}
	return true
}

func isFlacFileValid(path string) (bool, error) {
	cmdFlac := exec.Command("flac", "--test", path)
	err := cmdFlac.Run()
	if err != nil {
		return false, err
	}
	return true, nil
}
