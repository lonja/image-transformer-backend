package main

import (
	"api"
	"time"
	"runtime/debug"
)

func main() {
	server := api.New()
	memoryRelease(30)
	server.Start(8080)
}

func memoryRelease(interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for range ticker.C {
			debug.FreeOSMemory()
		}
	}()
}