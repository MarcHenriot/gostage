package html

import (
	"bytes"
	"html/template"
	"regexp"
	"slices"
)

type Form struct {
	template       *template.Template
	stringTemplate string
	funcMap        template.FuncMap
}

func NewForm() *Form {
	funcMap := template.FuncMap{
		"markdownToHTML": markdownToHTML,
		"toHTML":         func(s string) template.HTML { return template.HTML(s) },
		"contains":       slices.Contains[[]string],
	}
	return &Form{
		template: template.New("formTemplate").Funcs(funcMap),
		stringTemplate: `
			<form hx-post="/submit" hx-trigger="submit" hx-target="this" hx-swap="outerHTML">
				{{- range $index, $parameter := .Parameters }}
				<fieldset>
					<legend> {{ $parameter.Title }}:: </legend>
					{{- range $key, $propertie := $parameter.Properties }}
						<label for="{{ $key }}">{{ $propertie.Title }}:</label><br>
						<input {{ if contains $parameter.Required $key }}required{{ end }} value="{{ $propertie.Default }}" {{ if $propertie.Autofocus }}autofocus{{ end }} type="{{ $propertie.Type }}" id="{{ $key }}" name="{{ $key }}">
						<p>{{ if $propertie.EnableMarkdownInDescription }}{{ markdownToHTML $propertie.Description | toHTML }}{{ else }}{{ $propertie.Description | toHTML }}{{ end }}</p>
					{{- end }}
				</fieldset>
				{{- end }}
				<input class="mui-button" type="submit" value="Submit">
			</form>`,
		funcMap: template.FuncMap{
			"markdownToHTML": markdownToHTML,
			"toHTML":         func(s string) template.HTML { return template.HTML(s) },
			"contains":       slices.Contains[[]string],
		},
	}
}

func (f *Form) Generate(formSchema FormSchema) template.HTML {
	var buf bytes.Buffer
	tmpl, _ := f.template.Parse(f.stringTemplate)
	tmpl.Execute(&buf, formSchema)
	return template.HTML(buf.String())
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
