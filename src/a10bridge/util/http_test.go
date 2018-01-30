package util_test

import (
	bridgeTesting "a10bridge/testing"
	"a10bridge/util"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HttpUtilsTestSuite struct {
	suite.Suite
	helper     *util.TestHelper
	testServer *bridgeTesting.ServerConfig
}

func TestHttpUtils(t *testing.T) {
	tests := new(HttpUtilsTestSuite)
	tests.helper = new(util.TestHelper)
	tests.testServer = bridgeTesting.NewTestServer(t).Start()
	defer tests.testServer.Stop()
	suite.Run(t, tests)
}

func (suite HttpUtilsTestSuite) TestHTTPGet() {
	expectedName := "testEntity"
	expectedNumber := 10
	url := "{{.Url}}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	suite.testServer.Reset().
		AddRequest().
		Method("GET").
		Path("/path/"+request.Name).
		Header("Connection", "Close").
		Query("test", "1").
		Query("testing", "test").
		Response().
		Body(`{"name": "`+expectedName+`","num":`+strconv.Itoa(expectedNumber)+`}`, "application/json")

	err := util.HTTPGet(url, request, &response, headers)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedName, response.Name)
	suite.Assert().Equal(expectedNumber, response.Number)
}

func (suite HttpUtilsTestSuite) TestHTTPDelete() {
	expectedName := "testEntity"
	expectedNumber := 10
	url := "{{.Url}}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	suite.testServer.Reset().
		AddRequest().
		Method("DELETE").
		Path("/path/"+request.Name).
		Header("Connection", "Close").
		Query("test", "1").
		Query("testing", "test").
		Response().
		Body(`{"name": "`+expectedName+`","num":`+strconv.Itoa(expectedNumber)+`}`, "application/json")

	err := util.HTTPDelete(url, request, &response, headers)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedName, response.Name)
	suite.Assert().Equal(expectedNumber, response.Number)
}

func (suite HttpUtilsTestSuite) TestHTTPPost() {
	url := "{{.Url}}/path?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: "test name",
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	expectedName := "test name"
	expectedNumber := 10
	suite.testServer.Reset().
		AddRequest().
		Method("POST").
		Path("/path").
		Header("Connection", "Close").
		Query("test", "1").
		Query("testing", "test").
		Response().
		Body(`{"name": "`+expectedName+`","num":`+strconv.Itoa(expectedNumber)+`}`, "application/json")

	err := util.HTTPPost(url, "a10/v2/tpl/name.request", request, &response, headers)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedName, response.Name)
	suite.Assert().Equal(expectedNumber, response.Number)
}

func (suite HttpUtilsTestSuite) TestHTTPPut() {
	expectedName := "testEntity"
	expectedNumber := 10
	url := "{{.Url}}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	suite.testServer.Reset().
		AddRequest().
		Method("PUT").
		Path("/path/"+request.Name).
		Header("Connection", "Close").
		Query("test", "1").
		Query("testing", "test").
		Response().
		Body(`{"name": "`+expectedName+`","num":`+strconv.Itoa(expectedNumber)+`}`, "application/json")

	err := util.HTTPPut(url, "a10/v2/tpl/name.request", request, &response, headers)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedName, response.Name)
	suite.Assert().Equal(expectedNumber, response.Number)
}

func (suite HttpUtilsTestSuite) TestHTTP_incorrectUrlTemplate() {
	expectedName := "testEntity"
	url := "{{.Url}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	err := util.HTTPGet(url, request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_missingRequestTemplate() {
	expectedName := "testEntity"
	url := "{{.Url}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	err := util.HTTPPut(url, "i/dont/exist", request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_incorrectRequestTemplate() {
	requestTemplate := "/tmp/a10bridge/test.request"
	err := os.MkdirAll("/tmp/a10bridge/", 0700)
	if err != nil {
		suite.T().Errorf("Failed to create folder for the template, error: %s", err)
	}
	err = ioutil.WriteFile(requestTemplate, []byte(`{
		"name": "{{.Name}"
	  }`), 0700)
	if err != nil {
		suite.T().Errorf("Failed to write template, error: %s", err)
	}
	defer os.RemoveAll("/tmp/a10bridge/")

	expectedName := "testEntity"
	url := "{{.Url}}/path?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	err = util.HTTPPost(url, requestTemplate, request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_incorrectRequestForTemplate() {
	requestTemplate := "/tmp/a10bridge/test.request"
	err := os.MkdirAll("/tmp/a10bridge/", 0700)
	if err != nil {
		suite.T().Errorf("Failed to create folder for the template, error: %s", err)
	}
	err = ioutil.WriteFile(requestTemplate, []byte(`{
		"name": "{{.Test}}"
	  }`), 0700)
	if err != nil {
		suite.T().Errorf("Failed to write template, error: %s", err)
	}
	defer os.RemoveAll("/tmp/a10bridge/")

	expectedName := "testEntity"
	url := "{{.Url}}/path?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	err = util.HTTPPost(url, requestTemplate, request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_wrongUrl() {
	url := ":"
	expectedName := "testEntity"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	err := util.HTTPGet(url, request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_emptyUrl() {
	expectedName := "testEntity"
	url := ""
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	err := util.HTTPGet(url, request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_requestFails() {
	expectedName := "testEntity"
	url := "{{.Url}}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	suite.testServer.Reset().
		AddRequest().
		Response().
		StatusCode(500)

	err := util.HTTPGet(url, request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_wrongResponse() {
	expectedName := "testEntity"
	expectedNumber := 10
	url := "{{.Url}}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	suite.testServer.Reset().
		AddRequest().
		Response().
		Body(`{"name": ,"num":`+strconv.Itoa(expectedNumber)+`}`, "application/json")

	err := util.HTTPGet(url, request, &response, headers)
	suite.Assert().NotNil(err)
}

func (suite HttpUtilsTestSuite) TestHTTP_unreadableResponse() {
	original := suite.helper.SetIoutilReadAllFunc(func(r io.Reader) ([]byte, error) {
		return nil, errors.New("read failure")
	})
	defer suite.helper.SetIoutilReadAllFunc(original)

	expectedName := "testEntity"
	expectedNumber := 10
	url := "{{.Url}}/path/{{.Name}}?{{.QS}}"
	request := request{
		Url:  suite.testServer.GetURL(),
		QS:   "test=1&testing=test",
		Name: expectedName,
	}
	response := response{}
	headers := map[string]string{
		"Connection": "Close",
	}
	suite.testServer.Reset().
		AddRequest().
		Response().
		Body(`{"name": "`+expectedName+`","num":`+strconv.Itoa(expectedNumber)+`}`, "application/json")

	err := util.HTTPGet(url, request, &response, headers)
	suite.Assert().NotNil(err)
}

type request struct {
	Url  string
	QS   string
	Name string
}

type response struct {
	Name   string `json:"name"`
	Number int    `json:"num"`
}
