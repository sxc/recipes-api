// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/sxc/recipes-api.
//
//  Schemes: http
//  Host: localhost:8080
//  BasePath: /
//  Version: 1.0.0
//  Contact: Jim S
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta”

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"recepes-api/docs"
	"time"

	// gin-swagger middleware

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewRecipeHanler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHanler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHanler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	for i, r := range recipes {
		if r.ID == id {
			recipes[i] = recipe
			c.JSON(http.StatusOK, recipe)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "recipe not found"})
}

func DeleteREcipeHanler(c *gin.Context) {
	id := c.Param("id")
	for i, r := range recipes {
		if r.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "recipe deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "recipe not found"})
}

func SearchRecipesHanler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for _, r := range recipes {
		for _, t := range r.Tags {
			if t == tag {
				listOfRecipes = append(listOfRecipes, r)
			}
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

var recipes []Recipe

var ctx context.Context
var err error
var client *mongo.Client

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	json.Unmarshal([]byte(file), &recipes)

	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Cnnected to MongoDB")

	var listOfRecipes []interface{}
	for _, r := range recipes {
		listOfRecipes = append(listOfRecipes, r)
	}
	collection := client.Database("recipes").Collection("recipes")
	collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ", len(listOfRecipes), " recipes")
}

func main() {

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()
	router.POST("/recipes", NewRecipeHanler)
	router.GET("/recipes", ListRecipesHanler)
	router.PUT("/recipes/:id", UpdateRecipeHanler)
	router.DELETE("/recipes/:id", DeleteREcipeHanler)
	router.GET("/recipes/search", SearchRecipesHanler)
	router.Run()
}

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `josn:"instructions"`
	PublishedAt  time.Time `josn:"publishedAt"`
}
