package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// index, err := ioutil.ReadFile("public/index.html")
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "test",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
