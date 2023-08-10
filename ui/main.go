package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const (
	TEMPLATE_FILE_PATH = "examples/template/template.yaml"
)

type FormSchema struct {
	Title       string              `yaml:"title"`
	Description string              `yaml:"description"`
	Type        string              `yaml:"type"`
	Required    []string            `yaml:"required"`
	Properties  map[string]Property `yaml:"properties"`
}

type Property struct {
	Type                        string `yaml:"type"`
	Title                       string `yaml:"title"`
	Default                     string `yaml:"default,omitempty"`
	Autofocus                   bool   `yaml:"ui:autofocus,omitempty"`
	EmptyValue                  string `yaml:"ui:emptyValue,omitempty"`
	Placeholder                 string `yaml:"ui:placeholder,omitempty"`
	Autocomplete                string `yaml:"ui:autocomplete,omitempty"`
	EnableMarkdownInDescription bool   `yaml:"ui:enableMarkdownInDescription,omitempty"`
	Description                 string `yaml:"ui:description,omitempty"`
	Widget                      string `yaml:"ui:widget,omitempty"`
	UiTitle                     string `yaml:"ui:title,omitempty"`
}

func markdownToHTML(input string) string {
	// Convert **bold** to <b>bold</b>
	reBold := regexp.MustCompile(`\*\*(.*?)\*\*`)
	input = reBold.ReplaceAllString(input, "<b>$1</b>")

	// Convert *italic* to <i>italic</i>
	reItalic := regexp.MustCompile(`\*(.*?)\*`)
	input = reItalic.ReplaceAllString(input, "<i>$1</i>")

	return input
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}


func generateHTML(formSchema FormSchema) template.HTML {
	formHTML := `
		<div class="mui-form">
			<h1 class="mui-container">{{.Title}}</h1>
			<p class="mui-container">{{.Description}}</p>
			<form hx-post="/submit" hx-trigger="submit" hx-target="this" hx-swap="outerHTML">
				{{range $key, $value := .Properties}}
					<div class="mui-card">
						<label>{{$value.Title}}</label>
						<p>{{if $value.EnableMarkdownInDescription}}{{markdownToHTML $value.Description | toHTML}}{{else}}{{$value.Description | toHTML}}{{end}}</p>
						<input type="{{$value.Type}}" 
							name="{{$key}}" 
							placeholder="{{$value.Placeholder}}" 
							{{if $value.Autofocus}} autofocus {{end}}
							autocomplete="{{$value.Autocomplete}}"
							value="{{$value.Default}}"
							{{if contains $.Required $key}}required{{end}}>
					</div>
				{{end}}
				<input class="mui-button" type="submit" value="Submit">
			</form>
		</div>
	`	

	funcMap := template.FuncMap{
		"markdownToHTML": markdownToHTML,
		"toHTML":         func(s string) template.HTML { return template.HTML(s) },
		"contains":         contains,
	}

	tmpl, err := template.New("formTemplate").Funcs(funcMap).Parse(formHTML)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, formSchema)
	if err != nil {
		panic(err)
	}

	return template.HTML(buf.String())
}

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("ui/templates/*")
	r.Static("/static", "ui/static")
	r.StaticFile("/favicon.svg", "ui/static/favicon.svg")

	data, _ := os.ReadFile(TEMPLATE_FILE_PATH)

	var schema FormSchema
	yaml.Unmarshal(data, &schema)

	r.GET("/", func(c *gin.Context) {
		formHtml := generateHTML(schema)
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
			"lastName":  lastName,
			"age":       age,
		})
	})

	r.Run(":8080")
}
