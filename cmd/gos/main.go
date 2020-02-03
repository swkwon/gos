package main

import (
	"flag"
	"fmt"
)

var newProject string
var compile string

func flags() {
	flag.StringVar(&newProject, "new", "", "input your awesome project name")
	flag.StringVar(&compile, "compile", "", "input your messages directory")
	flag.Parse()
}

func main() {
	flags()
	if len(newProject) > 0 {
		// make folder
		// make files
		fmt.Println(newProject)
	} else if len(compile) > 0 {
		fmt.Println(compile)
	}
}
