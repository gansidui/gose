package extractutil

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestExtract(t *testing.T) {
	content, err := ioutil.ReadFile("./test.html")
	if err != nil {
		log.Fatal(err)
	}

	title := ExtractTitle(string(content))
	ioutil.WriteFile("./test_title.txt", []byte(title), os.ModePerm)

	body := ExtractBody(string(content))
	ioutil.WriteFile("./test_body.txt", []byte(body), os.ModePerm)
}
