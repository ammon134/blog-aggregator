package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ammon134/blog-aggregator/internal/database"
)

const (
	test_data = `
  <rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
    <channel>
      <title>Boot.dev Blog</title>
      <link>https://blog.boot.dev/</link>
      <description>Recent content on Boot.dev Blog</description>
      <generator>Hugo -- gohugo.io</generator>
      <language>en-us</language>
      <lastBuildDate>Wed, 01 May 2024 00:00:00 +0000</lastBuildDate>
      <atom:link href="https://blog.boot.dev/index.xml" rel="self" type="application/rss+xml"/>

      <item>
        <title>The Boot.dev Beat. April 2024</title>
        <link>https://blog.boot.dev/news/bootdev-beat-2024-04/</link>
        <pubDate>Wed, 03 Apr 2024 00:00:00 +0000</pubDate>
        <guid>https://blog.boot.dev/news/bootdev-beat-2024-04/</guid>
        <description>Pythogoras returned in our second community-wide boss battle. He was vanquished, and there was much rejoicing.</description>
      </item>

      <item>
        <title>Maybe You Do Need Kubernetes</title>
        <link>https://blog.boot.dev/education/maybe-you-do-need-kubernetes/</link>
        <pubDate>Fri, 08 Mar 2024 00:00:00 +0000</pubDate>
        <guid>https://blog.boot.dev/education/maybe-you-do-need-kubernetes/</guid>
        <description>Theo has this great video on Kubernetes, currently titled You Dont Need Kubernetes. Im a Kubernetes enjoyer, but Im not here to argue about that.</description>
      </item>
    </channel>
  </rss>
  `
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Title         string `xml:"title"`
		Link          string `xml:"link"`
		Description   string `xml:"description"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"` // TODO: parse this to time.Time
	Description string `xml:"description"`
}

func fetchDataFromFeed(url string) (*RSSFeed, error) {
	rssFeed := &RSSFeed{}
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, rssFeed)
	if err != nil {
		return nil, err
	}

	return rssFeed, nil
}

func runScrapingWorker(config *Config, batchSize int, interval time.Duration) {
	log.Printf("Collecting %d feeds concurrently every %v...\n", batchSize, interval)
	ctx := context.Background()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		feeds, err := config.DB.GetNextFeedListToFetch(ctx, int32(batchSize))
		if err != nil {
			log.Println("error getting next feed list to fetch")
			continue
		}
		log.Printf("Processing %d feeds...\n", len(feeds))
		var wg sync.WaitGroup
		for _, f := range feeds {
			wg.Add(1)
			go func(f database.Feed) {
				defer wg.Done()
				rss, err := fetchDataFromFeed(f.Url)
				if err != nil {
					log.Printf("error fetching data: %s\n", err.Error())
					return
				}

				f, err = config.DB.MarkFeedFetched(ctx, f.ID)
				if err != nil {
					log.Printf("error fetching data: %s\n", err.Error())
					return
				}

				for _, p := range rss.Channel.Item {
					log.Printf("Processed post: %s\n", p.Title)
				}
			}(f)
		}
		wg.Wait()
	}
}
