package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"runtime/debug"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type body struct {
	PageUrl string `json:"pageUrl"`
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()
	var b body
	err := json.Unmarshal([]byte(request.Body), &b)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Unable to parse body",
		}, nil
	}

	switch request.HTTPMethod {
	case "POST":
		var ingredientsList []string

		res, err := http.Get(b.PageUrl)
		if err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Unable to get page",
			}, nil
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Unable to create document",
			}, nil
		}
		uls := doc.Find("ul")

		uls.Each(func(i int, s *goquery.Selection) {
			itemsWithFirstAsDigit := 0
			s.Find("li").Each(func(ii int, ss *goquery.Selection) {
				trimmed := strings.TrimSpace(ss.Text())
				liRune := []rune(trimmed)
				fmt.Println(string(liRune))
				if trimmed != "" {
					if unicode.IsDigit(liRune[0]) {
						itemsWithFirstAsDigit++
					}
				}
			})
			if itemsWithFirstAsDigit >= 3 {
				fmt.Println("FOUND INGREDIENT LIST")
				s.Find("li").Each(func(ii int, ss *goquery.Selection) {
					trimmed := strings.ReplaceAll(ss.Text(), "\t", "")
					withoutSpaces := strings.ReplaceAll(trimmed, "\n", " ")
					ingredientsList = append(ingredientsList, withoutSpaces)
				})
			}
		})
		ingList, err := json.Marshal(ingredientsList)
		if err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Unable to marshal ingredients",
			}, nil
		}

		return &events.APIGatewayProxyResponse{
			StatusCode:        200,
			Headers:           map[string]string{"Content-Type": "application/json"},
			MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
			Body:              string(ingList),
		}, nil
	default:
		return &events.APIGatewayProxyResponse{
			StatusCode: 404,
			Headers:    map[string]string{"Content-Type": "text/plain"},
			Body:       "Incorrect method",
		}, nil
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
