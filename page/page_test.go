package page

import "testing"

func TestParseHTML(t *testing.T) {
	s := `
	<html>
	<body>
	<a href="https://www.google.com">search</a>
	<img src="https://images.google.com/puppy.jpg"/>
	</body>
	</html>
	`
	links, images := parseHTML(s)
	if len(links) != 1 {
		t.Errorf("Expected one link, but got %d", len(links))
	}
	if len(images) != 1 {
		t.Errorf("Expected one image, but got %d", len(images))
	}

	expected := []string{
		"https://www.google.com",
		"https://images.google.com/puppy.jpg"}

	if _, ok := links[expected[0]]; !ok {
		t.Errorf("Expected %s, but was not found", expected[0])
	}
	if _, ok := images[expected[1]]; !ok {
		t.Errorf("Expected %s, but was not found", expected[1])
	}

}
