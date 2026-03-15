package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprint(w, "Running")
		if err != nil {
			return
		}
	})

	fmt.Println(":8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
