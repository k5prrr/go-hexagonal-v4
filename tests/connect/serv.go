package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

var number uint64

func handler(w http.ResponseWriter, r *http.Request) {
	currentNumber := atomic.AddUint64(&number, 1)
	fmt.Fprintf(w, "Hello, World! Request number: %d", currentNumber)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running on :3000")
	http.ListenAndServe(":3000", nil)
}
