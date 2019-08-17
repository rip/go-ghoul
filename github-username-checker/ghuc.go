package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	u := flag.String("u", "u.txt", "list of usernames")
	flag.Parse()
	f, err := os.Open(*u)
	if err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		r, err := http.Get("https://github.com/" + s.Text())
		if err != nil {
			log.Println(err)
		}
		defer r.Body.Close()
		if r.StatusCode == 404 {
			log.Println(s.Text())
		}
	}
}
