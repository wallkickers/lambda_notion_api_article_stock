package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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

	byteArray, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	siteTitle := ""
	if strings.Contains(string(byteArray), "<title>") {
		splitArray := strings.Split(string(byteArray), "<title>")
		splitArray = strings.Split(splitArray[1], "</title>")
		siteTitle = splitArray[0]
	}
	fmt.Println(siteTitle)
}
