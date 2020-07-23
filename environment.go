package main

import (
	"log"
	"os"
	"os/exec"
)

func printEnvironment() {
	for i, v := range os.Args {
		log.Printf("Args[%d] == %s\n", i, v)
	}
	//
	// OS name, version, free memory, etc..
	//
}

func printCmd(c *exec.Cmd) {
	log.Printf("Cmd is  : %v", c)
	log.Printf("Path is : %v", c.Path)
	log.Printf("Args are: %v", c.Args)
	log.Printf("Env  is : %v", c.Env)

}

func verifyDep(cmd string) bool {
	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Printf("NO... Could not find `%v`. Error: \"%v\"", cmd, err.Error())
		return false
	}
	log.Printf("YES... Found `%v` at `%v`\n", cmd, path)
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
