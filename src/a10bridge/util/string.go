package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

//ToJSON converts a entity into indented json string
func ToJSON(entity interface{}) string {
	json, err := json.MarshalIndent(entity, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(json)
}

//ApplyTemplate processes a string template using the provided data entity for lookups
func ApplyTemplate(data interface{}, tpl string) (string, error) {
	var url string
	tmpl, err := template.New(tpl).Parse(tpl)
	if err != nil {
		return url, err
	}

	var writer bytes.Buffer
	err = tmpl.Execute(&writer, data)
	url = writer.String()

	fmt.Println(url)

	return url, err
}

//Contains check if a string slice contains string
func Contains(slice []string, lookFor string) bool {
	for _, item := range slice {
		if item == lookFor {
			return true
		}
	}
	return false
}
