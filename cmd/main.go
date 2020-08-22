package main

import (
	"sync"
)

var registry sync.Map

func main() {
	registry = sync.Map{}
}
