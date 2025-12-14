package rss

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	PubDate     string `xml:"pubDate"`
	Episode     int    `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd episode"`
	Duration    string `xml:"http://www.itunes.com/dtds/podcast-1.0.dtd duration"`
	Description string `xml:"description"`
}
