package main

import (
	"fmt"
	"encoding/json"
	"os"
	"sync"
	"io/ioutil"
	"runtime"
)

type SMTPConfig struct {
	Username string
	Password string
	Host string
	OutgoingServer string `json:"outgoing_server"`
	From string
}

type Config struct {
	ToEmail []string `json:"to_email"`
	FeedURLs []string `json:"feed_urls"`
	SMTP SMTPConfig
}

var wg sync.WaitGroup
var config Config

func main() {

	runtime.GOMAXPROCS(4)

	configFile, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	json.Unmarshal(configFile, &config)

	for _, url := range config.FeedURLs {
		wg.Add(1)
		go fetchFeed(url, 1)
	}

	wg.Wait()

}
