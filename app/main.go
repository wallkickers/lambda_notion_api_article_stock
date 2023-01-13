package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	}
}

// パラメータのURLでcurlを叩き、サイトのタイトルを取得&返却
func httpGetUrl(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

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
	return true
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
