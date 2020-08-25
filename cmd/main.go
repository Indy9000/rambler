package main

import (
	"log"
	"net/url"
	"sync"

	"github.com/indy9000/rambler/page"
)

var registry sync.Map

func main() {
	registry = sync.Map{}
	domain, _ := url.Parse("https://www.cuvva.com/")
	page.URLProcessor(domain, &registry, domain.String())
	log.Printf("Done!")
}
