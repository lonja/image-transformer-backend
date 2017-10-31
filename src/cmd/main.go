package main

import (
	"api"
)

func main() {
	server := api.New()
	server.Start(8080)
}
