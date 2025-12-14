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
