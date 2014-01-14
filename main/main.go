package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	
	"github.com/julien-c/epubserver"
)


func openChrome() {
	cmd := exec.Command("open", "http://localhost:8080")
	cmd.Start()
}


func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("You need to specify a file name.")
		os.Exit(1)
	}
	
	epub, err := epubserver.Open(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	// go openChrome()
	
	epub.Serve()
}


