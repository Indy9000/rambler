package page

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func parseHTML(htmlContent string) (map[string]bool, map[string]bool) {
	links := map[string]bool{}
	images := map[string]bool{}
	htmlContent = strings.TrimSpace(htmlContent)
	if htmlContent == "" {
		return links, images
	}
	reader := strings.NewReader(htmlContent)
	node, err := html.Parse(reader)
	if err != nil {
		// TODO handle error
		return links, images
	}

	extractor := func(n *html.Node, k string, m map[string]bool) {
		for _, a := range n.Attr {
			if a.Key == k {
				log.Printf("found : %s", a.Val)
				m[a.Val] = true
			}
		}
	}

	var nodeNavigator func(n *html.Node)
	nodeNavigator = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "a" { // extract link
				extractor(n, "href", links)
			}
			if n.Data == "img" { // extract image URL
				extractor(n, "src", images)
			}
		}
		// depth first search
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			nodeNavigator(c)
		}
	}

	nodeNavigator(node)
	return links, images
}

func downloadLink(URL string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf("Unable to download %s. Error: %s", URL, err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read downloaded %s. Error: %s", URL, err.Error())
		return "", err
	}
	return string(body), nil
}

// validateLink checks if the URL given is inside the domain
// and returns the fully qualified url.URL object
func validateLink(domain *url.URL, URL string) *url.URL {
	u, e := url.Parse(URL)
	if e != nil {
		return nil
	}
	if u.IsAbs() {
		if domain.Host == u.Host { // inside this domain
			return u
		}
		return nil
	}

	l, e := domain.Parse(URL)
	if e != nil {
		log.Println(e)
		return nil
	}
	return l
}
