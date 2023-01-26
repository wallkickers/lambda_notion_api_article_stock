package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

type Response struct {
	RequestBody string `json:"RequestBody"`
}

type Event struct {
	Events []struct {
		Message struct {
			Text string `json:"text"`
		} `json:"message"`
		Source struct {
			UserID string `json:"userId"`
		} `json:"source"`
	}
}

func handler(request events.APIGatewayProxyRequest) {
	var event Event

	fmt.Println(request)
	fmt.Println(request.Body)

	body := request.Body
	res := Response{
		RequestBody: body,
	}

	json.Unmarshal([]byte(res.RequestBody), &event)

	userid := fmt.Sprintf("%v", event.Events[0].Source.UserID)
	text := fmt.Sprintf("%v", event.Events[0].Message.Text)

	siteTitle := httpGetUrl(text)
	isApiSuccess := postNotionApiStockArticle(siteTitle, text)
	if isApiSuccess {
		postLineMessage(userid, text)
	} else {
		postLineMessage(userid, "記事の保存に失敗しました。")
	}
}

// パラメータのURLでcurlを叩き、サイトのタイトルを取得&返却
func httpGetUrl(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal("Unexpected Statuscode:", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "utf") {
		log.Fatal("Unexpected Content-Type:", contentType)
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
	return siteTitle
}

// パラメータのタイトルとURLでNotionAPIを叩く
func postNotionApiStockArticle(siteTitle string, url string) bool {
	jsonBody, err := os.ReadFile("postNotionApiStockArticleRequest.json")
	if err != nil {
		fmt.Println(err)
	}
	relacedJsonBody := replaceParameter(string(jsonBody), siteTitle, url)

	body := strings.NewReader(string(relacedJsonBody))
	req, err := http.NewRequest("POST", "https://api.notion.com/v1/pages", body)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	fmt.Println("------------------------------------------")
	fmt.Println(req)
	fmt.Println("------------------------------------------")
	fmt.Println(response)

	if response.StatusCode != 200 {
		return false
	}
	return true
}

func replaceParameter(jsonBody, content, url string) string {
	replacement := map[string]string{
		"%NOTION_DATABASE_ID%": os.Getenv("NOTION_DATABASE_ID"),
		"%CONTENT%":            content,
		"%SITEURL%":            url,
	}
	for key, value := range replacement {
		jsonBody = strings.Replace(jsonBody, key, value, -1)
	}
	return jsonBody
}

func postLineMessage(userid string, text string) {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}

	if _, err := bot.PushMessage(userid, linebot.NewTextMessage(text)).Do(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	lambda.Start(handler)
}
