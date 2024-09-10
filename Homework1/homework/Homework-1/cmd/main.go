package main

import (
	"flag"
	"homework/Homework-1/executor"
)

func main() {
	serviceMode := flag.Bool("service", false, "Run in service mode")
	flag.Parse()

	if *serviceMode {
		executor.RunServiceMode()
	} else {
		executor.RunInteractiveMode()
	}

}
