package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mongodb/mongo-go-driver/mongo"
)

const (
	signature = "drowssap"
)

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type accessToken struct {
	Token     string `json:"accessToken"`
	ExpiresIn int64  `json:"expiresIn"`
}

type signupName struct {
	Username string
	Password string
	Email    string
}

func signup(c echo.Context) error {
	var m signupName
	err := c.Bind(&m)
	if err != nil {
		c.Error(err)
	}

	if m.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "email is required",
		})
	}

	sum := sha256.Sum256([]byte(m.Password))
	m.Password = fmt.Sprintf("%x", sum)

	err = mgoInsert(m)
	if err != nil {
		c.Error(err)
	}

	token, err := generateToken()
	if err != nil {
		c.Error(err)
	}

	return c.JSON(http.StatusOK, accessToken{
		Token:     token,
		ExpiresIn: int64(time.Hour.Minutes()),
	})
}

func generateToken() (string, error) {
	mySigningKey := []byte(secret)

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Issuer:    "odds",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func mongodbInsert(data signupName) error {
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		return err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}
	collection := client.Database("odds").Collection("credential")
	_, err = collection.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

func mgoInsert(data signupName) error {
	url := "mongodb://localhost:27017"
	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	col := session.DB("odds").C("credential")
	_, err = col.Upsert(bson.M{"email": data.Email}, data)
	return err
}

var (
	secret string
)

func init() {
	s := flag.String("secret", signature, "-scecret=yourpassword")
	secret = *s
}

func main() {
	flag.Parse()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/signup", signup)
	e.GET("/posts/:id", getPostsHandler, middleware.JWT([]byte(secret)))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func getPostsHandler(c echo.Context) error {
	p, err := getPosts()
	if err != nil {
		c.Error(err)
	}

	return c.JSON(http.StatusOK, p)
}

func getPosts() ([]comment, error) {
	url := "https://jsonplaceholder.typicode.com/posts"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var comments []comment
	err = json.Unmarshal(body, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

type comment struct {
	UserID int64  `json:"userId"`
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
