package server

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

type GostageServer struct {
	router *gin.Engine
}

func New() *GostageServer {
	return &GostageServer{
		router: gin.Default(),
	}
}

func (gs *GostageServer) Init() {
	data, _ := os.ReadFile(TEMPLATE_FILE_PATH)
	var schema html.FormSchema
	yaml.Unmarshal(data, &schema)

	gs.router.LoadHTMLGlob("ui/templates/*")
	gs.router.Static("/static", "ui/static")
	gs.router.StaticFile("/favicon.svg", "ui/static/favicon.svg")
	gs.router.GET("/", getForm(schema))
	gs.router.POST("/submit", submitForm())
}

func (gs *GostageServer) Run() {
	gs.router.Run(":8080")
}

func getForm(schema html.FormSchema) gin.HandlerFunc {
	return func(c *gin.Context) {
		formHtml := html.NewForm()
		c.Set("renderUI", schema)
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.HTML(200, "index.html", gin.H{
			"FormHTML": formHtml.Generate(schema),
		})
	}
}

func submitForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		firstName := c.PostForm("firstName")
		lastName := c.PostForm("lastName")
		age := c.PostForm("age")
		c.JSON(http.StatusOK, gin.H{
			"firstName": firstName,
			"lastName":  lastName,
			"age":       age,
		})
	}
}
