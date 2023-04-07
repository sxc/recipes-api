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

func init() {
	recipes = make([]Recipe, 0)
	file, _ := ioutil.ReadFile("recipes.json")
	json.Unmarshal([]byte(file), &recipes)
}

func main() {
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
