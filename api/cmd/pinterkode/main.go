package main

import (
	"fmt"

	"gitlab.com/pinterkode/pinterkode/api/pkg/server"
)

func main() {
	addr := "localhost:8080"
	if err := server.Run(addr, false); err != nil {
		fmt.Println(err)
	}
}
