package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
)

//User user
type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func helloHandler(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}

func welcomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World from Go")
}

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to gin lambda server.",
	})
}

func registerHandle(c *gin.Context) {
	request := User{}
	c.BindJSON(&request)
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("user-oa"),
		Item: map[string]types.AttributeValue{
			"name":     &types.AttributeValueMemberS{Value: request.Name},
			"password": &types.AttributeValueMemberS{Value: request.Password},
		},
	})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(out.Attributes)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    "ok",
	})
}

func loginHandle(c *gin.Context) {
	request := User{}
	c.BindJSON(&request)
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	_, err = svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("user-oa"),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: request.Name},
		},
	})

	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    "ok",
	})
}

func listUserHandle(c *gin.Context) {
	request := User{}
	c.BindJSON(&request)
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("user-oa"),
		Item: map[string]types.AttributeValue{
			"name":     &types.AttributeValueMemberS{Value: request.Name},
			"password": &types.AttributeValueMemberS{Value: request.Password},
		},
	})

	if err != nil {
		log.Println(err)
	}

	fmt.Println(out.Attributes)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    "ok",
	})
}

func routerEngine() *gin.Engine {
	// set server mode
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/welcome", welcomeHandler)
	r.GET("/user/:name", helloHandler)
	r.StaticFile("/", "login_demo/index.html")

	r.POST("/api/register", registerHandle)
	r.POST("/api/login", loginHandle)
	r.POST("/api/list_user", listUserHandle)

	return r
}

func main() {
	addr := ":" + os.Getenv("PORT")
	log.Fatal(gateway.ListenAndServe(addr, routerEngine()))
}
