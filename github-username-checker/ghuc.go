package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"net/http"
	"os"
)

func main() {
	// elite ascii art is a prerequisite for hacking properly
	black := color.FgBlack.Render
	blue := color.FgBlue.Render
	fmt.Printf("%s%s\n", black("https://github.com/rip/go-ghoul/"), blue("github-username-checker"))
	color.Magenta.Println("                        _                _ ")
	color.Magenta.Println("  __ _  ___  ___  __ _ | |_   ___  _  _ | |")
	color.Magenta.Println(" / _` |/ _ \\|___|/ _` || ' \\ / _ \\| || || |")
	color.Magenta.Println(" \\__, |\\___/     \\__, ||_||_|\\___/ \\_,_||_|")
	color.Magenta.Println(" |___/           |___/")
	// user input
	usernamesPath := flag.String("u", "", "path to file containing usernames")
	threads := flag.Int("t", 99, "amount of simultaneous checking threads")
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	// validate input
	if *usernamesPath == "" || *threads < 1 {
		color.Red.Println("-h for help")
		os.Exit(1)
	}
	// initiate janky golang threadpooling
	semaphore := make(chan bool, *threads)
	// read file
	f, err := os.Open(*usernamesPath)
	if err != nil {
		fmt.Println(err)
	}
	// populate slices with requisite data
	s := bufio.NewScanner(f)
	// for slice in slices
	for s.Scan() { // todo turn this into a slice to remove successfully checked and iterate over until all are checked because of request/proxy errors
		username := s.Text()
		// mark thread "in use"
		semaphore <- true
		// invoke reave with a goroutine
		go func(username string) {
			// mark thread available after anonymous function has completed
			defer func() { <-semaphore }()
			// do stuff
			r, err := http.Get("https://github.com/" + username)
			if err != nil {
				color.Yellow.Println(err)
			} else { // don't panic on errors
				// friendly reminder to close stuff
				r.Close = true
				defer r.Body.Close()
				// psst... person reading this, ur a nerd 🤓
				switch r.StatusCode {
				case 404:
					color.Green.Println(username)
				default:
					if *verbose {
						color.Red.Println(username)
					}
				}
			}
		}(username) // arg... removing this seems to just skip first one?
	}
	// clean up thread pool on loop completion
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}
}
