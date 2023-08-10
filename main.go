package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	TEMPLATE_FOLDER    = "examples/template"
	DESTINATION_FOLDER = "examples/destination"
	JSON_DATA_FILE     = "examples/template/data.json"
)

type TemplateData struct {
	Values map[string]interface{}
}

func main() {
	var templateData TemplateData
	
	content, err := os.ReadFile(JSON_DATA_FILE)
	if err != nil {
		slog.Error("Error when opening file: ", "err", err)
	}
	err = json.Unmarshal(content, &templateData.Values)
	if err != nil {
		slog.Error("Error during Unmarshal(): ", "err", err)
	}

	templater := template.New("templater").Delims("${{", "}}")

	dirs, files, _ := listFilesAndDirs(TEMPLATE_FOLDER)
	if err != nil {
		slog.Error("Error getting pathes: ", "err", err)
	}
	slog.Info("Listing pathes", "files", files, "dirs", dirs)

	for _, dir := range dirs {
		templatedDir, _ := templatePath(dir, templater, templateData)
		os.Mkdir(filepath.Join(DESTINATION_FOLDER, templatedDir), os.ModePerm)
	}
	for _, inputPath := range files {
		outputPath, _ := templatePath(inputPath, templater, templateData)
		
		outputFile, _ := os.Create(filepath.Join(DESTINATION_FOLDER, outputPath))
		defer outputFile.Close()
		
		inputFile, _ := os.ReadFile(inputPath)
		
		parsedTemplate, _ := templater.Parse(string(inputFile))
		parsedTemplate.Execute(outputFile, templateData)
	}
}

func templatePath(path string, templater *template.Template, data TemplateData) (string, error) {
	var buffer bytes.Buffer
	parsedTemplate, err := templater.Parse(path)
	err = parsedTemplate.Execute(&buffer, data)
	return strings.TrimPrefix(buffer.String(), TEMPLATE_FOLDER), err
}

func listFilesAndDirs(rootPath string) (files []string, dirs []string, err error) {
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		} else {
			files = append(files, path)
		}
		return nil
	})
	return dirs, files, err
}
