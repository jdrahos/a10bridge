package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

func HttpGet(url string, response interface{}) error {
	httpResponse, err := httpClient().Get(url)
	if err != nil {
		return err
	}
	return processResponse(httpResponse, &response)
}

func HttpPost(url string, tplPath string, request interface{}, response interface{}) error {
	tmpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}

	var writer bytes.Buffer
	err = tmpl.Execute(&writer, request)
	if err != nil {
		return err
	}

	fmt.Println(string(writer.Bytes()))

	requestReader := bytes.NewBuffer(writer.Bytes())

	httpResponse, err := httpClient().Post(url, "application/json", requestReader)
	if err != nil {
		return err
	}

	return processResponse(httpResponse, &response)
}

func httpClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: transport}
}

func processResponse(httpResponse *http.Response, response interface{}) error {
	defer httpResponse.Body.Close()
	binary, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(binary))

	err = json.Unmarshal(binary, &response)
	if err != nil {
		return err
	}

	return err
}
