package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"html/template"
)

func generateFormFromSchema(schema map[string]interface{}) template.HTML {
	properties := schema["properties"].(map[string]interface{})
	formHtml := `<form hx-post="/submit" hx-swap="outerHTML">`
	for key, valueMap := range properties {
		value := valueMap.(map[string]interface{})
		formHtml += `<label>` + value["title"].(string) + `</label>`
		switch value["type"].(string) {
		case "string":
			formHtml += `<input type="text" name="` + key + `">`
		case "integer":
			formHtml += `<input type="number" name="` + key + `">`
		}
	}
	formHtml += `<button type="submit">Submit</button></form>`
	return template.HTML(formHtml)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("ui/templates/*")
	r.Static("/static", "ui/static")
	r.StaticFile("/favicon.svg", "ui/static/favicon.svg")


	schema := map[string]interface{}{
		"title":       "A registration form",
		"description": "A simple form example.",
		"required":    []string{"firstName", "lastName"},
		"type":        "object",
		"properties": map[string]interface{}{
			"firstName": map[string]interface{}{
				"type":  "string",
				"title": "First name",
			},
			"lastName": map[string]interface{}{
				"type":  "string",
				"title": "Last name",
			},
			"age": map[string]interface{}{
				"type":  "integer",
				"title": "Age",
			},
			"phone": map[string]interface{}{
				"type":  "string",
				"title": "Phone",
			},
		},
	}

	r.GET("/", func(c *gin.Context) {
		formHtml := generateFormFromSchema(schema)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.HTML(200, "index.html", gin.H{
			"FormHTML": formHtml,
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
			"lastName": lastName,
			"age": age,
		})
	})

	r.Run(":8080")
}
