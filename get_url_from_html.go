package main

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"net/url"
	"strings"
)

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}

	var links []string
	doc, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		log.Fatal(err)
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link := a.Val
					if !strings.HasPrefix(link, "http") {
						absoluteURL, err := baseURL.Parse(link)
						if err != nil {
							log.Println("Error parsing URL:", err)
							continue
						}
						link = absoluteURL.String()
					}
					links = append(links, link)
					break
				}
			}
		}
	}

	return links, nil
}
