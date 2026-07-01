package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	srv := server.CreateServerWithRouter(logger)

	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
