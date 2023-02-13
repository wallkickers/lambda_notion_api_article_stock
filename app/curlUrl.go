package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func main() {
	messageText := "https://qiita.com/sayama0402/items/e32814e38375fafa919a"

	// LINEで転送された際にメッセージ冒頭に改行が入るため削除
	messageTextRemovedNewLine := strings.ReplaceAll(messageText, "\n", "")

	// urlチェック
	parseUrl, _ := url.Parse(messageTextRemovedNewLine)
	inputUrl := strings.Join(strings.Fields(parseUrl.String()), "")

	res, err := http.Get(inputUrl)
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
		r := regexp.MustCompile(`(?s)<title.*?>(.*?)</title>`)
		match := r.FindStringSubmatch(string(byteArray))
		if len(match) > 1 {
			siteTitle = strings.Join(strings.Fields(match[1]), "")
		}
	}
	fmt.Println(siteTitle)
}
