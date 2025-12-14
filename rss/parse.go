package rss

import "encoding/xml"

// Parse lit le RSS pour appliquer les règles,
// sans jamais être utilisé pour la sortie XML.
func Parse(data []byte) (RSS, error) {
	var feed RSS
	err := xml.Unmarshal(data, &feed)
	return feed, err
}
