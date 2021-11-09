package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

	"runtime/debug"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	f "github.com/fauna/faunadb-go/v4/faunadb"
	"github.com/joho/godotenv"
)

type body struct {
	PageUrl string `json:"pageUrl"`
}

type Recipe struct {
	Ingredients []string `json:"ingredients"`
	Title       string   `json:"title"`
	Url         string   `json:"url"`
	Cached      bool     `json:"cached"`
}

func ErrorEvent(status int, message string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       message,
	}
}

func JsonEvent(status int, json string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       json,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func getTitle(document *goquery.Document) string {
	dividersRegexp := regexp.MustCompile(`-|\||by`)
	title := strings.TrimSpace(dividersRegexp.Split(document.Find("head title").Text(), -1)[0])
	return title
}

func getIngredients(document *goquery.Document) []string {
	var ingredients []string
	uls := document.Find("ul")

	uls.Each(func(i int, s *goquery.Selection) {
		itemsWithFirstAsDigit := 0
		s.Find("li").Each(func(ii int, ss *goquery.Selection) {
			trimmed := strings.TrimSpace(ss.Text())
			trimmedSr := strings.TrimSpace(strings.ReplaceAll(trimmed, "▢", ""))
			liRune := []rune(trimmedSr)
			if trimmed != "" {
				if unicode.IsDigit(liRune[0]) {
					itemsWithFirstAsDigit++
				}
			}
		})
		if itemsWithFirstAsDigit >= 3 {
			s.Find("li").Each(func(ii int, ss *goquery.Selection) {
				trimmed := strings.ReplaceAll(ss.Text(), "\t", "")
				withoutSpaces := strings.ReplaceAll(trimmed, "\n", " ")
				trimmedSr := strings.TrimSpace(strings.ReplaceAll(withoutSpaces, "▢", ""))
				ingredients = append(ingredients, trimmedSr)
			})
		}
	})

	return removeDuplicateStr(ingredients)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	godotenv.Load(".env")

	var b body
	err := json.Unmarshal([]byte(request.Body), &b)
	if err != nil {
		return ErrorEvent(400, "Unable to parse body"), nil
	}

	switch request.HTTPMethod {
	case "POST":
		client := f.NewFaunaClient(
			os.Getenv("FAUNA_SECRET"),
			f.Endpoint("https://db.eu.fauna.com"),
		)

		var response Recipe

		recipeFromDB, dbErr := client.Query(f.Get(f.MatchTerm(f.Index("recipes_search_by_url"), b.PageUrl)))
		if dbErr != nil {
			// The recipe is not in the db, so fetch from scratch
			fmt.Println("Not found in db")

			res, err := http.Get(b.PageUrl)

			if err != nil {
				return ErrorEvent(400, "Unable to get page"), nil
			}
			defer res.Body.Close()

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				return ErrorEvent(400, "Unable to create document"), nil
			}

			response.Title = getTitle(doc)
			response.Ingredients = getIngredients(doc)
			response.Url = b.PageUrl

			_, err = client.Query(f.Create(f.Collection("recipes"), f.Obj{
				"data": response,
			}))
			if err != nil {
				fmt.Println(err)
				return ErrorEvent(400, "Unable to add to database"), nil
			}
		} else {
			// Recipe found in db
			if err := recipeFromDB.At(f.ObjKey("data")).Get(&response); err != nil {
				return ErrorEvent(400, "Unable to marshal data into struct"), nil
			}
			response.Cached = true
		}

		resp, err := json.Marshal(response)
		if err != nil {
			return ErrorEvent(400, "Unable to marshal response"), nil

		}

		return JsonEvent(200, string(resp)), nil

	default:
		return ErrorEvent(404, "Incorrect method"), nil
	}
}

func main() {
	lambda.Start(handler)
}
