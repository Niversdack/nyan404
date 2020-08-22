package main

import (
	"flag"
	"log"

	"github.com/nyan404/internal/app/apiserver"
)

var (
	configPath string
)

func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
