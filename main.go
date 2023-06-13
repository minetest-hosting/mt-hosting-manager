package main

import (
	"fmt"
	"mt-hosting-manager/web"
)

func main() {
	fmt.Println("Listening on port 8080")
	err := web.Serve()
	if err != nil {
		panic(err)
	}
}
