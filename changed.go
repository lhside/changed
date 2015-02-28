package main

import (
  "flag"
  "fmt"
  "log"
  "os/exec"
  "strings"

  "github.com/howeyc/fsnotify"
)

type Changed struct {
  Watcher *fsnotify.Watcher
  File    string
  Action  string
}

func NewChanged(ifFile, doAction string) *Changed {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatal(err)
  }
  // watcher.Watch(ifFile)
  watcher.Watch(".")

  return &Changed{
    Watcher: watcher,
    File:    ifFile,
    Action:  doAction,
  }
}

func (c *Changed) Do() {
  parts := strings.Fields(c.Action)
  cmd := exec.Command(parts[0], parts[1:]...)
  out, err := cmd.CombinedOutput()
  if err != nil {
    log.Println(err)
  }
  fmt.Printf("%s\n", out)
}

func main() {
  var ifFile, doAction string

  flag.StringVar(&ifFile, "if", "", "specify path")
  flag.StringVar(&doAction, "do", "", "specify command")
  flag.Parse()

  changed := NewChanged(ifFile, doAction)
  done := make(chan bool)

  go func() {
    for {
      select {
      case ev := <-changed.Watcher.Event:
        if ev.Name == changed.File && (ev.IsModify()) {
          changed.Do()
        }
      case err := <-changed.Watcher.Error:
        log.Println("error:", err)
        done <- true
      }
    }
  }()

  <-done
  changed.Watcher.Close()
}
