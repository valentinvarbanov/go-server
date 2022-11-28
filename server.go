
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to my website!")
    })

	const port = 8080;
	fmt.Printf("Starting server at port %v\n", port);

	var host = fmt.Sprintf(":%v", port);
    http.ListenAndServe(host, nil)
}