package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
)

var (
	TELEGRAM_BOT_TOKEN string
	TELEGRAM_CHAT_ID   string
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var data url.Values

	var err error

	if data, err = url.ParseQuery(request.Body); err != nil {

		return events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil

	}

	if data.Get("honeypot") != "" {

		return events.APIGatewayProxyResponse{
			StatusCode: 403,
		}, nil

	}

	text := fmt.Sprintf("<b>%s</b><pre>\n</pre><b>%s</b><pre>\n</pre><pre>%s</pre>", data.Get("name"), data.Get("email"), data.Get("message"))

	if req, err := http.NewRequest("GET", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TELEGRAM_BOT_TOKEN), nil); err != nil {

		log.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil

	} else {

		q := req.URL.Query()

		q.Add("text", text)

		q.Add("chat_id", TELEGRAM_CHAT_ID)

		q.Add("parse_mode", "html")

		req.URL.RawQuery = q.Encode()

		client := http.Client{}

		if _, err := client.Do(req); err != nil {

			log.Println(err)

			return events.APIGatewayProxyResponse{
				StatusCode: 503,
			}, nil

		}

	}

	return events.APIGatewayProxyResponse{
		StatusCode: 204,
	}, nil

}

func main() {
	lambda.Start(handler)
}
