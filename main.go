package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
