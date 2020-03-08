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
// *
// * command line arguments
// * MAXPROCS
// * logging
// *
//
// To understand:
// * pkg/sync
// * panic / recover example
// * pkg/runtime pkg/reflect
// * file io
// * string formatting
// * atos tokenizer
// *

//
// log.Fatal() == os.Exit(1)
// log.Panic() ==
// log.Printf("(%T)%v", obj, obj)
// log.Println(v ...interface{})
//

/*
 * TODO : need a recursive glob function similar to ruby
 */

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
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// current flac settings
var presetType PresetType

func main() {
	printGreeting()
	printEnvironment()

	if !verifyDeps() {
		log.Println("One or more dependencies count not be found. Install those and come back later.")
		return
	}

	// defaults
	dir := "/tmp/2000-04-04"
	presetType = Extreme
	// Parse the given argument

	if len(os.Args) == 1 {
		log.Printf("We were not given a directory. Using our default (%s)\n", dir)
	}
	if len(os.Args) > 1 {
		dir = os.Args[1]
		log.Printf("Superflac is starting using directory (%s)\n", dir)
	}
	if len(os.Args) > 2 {
		presetType = presetTypeFromString(os.Args[2])
	}

	if notExists(dir) {
		log.Fatalf("Superflac is failing '%s' does not exist.", dir)
	}

	log.Printf("Superflac is starting. dir = (%s) presetType = (%s)", dir, presetType.String())

	if err := filepath.Walk(dir, walkFunc); err != nil {
		log.Fatalf("Superflac failed with err (%v)", err)
	}
}

// walkFunc is called for each path. We'll use this to determine
// what the file extension is. If a flac file is encountered,
// we'll encode it.
func walkFunc(fileName string, info os.FileInfo, err error) error {

	if err != nil {
		log.Fatalln("walkFunc received an error : %v", err)
		return err
	}
	if info.IsDir() {
		log.Printf("walkFunc skipping directory at \"%v\"\n", fileName)
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

	if notExists(inPath) {
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
	cmdLame := exec.Command("lame", "--preset", presetType.String(), "--quiet", "-", outPath)
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

	if notExists(dirPath) {
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

func notExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	}
	return false
}
