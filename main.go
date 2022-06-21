package main

import (
	"context"
	"sync"

	"github.com/KaanSK/shomon/log"
	"github.com/KaanSK/shomon/service"
)

func main() {
	var wg sync.WaitGroup
	ctx := context.Background()

	srv, err := service.New(&wg, ctx)
	if err != nil {
		log.Error(err.Error())
	}

	wg.Add(1)
	go srv.ListenStream()
	wg.Wait()

}
