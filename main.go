package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	sessionKey = "session_id"
)

var (
	sessionMap sync.Map
)

//User user
type User struct {
	Name     string `json:"name" dynamodbav:"name"`
	Password string `json:"password" dynamodbav:"password"`
}

func helloHandler(c *gin.Context) {
	name := c.Param("name")
	c.String(http.StatusOK, "Hello %s", name)
}

func welcomeHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello World from Go")
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
	user := User{
		Name: request.Name,
	}
	sessionID, err := c.Cookie(sessionKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			sessionID = uuid.New().String()
			c.SetCookie(sessionKey, sessionID, 86400,
				"/", "", false, true)
		} else {
			c.Error(fmt.Errorf("unexpect error occurs, request: %+v, err: %s", c.Request, err.Error()))
		}
	}
	sessionMap.Store(sessionID, &user)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    "ok",
	})
}

func listUserHandle(c *gin.Context) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}
	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("user-oa"),
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out.Items)
	userList := make([]User, 0)
	err = attributevalue.UnmarshalListOfMaps(out.Items, &userList)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal Dynamodb Scan Items, %v", err))
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    userList,
	})
}

func getUserHandle(c *gin.Context) {
	var user *User
	sessionID, err := c.Cookie(sessionKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "error",
			"data":    nil,
		})
		return
	}
	if v, ok := sessionMap.Load(sessionID); ok {
		user = v.(*User)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    user,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"message": "error",
		"data":    nil,
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
	r.GET("/api/get_user", getUserHandle)

	return r
}

func main() {
	addr := ":" + os.Getenv("PORT")
	log.Fatal(gateway.ListenAndServe(addr, routerEngine()))
}
