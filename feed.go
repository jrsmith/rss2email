package main

import (
    rss "github.com/jteeuwen/go-pkg-rss"
    "fmt"
    "time"
    "errors"
)

var afterDate = time.Date(2013, 3, 25, 0, 0, 0, 0, time.UTC)
var dateFormats = []string{time.ANSIC, time.UnixDate, time.RubyDate, time.RFC822, time.RFC822Z, time.RFC850, time.RFC1123, time.RFC1123Z, time.RFC3339, time.RFC3339Nano}

func fetchFeed(url string, timeout int) {

    fmt.Printf("Fetching feed %s\n", url)

    feed := rss.New(timeout, true, nil, itemHandler)

    for {

        if err := feed.Fetch(url, nil); err != nil {
            fmt.Printf("%s: %s\n", url, err)
            return
        }
        <-time.After(time.Duration(feed.SecondsTillUpdate() * 1e9))

    }

    wg.Done()

}

func parseDate(date string) (time.Time, error) {

    var (
        parsedDate time.Time
        err error
    )

    if date == "" {
        return parsedDate, errors.New("Null date given")
    }

    for _, format := range dateFormats {

        if dateObject, err := time.Parse(format, date); err == nil {
            return dateObject, nil
        }

        err = err

    }

    return parsedDate, err

}

func itemHandler(feed *rss.Feed, ch *rss.Channel, newItems []*rss.Item) {

    defer func() {
        if r := recover(); r != nil {
            fmt.Println("goroutine panicked:", r)
        }
    }()

    fmt.Printf("Got %d items for %s\n", len(newItems), feed.Url)

    for _, item := range newItems {

        if item.PubDate == "" {
            fmt.Println("Pubdate is null")
            return
        }

        parsedPubDate, err := parseDate(item.PubDate)

        if err != nil {
            fmt.Printf("Error parsing date: %s\n", err)
            return
        }

        if parsedPubDate.After(afterDate) {

            var title = item.Title
            var content string = ""

            fmt.Printf("[%v] %v\n", parsedPubDate, title)

            switch feed.Type {
                case "rss":
                    content = item.Description
                case "atom":
                    content = item.Content.Text
                default:
		    return
            }

            sendItem(title, content)

        }

    }

}
