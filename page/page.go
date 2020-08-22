package page

import (
	"net/url"
	"strings"
	"sync"
)

// DownloadState defines the state of download progress
type DownloadState int

// The following enum defines the various stages of
// download progress
const (
	DownloadNotStarted DownloadState = iota
	DownloadStarted
	DownloadFinished
	DownloadFailed
)

// Page defines an HTML page and its content
type Page struct {
	URL    *url.URL
	HTML   string
	links  []*Page
	images []string
	state  DownloadState
}

// URLProcessor Processes a given url
func URLProcessor(domain *url.URL, sm *sync.Map, URL string) *Page {
	if sm == nil {
		return nil
	}
	URL = strings.TrimSpace(URL)
	if URL == "" {
		return nil
	}

	link := validateLink(domain, URL)
	if link == nil {
		return nil
	}

	var page *Page
	if p, ok := sm.Load(URL); ok {
		// found already
		page = p.(*Page)
		return page // return existing page
	}

	page = &Page{URL: link, state: DownloadNotStarted}
	// record this page in registry
	sm.Store(link.String(), page)

	// download page and extract
	html, err := downloadLink(URL)
	if err != nil {
		page.state = DownloadFailed
		// TODO: retry logic
		return page
	}
	page.state = DownloadFinished
	page.HTML = html

	// Extract links and images
	links, images := parseHTML(html)
	for im := range images {
		page.images = append(page.images, im)
	}

	// concurrently process the links
	m := &sync.Mutex{}
	for li := range links {
		go func(pp *Page, link string) {
			p := URLProcessor(domain, sm, link)
			m.Lock()
			pp.links = append(pp.links, p)
			m.Unlock()
		}(page, li)
	}

	return page
}
