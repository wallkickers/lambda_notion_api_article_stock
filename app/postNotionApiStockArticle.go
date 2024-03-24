package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type PostNotionApiStockArticleRequestData struct {
	Parent struct {
		DatabaseID string
	} `json:"parent"`
	Properties struct {
		Title struct {
			Title []struct {
				Text struct {
					Content string `json:"content"`
				} `json:"text"`
			} `json:"title"`
		} `json:"title"`
		Status struct {
			Select struct {
				Name string `json:"name"`
			} `json:"select"`
		} `json:"status"`
		URL struct {
			Url string `json:"URL"`
		} `json:"URL"`
	} `json:"properties"`
}

func main() {
	jsonBody, err := os.ReadFile("postNotionApiStockArticleRequest.json")
	if err != nil {
		fmt.Println(err)
	}
	relacedJsonBody := replaceParameter(string(jsonBody), "置換後content", "置換後url")

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
