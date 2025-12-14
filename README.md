# rss-proxy

A lightweight RSS proxy for podcasts that allows filtering episodes
based on declarative rules without modifying the original audio or metadata.

Designed to be fully compatible with Apple Podcasts.
It will probably work with other readers too.

---

## Features

* RSS item filtering without XML re-serialization
* Apple Podcasts–safe (no enclosure rewriting)
* Declarative YAML configuration
* HTTP cache with ETag / If-Modified-Since support
* Fully testable, no external dependencies

---

## Non-goals

* ❌ Audio processing
* ❌ RSS transformation or rewriting
* ❌ Podcast analytics
* ❌ Stateful episode aggregation

---

## Example use cases

* Keep episodes from a podcast with a specific regex in the title
* Keep only the last episode of multipart stories (`[x/y]`)
* Skip early episodes of long-running podcasts
* Create personal filtered feeds

---

## Configuration

```yaml
feeds:
  - id: legend-rediff
    source: https://example.com/feed.xml
    rules:
      - type: title_contains
        value: "[REDIFF]"
```

---

## Supported rules

| Rule                    | Description                                   |
| ----------------------- | --------------------------------------------- |
| `title_contains`        | Keep episodes whose title contains a string   |
| `title_fraction_equals` | Keep only episodes where `[x/y]` and `x == y` |
| `episode_number_min`    | Keep episodes with episode number ≥ N         |

---

## Running locally

```bash
make run
```

---

## Docker

```bash
make docker-run
```

---

## Apple Podcasts

Use the generated feed URLs directly in Apple Podcasts.
The proxy does not modify audio files.

---

## Disclaimer

This project is intended for personal use.
Users are responsible for complying with podcast distribution rights.

---

## License

MIT

```text
Copyright © 2025 Amaury Decrême

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
