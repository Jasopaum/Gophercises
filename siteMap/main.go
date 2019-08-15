package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"gophercises/htmlParser"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var xmlDoctypeDeclaration = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"

func main() {
	var (
		urlPath  = flag.String("url", "https://gophercises.com/", "url of site you want to build the map of")
		maxDepth = flag.Int("depth", -2, "max depth to explore")
	)
	flag.Parse()

	reachable := bfs(*urlPath, *maxDepth)

	if err := convertToXml(reachable); err != nil {
		log.Println("Error when converting to xml:", err)
	}
}

type urlset struct {
	Locs  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}
type loc struct {
	Loc string `xml:"loc"`
}

func convertToXml(reachable map[string]struct{}) error {
	toXml := urlset{
		Locs:  make([]loc, len(reachable)),
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}
	i := 0
	for k, _ := range reachable {
		toXml.Locs[i] = loc{k}
		i++
	}

	writer, err := os.Create("toto.xml")
	if err != nil {
		return err
	}
	defer writer.Close()

	writer.Write([]byte(xmlDoctypeDeclaration))
	enc := xml.NewEncoder(writer)
	enc.Indent("", " ")
	if err := enc.Encode(toXml); err != nil {
		log.Println("Error while encoding:", err)
		return err
	}
	return nil
}

func bfs(url string, maxDepth int) map[string]struct{} {
	var (
		seen = map[string]struct{}{}
		nq   = map[string]struct{}{url: {}}
		q    map[string]struct{}
	)
	for d := 0; d != maxDepth+1 && len(nq) > 0; d++ {
		q, nq = nq, map[string]struct{}{}
		for u, _ := range q {
			if _, ok := seen[u]; ok {
				continue
			}
			seen[u] = struct{}{}
			urls, _ := getUrls(u)
			for _, newUrl := range urls {
				if _, ok := seen[newUrl]; !ok {
					nq[newUrl] = struct{}{}
				}
			}
		}
	}
	return seen
}

func getUrls(urlPath string) ([]string, error) {
	resp, err := http.Get(urlPath)
	if err != nil {
		log.Println("Error while getting the url:", err)
	}
	defer resp.Body.Close()

	base := getBaseUrl(resp)
	foundUrls, err := extractUrls(resp, base)
	if err != nil {
		fmt.Println("Error while parsing: ", err)
		return nil, err
	}
	return foundUrls, nil
}

func getBaseUrl(resp *http.Response) string {
	baseUrl := url.URL{
		Scheme: resp.Request.URL.Scheme,
		Host:   resp.Request.URL.Host,
	}
	base := baseUrl.String()

	return base
}

func extractUrls(resp *http.Response, base string) ([]string, error) {
	links, err := htmlParser.ParseHtml(resp.Body)
	if err != nil {
		return nil, err
	}

	var foundUrls []string
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			foundUrls = append(foundUrls, base+link.Href)
		case strings.HasPrefix(link.Href, base):
			foundUrls = append(foundUrls, link.Href)
		}
	}
	return foundUrls, nil
}
