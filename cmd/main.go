package main

import (
	"log"
	"sync"

	"github.com/ochiengotieno304/oneotp/cmd/handlers"
	"github.com/ochiengotieno304/oneotp/cmd/servers"
)

func runGrpc(wg *sync.WaitGroup) {
	defer wg.Done()
	servers.StartRPC()
}

func runHttp(wg *sync.WaitGroup) {
	defer wg.Done()
	handlers.RunHandlers()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Server Error:", r)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)
	go runGrpc(&wg)
	go runHttp(&wg)
	wg.Wait()
}
