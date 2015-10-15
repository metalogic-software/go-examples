package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	res, err := http.Get("http://golang.org/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	f, err := os.Create("/tmp/test.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Contents of http://golang.org/index.html saved to /tmp/test.html")
}

