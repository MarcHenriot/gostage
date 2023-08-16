package main

import (
	"net/http"
	"os"

	"github.com/MarcHenriot/gostage/pkg/html"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const (
	TEMPLATE_FILE_PATH = "examples/template/template.yaml"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("ui/templates/*")
	r.Static("/static", "ui/static")
	r.StaticFile("/favicon.svg", "ui/static/favicon.svg")

	data, _ := os.ReadFile(TEMPLATE_FILE_PATH)

	var schema html.FormSchema
	yaml.Unmarshal(data, &schema)

	r.GET("/", func(c *gin.Context) {
		formHtml := html.NewForm()
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.HTML(200, "index.html", gin.H{
			"FormHTML": formHtml.Generate(schema),
		})
	})

	r.POST("/submit", func(c *gin.Context) {
		// Capture the form data
		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		age := c.PostForm("age")

		// Return as JSON
		c.JSON(http.StatusOK, gin.H{
			"firstName": firstName,
			"lastName":  lastName,
			"age":       age,
		})
	})

	r.Run(":8080")
}
