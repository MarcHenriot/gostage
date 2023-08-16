package html

type FormSchema struct {
	Title       string      `yaml:"title"`
	Description string      `yaml:"description"`
	Parameters  []Parameter `yaml:"parameters"`
}

type Parameter struct {
	Title      string               `yaml:"title"`
	Required   []string             `yaml:"required"`
	Properties map[string]Propertie `yaml:"properties"`
}

type Propertie struct {
	Title                       string `yaml:"title"`
	Type                        string `yaml:"type"`
	Default                     string `yaml:"default,omitempty"`
	Autofocus                   bool   `yaml:"ui:autofocus,omitempty"`
	EnableMarkdownInDescription bool   `yaml:"ui:enableMarkdownInDescription,omitempty"`
	Description                 string `yaml:"ui:description,omitempty"`
}
