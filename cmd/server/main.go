package main

import (
	"log"

	"github.com/1ch0/tv2okx/cmd/server/app"
)

func main() {
	cmd := app.NewAPIServerCommand()
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
