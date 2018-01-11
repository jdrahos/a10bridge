package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
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
func ApplyTemplate(data interface{}, name string, tpl string) (string, error) {
	var url string
	tmpl, err := template.New(name).Parse(tpl)
	if err != nil {
		return url, err
	}

	var writer bytes.Buffer
	err = tmpl.Execute(&writer, data)
	url = writer.String()

	fmt.Println(url)

	return url, err
}
