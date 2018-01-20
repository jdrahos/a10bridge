package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var httpClient = buildHTTPClient()
var templateRoot = buildTemplateRoot()

//HTTPGet performs http GET call
func HTTPGet(url string, request interface{}, response interface{}, headers map[string]string) error {
	return httpCall("GET", url, "", request, response, headers)
}

//HTTPDelete performs http DELETE call
func HTTPDelete(url string, request interface{}, response interface{}, headers map[string]string) error {
	return httpCall("DELETE", url, "", request, response, headers)
}

//HTTPPost performs http POST call
func HTTPPost(url string, tplPath string, request interface{}, response interface{}, headers map[string]string) error {
	return httpCall("POST", url, tplPath, request, response, headers)
}

//HTTPPut performs http PUT call
func HTTPPut(url string, tplPath string, request interface{}, response interface{}, headers map[string]string) error {
	return httpCall("PUT", url, tplPath, request, response, headers)
}

func buildHTTPClient() *http.Client {
	timeout := time.Second * 30
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

func buildTemplateRoot() string {
	goroot := os.Getenv("GOPATH")
	if len(goroot) > 0 && !strings.HasSuffix(goroot, "/") {
		goroot = goroot + "/"
	}
	return goroot + "src/a10bridge/"
}

func httpCall(method string, urlTpl string, tplPath string, request interface{}, response interface{}, headers map[string]string) error {
	var requestReader io.Reader

	url, err := ApplyTemplate(request, urlTpl)
	if err != nil {
		return err
	}

	if len(tplPath) > 0 {
		if !filepath.IsAbs(tplPath) {
			tplPath = templateRoot + tplPath
		}
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

		requestReader = bytes.NewBuffer(writer.Bytes())
	} else {
		//		requestReader = nil
	}

	httpRequest, err := http.NewRequest(method, url, requestReader)
	if err != nil {
		return err
	}
	addHeaders(httpRequest, headers)

	if request != nil {
		httpRequest.Header.Add("Content-Type", "application/json")
	}

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		return err
	}

	return processResponse(httpResponse, &response)
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

func addHeaders(request *http.Request, headers map[string]string) {
	if headers != nil && len(headers) > 0 {
		for header, headerVal := range headers {
			request.Header.Add(header, headerVal)
		}
	}
}
