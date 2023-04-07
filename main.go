package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	json.Unmarshal([]byte(file), &recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHanler)
	router.GET("/recipes", ListRecipesHanler)
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
