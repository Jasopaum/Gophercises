package htmlParser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ParseHtml(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := findLinks(doc)

	return links, nil
}

func findLinks(n *html.Node) []Link {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []Link{buildLinkFromNode(n)}
	}
	var ret []Link
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, findLinks(c)...)
	}
	return ret
}

func buildLinkFromNode(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
		}
	}
	ret.Text = linkText(n)
	return ret
}

func linkText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += linkText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
