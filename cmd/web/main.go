package main

import (
	"log"
	"os"

	"github.com/piotrstrzalka/contributors/pkg/service"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	service.Run()
}
