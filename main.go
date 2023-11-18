package main

import (
	"github.com/koshkaj/lofictl/cmd"
	"log"
)

func main() {
	err := cmd.CreateRootCommand().Execute()
	if err != nil {
		log.Fatal(err)
	}
}
