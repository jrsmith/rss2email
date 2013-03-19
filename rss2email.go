package main

import (
	rss "github.com/jteeuwen/go-pkg-rss"
	"fmt"
	"encoding/json"
	"os"
	"net/smtp"
	"time"
	"sync"
	"io/ioutil"
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

func fetchFeed(url string, timeout int) {

	feed := rss.New(timeout, true, nil, itemHandler)

	for {

		fmt.Printf("Fetching %s\n", url)

		if err := feed.Fetch(url, nil); err != nil {
			fmt.Fprintf(os.Stderr, "[e] %s: %s", url, err)
			return
		}
		<-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))

	}

	wg.Done()

}

func sendItem(subject string, content string) {

	fmt.Println("Sending mail")

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n";

	auth := smtp.PlainAuth(
		"",
		config.SMTP.Username,
		config.SMTP.Password,
		config.SMTP.Host,
	)

	err := smtp.SendMail(
		config.SMTP.OutgoingServer,
		auth,
		config.SMTP.From,
		config.ToEmail,
		[]byte("Subject: "+subject+"\n"+mime+content),
	)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Mail sent")
	}

}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {

	fmt.Printf("%d new item(s) in %s\n", len(newItems), feed.Url)

	for _, item := range newItems {

		var title = item.Title
		var content string

		switch feed.Type {
			case "rss":
				content = item.Description
			case "atom":
				content = item.Content.Text
		}

		wg.Add(1)
		go sendItem(title, content)

	}

}
