package main

import (
	"flag"
	"log"
	"os"
	"path"
)

var newProject string

func flags() {
	flag.StringVar(&newProject, "new", "", "input your awesome project name")
	flag.Parse()
}

func createFolder(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		pathAll := path.Join(dir, "message")
		log.Println(pathAll)
		return os.MkdirAll(pathAll, 0755)
	}
	return os.ErrExist
}

func main() {
	flags()
	if len(newProject) > 0 {
		// make folder
		createFolder(newProject)
		// make files
	}
}
