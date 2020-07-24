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
	"runtime"
	"strings"

	"github.com/gammazero/workerpool"
)

type options struct {
	quality PresetType
	rootDir string
}

func main() {
	opts, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !verifyDeps() {
		log.Println("One or more dependencies count not be found. Install those and come back later.")
		os.Exit(1)
	}
	log.Printf("superflac is starting. dir = (%s) quality = (%s)\n", opts.rootDir, opts.quality)

	candidates := make([]string, 0)
	err = filepath.Walk(opts.rootDir, func(fileName string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		cleanPath := path.Clean(fileName)
		if !info.IsDir() && strings.EqualFold(path.Ext(cleanPath), ".flac") {
			candidates = append(candidates, cleanPath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("superflac failed - %v\n", err)
		os.Exit(1)
	}
	if err := encodeFlacToMp3(candidates, opts); err != nil {
		fmt.Printf("superflac failed - %v\n", err)
		os.Exit(1)
	}
}

func parseArgs() (options, error) {
	q := flag.String(
		"quality",
		string(PresetTypeExtreme),
		"Specifies the lame preset type to use. Possible values are 'standard', 'extreme', and 'insane'")
	flag.Parse()

	if len(flag.Args()) == 0 {
		return options{}, fmt.Errorf("directory not specified. usage: superflac -quality [standard|extreme|insane] /path/to/encode")
	}
	dir := flag.Args()[0]
	if !isDir(dir) {
		return options{}, fmt.Errorf("directory does not exist: %s", dir)
	}
	return options{
		quality: PresetType(*q),
		rootDir: flag.Args()[0],
	}, nil
}

func encodeFlacToMp3(candidates []string, opts options) error {
	wg := workerpool.New(runtime.NumCPU())

	for _, candidate := range candidates {
		candidate := candidate
		wg.Submit(func() {
			if ok, err := isFlacFileValid(candidate); !ok {
				log.Printf("unable to encode file. `flac` says this file is invalid : %s error : %v\n", candidate, err)
				return
			}

			// create MP3 directory
			outDir := filepath.Join(path.Dir(candidate), "mp3")
			outBase := filepath.Base(candidate)
			outPath := filepath.Join(outDir, strings.TrimSuffix(outBase, filepath.Ext(outBase)))
			outPath = outPath + ".mp3"

			if !ensureDir(outDir) {
				log.Printf("Unable to create output directory at %s\n", outDir)
				return
			}

			if err := os.RemoveAll(outPath); err != nil {
				log.Printf("Unable to remove existing path at %s Error:%v\n", outPath, err)
				return
			}

			// OK : we are ready to encode to MP3

			// Before we start the decode/encode, determine if the flac file is valid
			cmdFlac := exec.Command("flac", "--decode", "--stdout", "--totally-silent", candidate)
			cmdLame := exec.Command("lame", "--preset", string(opts.quality), "--quiet", "-", outPath)

			fmt.Printf("encoding `%s` to `%s`\n", candidate, outPath)
			cmdLame.Stdin, _ = cmdFlac.StdoutPipe()
			cmdLame.Stdout = os.Stdout
			cmdLame.Stderr = os.Stderr
			_ = cmdLame.Start()
			_ = cmdFlac.Start()
			_ = cmdLame.Wait()

			// TODO - @dra - check the return codes from either process. If
			// either fails, error out.
		})
	}
	wg.StopWait()
	// @dra - better error handling?
	return nil
}

// ensureDir will create a directory at `fileName` if one
// does not already exist. If a file already exists, or if
// we are not able to create a direcotry, returns false
func ensureDir(dirPath string) bool {
	if !isDir(dirPath) {
		err := os.MkdirAll(dirPath, 0777)
		return err == nil
	} else {
		return true
	}
}

func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return stat.IsDir()
}
