//
// superflac
//
// Recursively encodes flac files into MP3
//
//
// TODO :
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
// log.Panic() == calls
// log.Printf("(%T)%v", obj, obj)
// log.Println(v ...interface{})
//
package main

import (
  // environment information
  "os/user"
  "log"
  "os/exec"
)

func printGreeting() {
  u, err := user.Current()
  if err != nil {
    log.Fatalln("Who are you? I'm unable to find the current user.");
    return
  }
  n := u.Name
  if n == "" {
    n = u.Username
  }
  log.Printf("Welcome, %v, to superflac!\n", n)
}

func printEnvironment() {
  //
  // OS name, version, free memory, etc..
  //
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
  var success = true
  if !verifyDep("lame") {
    success = false
  }
  if !verifyDep("flac") {
    success = false
  }
  return success
}

func main() {
  printGreeting()
  printEnvironment()
  if !verifyDeps() {
    log.Println("One or more dependencies count not be found. Install those and come back later.")
    return
  }
  log.Println("Actually do something here...")
}
