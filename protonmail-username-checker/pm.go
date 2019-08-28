package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gookit/color"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	x, verbose bool
	proxies          []string
)

// some comments and code passed down from ytcracker
func ghoul(u string) bool {
	// initialize http client with 10 second timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// if proxies was populated
	if len(proxies) > 0 {
		// seed the random number generator and pick a proxy from the list at random
		rand.Seed(time.Now().UTC().UnixNano())
		proxy := "http://" + string(proxies[rand.Int()%len(proxies)])
		proxyURL, _ := url.Parse(proxy)
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   20 * time.Second,
		}
	}

	req, _ := http.NewRequest("GET", "https://mail.protonmail.com/api/users/available?Name="+u, nil)
	req.Header.Add("x-pm-apiversion", "3")
	req.Header.Add("x-pm-appversion", "Web_3.16.3")

	// submit the request
	resp, err := client.Do(req)
	if err != nil { // _ err causes panic
		// request failed, probably due to a bad proxy so return false to try checking again.
		return false
	}

	// friendly reminder to always close everything
	resp.Close = true
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	switch resp.StatusCode {
	case 200:
		// avail
		color.Green.Println(u)
		return true
	case 409:
		// taken
		if verbose {
			color.Red.Println(u)
		}
		return true
	default:
		color.Yellow.Println("error")
	}

	return false
}

// readLines reads a file into memory and returns a slice of its lines
// appropriated bufio parts from stackoverflow
func readLines(path string) []string {
	// open
	file, err := os.Open(path)
	if err != nil {
		// can use log.Fatal() here but why bother to import the log package for exceptions
		color.Red.Println("error reading " + path + "!")
		color.Yellow.Println("use -h for help")
		os.Exit(1)
	}
	//close
	defer file.Close() // friendly reminder to close files and stuff :)
	// read
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func main() {
	// elite ascii art is a prerequisite for hacking properly
	black := color.FgBlack.Render
	blue := color.FgBlue.Render
	fmt.Printf("%s%s\n", black("https://github.com/rip/go-ghoul/"), blue("protonmail-username-checker"))
	color.Magenta.Println("                        _                _ ")
	color.Magenta.Println("  __ _  ___  ___  __ _ | |_   ___  _  _ | |")
	color.Magenta.Println(" / _` |/ _ \\|___|/ _` || ' \\ / _ \\| || || |")
	color.Magenta.Println(" \\__, |\\___/     \\__, ||_||_|\\___/ \\_,_||_|")
	color.Magenta.Println(" |___/           |___/")
	// user input
	usernamesPath := flag.String("u", "", "path to file containing usernames")
	proxiesPath := flag.String("p", "", "path to file containing http proxies (ip:port)")
	threads := flag.Int("t", 25, "amount of simultaneous checking threads")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.Parse()
	// validate input
	if *usernamesPath == "" || *threads < 1 {
		color.Red.Println("-h for help")
		os.Exit(1)
	}
	// populate slices with requisite data
	usernames := readLines(*usernamesPath)
	// if a proxy file is specified, load it into a slice
	if *proxiesPath != "" {
		proxies = readLines(*proxiesPath)
	}
	// initiate janky golang threadpooling
	semaphore := make(chan bool, *threads)
	// iterate over username file
	for _, username := range usernames {
		semaphore <- true
		// invoke ghoul with a goroutine
		go func(username string) bool {
			// mark thread available after anonymous function has completed
			defer func() { <-semaphore }()
			// toggle the success flag to break the for loop and continue on to the next username
			x = ghoul(username)
			if x {
				return true
			}
			return false
		}(username)
		if x { //////////doesntseemtoworklol
			x = false
			break
		}
	}
	// clean up thread pool on loop completion
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}
}
