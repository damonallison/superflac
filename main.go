//
// superflac
//
// Recursively encodes flac files into MP3
//
// usage:
//   $ superflac /tmp/dir [preset-type]
//
//     preset-type = [standard | medium | extreme (default) | insane]
//
// TODO :
// * Command line arguments
// * Concurrency: spin up MAXPROCS goroutines to encode simultaneously
// * Need a recursive glob function similar to Ruby (or Python)
// * Logging: Find a true logging package, not fmt.
//
// To understand:
// * pkg/sync
// * pkg/runtime pkg/reflect
// * file io
// * string formatting
// * atos tokenizer

// Start walking our file path.
// try walk, glob
// Glob() will not recurse directories.
// We'll need walk.
// matches, err := filepath.Glob("/tmp/test/*/*.flac")
// if err != nil {
// 	log.Printf("Error looking for flac files == %v\n", err)
// 	return
// }
// log.Printf("glob found %v matches\n", len(matches))
//
// for _, v := range matches {
// 	log.Printf("found flac file at %v\n", v)
// }

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var qFlag string

func main() {

	flag.StringVar(&qFlag, "quality", string(PresetTypeExtreme), "Specifies the lame preset type to use. Possible values are 'standard', 'extreme', and 'insane'")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Printf("Usage: superflac -quality [standard|extreme|insane] /path/to/dir")
		os.Exit(1)
	}
	dir := flag.Args()[0]

	if !verifyDeps() {
		log.Println("One or more dependencies count not be found. Install those and come back later.")
		os.Exit(1)
	}
	if !exists(dir) {
		log.Printf("'%s' does not exist", dir)
		os.Exit(1)
	}

	log.Printf("Superflac is starting. dir = (%s) quality = (%s)", dir, qFlag)

	if err := filepath.Walk(dir, walkFunc); err != nil {
		log.Fatalf("Superflac failed with err (%v)", err)
	}
}

// walkFunc is called for each path. We'll use this to determine
// what the file extension is. If a flac file is encountered,
// we'll encode it.
func walkFunc(fileName string, info os.FileInfo, err error) error {

	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	cleanPath := path.Clean(fileName)
	log.Printf("Base == \"%s\" Dir == \"%s\" Ext == \"%s\"\n",
		path.Base(cleanPath),
		path.Dir(cleanPath),
		path.Ext(cleanPath))

	if path.Ext(cleanPath) == ".flac" {
		if !encodeFlacToMp3(cleanPath) {
			log.Printf("Encode failed for file at %s\n", cleanPath)
		} else {
			log.Printf("Encode success %s\n", cleanPath)
		}
	} else {
		log.Printf("skipping non-flac file at path == \"%s\"\n", cleanPath)
	}
	return nil
}

// convertFlacToMp3 will convert `file` and return a boolean
// indicating if the file was successfully encoded and the path
// where it was encoded to.
func encodeFlacToMp3(inPath string) bool {

	if !exists(inPath) {
		log.Printf("Unable to encode file. path does not exist : %s\n", inPath)
		return false
	}

	// TODO : better way to determine if a file exists than .Stat()?
	fi, err := os.Stat(inPath)
	if err != nil {
		log.Printf("Unable to enode file. Error == %v", err)
		return false
	}

	if fi.IsDir() {
		log.Printf("Unable to encode file. path is not a valid flac file : %s\n", inPath)
		return false
	}

	if ok, err := isFlacFileValid(inPath); !ok {
		log.Printf("Unable to encode file. `flac` says this file is invalid : %s Error : %v\n",
			inPath, err)
		return false
	}

	log.Printf("Attempting to encode flac file at \"%s\"", inPath)

	// create MP3 directory
	outDir := filepath.Join(path.Dir(inPath), "mp3")
	outBase := filepath.Base(inPath)
	outPath := filepath.Join(outDir, strings.TrimSuffix(outBase, filepath.Ext(outBase)))
	outPath = outPath + ".mp3"

	log.Printf("output mp3 filename = %s", outPath)

	if !ensureDir(outDir) {
		log.Printf("Unable to create output directory at %s\n", outDir)
		return false
	}

	if err := os.RemoveAll(outPath); err != nil {
		log.Printf("Unable to remove existing path at %s Error:%v\n", outPath, err)
		return false
	}

	// OK : we are ready to encode to MP3

	// Before we start the decode/encode, determine if the flac file is valid
	cmdFlac := exec.Command("flac", "--decode", "--stdout", "--totally-silent", inPath)
	// --quiet
	cmdLame := exec.Command("lame", "--preset", qFlag, "--quiet", "-", outPath)
	printCmd(cmdFlac)
	printCmd(cmdLame)

	cmdLame.Stdin, _ = cmdFlac.StdoutPipe()
	cmdLame.Stdout = os.Stdout
	cmdLame.Stderr = os.Stderr
	_ = cmdLame.Start()
	_ = cmdFlac.Start()
	_ = cmdLame.Wait()

	// TODO : check the return codes from either process.
	//        if either fails, error out.
	return true
}

// ensureDir will create a directory at `fileName` if one
// does not already exist. If a file already exists, or if
// we are not able to create a direcotry, returns false
func ensureDir(dirPath string) bool {

	if !exists(dirPath) {
		log.Printf("Making directory at (%v)", dirPath)
		err := os.MkdirAll(dirPath, 0777)
		return err == nil
	}

	// Exists - make sure we have a file here.
	fi, err := os.Stat(dirPath)
	if err != nil {
		log.Printf("Unable to stat file at %v", dirPath)
		return false
	}
	return fi.IsDir()
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
