package rss

import (
	"strings"
	"testing"
)

const sampleRSS = `
<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd"
     xmlns:podcast="https://podcastindex.org/namespace/1.0">
  <channel>
    <title>Test</title>

	<itunes:new-feed-url>https://feeds.podcastics.com/podcastics/podcasts/rss/7543_22adb373093deb54e1ec644c2a7adec7.rss</itunes:new-feed-url>

    <item>
      <title>KEEP ME</title>
      <description><![CDATA[<p>Hello <strong>world</strong></p>]]></description>
      <enclosure url="https://example.com/audio.mp3" type="audio/mpeg"/>
      <itunes:image href="https://example.com/image.png"/>
      <podcast:person role="host">Alice</podcast:person>
    </item>

    <item>
      <title>DROP ME</title>
      <description><![CDATA[<p>Should not appear</p>]]></description>
    </item>

  </channel>
</rss>
`

func TestFilterXMLPreservesContent(t *testing.T) {
	keep := map[string]bool{
		"KEEP ME": true,
	}

	out, err := FilterXML([]byte(sampleRSS), keep)
	if err != nil {
		t.Fatalf("FilterXML error: %v", err)
	}

	s := string(out)

	// kept item present
	if !strings.Contains(s, "<title>KEEP ME</title>") {
		t.Fatal("kept item missing")
	}

	// dropped item gone
	if strings.Contains(s, "<title>DROP ME</title>") {
		t.Fatal("dropped item still present")
	}

	// content preserved (HTML content still there, escaped or not)
	if !strings.Contains(s, "Hello") || !strings.Contains(s, "world") {
		t.Fatal("description content missing")
	}

	// enclosure preserved
	if !strings.Contains(s, "<enclosure") {
		t.Fatal("enclosure missing")
	}

	// itunes:image preserved
	if !strings.Contains(s, "itunes:image") {
		t.Fatal("itunes:image missing")
	}

	// podcast namespace preserved
	if !strings.Contains(s, "podcast:person") {
		t.Fatal("podcast:person missing")
	}
}

func TestFilterXMLRewriteNewFeedURL(t *testing.T) {
	keep := map[string]bool{
		"KEEP ME": true,
	}

	const rewritten = "https://podcasts.decre.me/rss/bible-en-un-an.xml"

	out, err := FilterXMLWithOptions([]byte(sampleRSS), keep, FilterXMLOptions{RewriteNewFeedURL: rewritten})
	if err != nil {
		t.Fatalf("FilterXMLWithOptions error: %v", err)
	}

	s := string(out)

	if !strings.Contains(s, "<itunes:new-feed-url>"+rewritten+"</itunes:new-feed-url>") {
		t.Fatal("expected itunes:new-feed-url to be rewritten")
	}
	if strings.Contains(s, "feeds.podcastics.com/podcastics/podcasts/rss/") {
		t.Fatal("expected upstream itunes:new-feed-url to be removed/replaced")
	}
}

func TestFilterXMLNoRewriteKeepsUpstreamNewFeedURL(t *testing.T) {
	keep := map[string]bool{
		"KEEP ME": true,
	}

	out, err := FilterXMLWithOptions([]byte(sampleRSS), keep, FilterXMLOptions{})
	if err != nil {
		t.Fatalf("FilterXMLWithOptions error: %v", err)
	}

	s := string(out)
	if !strings.Contains(s, "<itunes:new-feed-url>https://feeds.podcastics.com/podcastics/podcasts/rss/7543_22adb373093deb54e1ec644c2a7adec7.rss</itunes:new-feed-url>") {
		t.Fatal("expected upstream itunes:new-feed-url to be preserved when no rewrite requested")
	}
}
