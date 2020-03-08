package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilePathAbs(t *testing.T) {

	testPath := "/test/path.txt"
	if path, _ := filepath.Abs(testPath); path != "/test/path.txt" {
		t.Error("Abs() should return the same absolute value given")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not obtain the current working directory")
	}
	t.Logf("Current working directory : %s", cwd)

	expectedAbsPath := filepath.Join(cwd, "test/path.txt")
	actualAbsPath, _ := filepath.Abs("test/path.txt")
	if actualAbsPath != expectedAbsPath {
		t.Errorf("Abs() should append the given relative path (%s) to cwd (%s).\nActual:%s\nExpected:%s\n",
			"test/path.txt", cwd, actualAbsPath, expectedAbsPath)
	}
}

// path.Base returns the last path component, with extension
// there are special cases for paths with all /// chars
func TestFilePathBase(t *testing.T) {

	// absolute
	if filepath.Base("/a/b/c/path.txt") != "path.txt" {
		t.Error("Abs of a base path failed")
	}
	// relative
	if filepath.Base("a/b/c/path.txt") != "path.txt" {
		t.Error("Abs of a relative path failed")
	}
	// root
	if filepath.Base("/") != "/" {
		t.Error("Abs of / failed")
	}

	// negative
	if filepath.Base("a/b/c/////path.txt") != "path.txt" {
		t.Error("Abs with multiple separators failed")
	}
	if filepath.Base("./dir/.././path.txt") != "path.txt" {
		t.Error("Abs with a non-clean path failed")
	}
	if filepath.Base("//////") != "/" {
		t.Error("Multi-slashed root path failed")
	}
}
