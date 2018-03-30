package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"
)

var (
	WEBDONA_BOT_TOKEN       string
	WEBDONA_CONTACT_CHAT_ID string
)

type Data struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Body  string `json:"body"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	data := Data{}

	if err := json.Unmarshal([]byte(request.Body), &data); err != nil {

		log.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil

	}

	text := fmt.Sprintf("<b>%s</b><pre>\n</pre><b>%s</b><pre>\n</pre><pre>%s</pre>", data.Name, data.Email, data.Body)

	if req, err := http.NewRequest("GET", fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", WEBDONA_BOT_TOKEN), nil); err != nil {

		log.Println(err)

		return events.APIGatewayProxyResponse{
			StatusCode: 503,
		}, nil

	} else {

		q := req.URL.Query()

		q.Add("text", text)

		q.Add("chat_id", WEBDONA_CONTACT_CHAT_ID)

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
