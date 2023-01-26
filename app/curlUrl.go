package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	res, err := http.Get("https://qiita.com/sayama0402/items/e32814e38375fafa919a")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal("Unexpected Statuscode:", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if !(strings.Contains(contentType, "utf") || (strings.Contains(contentType, "UTF"))) {
		log.Fatal("Unexpected Content-Type:", contentType)
	}

	byteArray, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	siteTitle := ""
	if strings.Contains(string(byteArray), "title") {
		r := regexp.MustCompile(`<title.*?>(.*?)</title>`)
		match := r.FindStringSubmatch(string(byteArray))
		if len(match) > 1 {
			siteTitle = match[1]
		}
	}
	fmt.Println(siteTitle)
}
