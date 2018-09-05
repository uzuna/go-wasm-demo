package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("a", ":8080", "address:port")
	flag.Parse()

	log.Println(fmt.Printf("Listening on %q", *addr))
	log.Fatal(http.ListenAndServe(*addr, http.FileServer(http.Dir("./desc"))))
}
