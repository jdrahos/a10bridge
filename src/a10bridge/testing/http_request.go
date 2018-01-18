package testing

import (
	"net/http"
	"strings"
	"testing"

	tassert "github.com/stretchr/testify/assert"
)

type HttpRequestCheck struct {
	request http.Request
	t       *testing.T
	assert  *tassert.Assertions

	expectedPath    string
	expectedMethod  string
	expectedBody    string
	bodyInspector   StringInspector
	expectedQuery   map[string]string
	expectedHeaders map[string]string

	response *HttpResponseConfig
}

func NewHttpRequestCheck(t *testing.T) *HttpRequestCheck {
	return &HttpRequestCheck{
		t:               t,
		assert:          tassert.New(t),
		expectedMethod:  "",
		expectedBody:    "",
		expectedQuery:   make(map[string]string),
		expectedHeaders: make(map[string]string),
		response: &HttpResponseConfig{
			body:       "",
			statuscode: 200,
			headers:    make(map[string]string),
		},
	}
}

func (check *HttpRequestCheck) Path(path string) *HttpRequestCheck {
	check.expectedPath = path
	return check
}

func (check *HttpRequestCheck) Query(key, val string) *HttpRequestCheck {
	check.expectedQuery[key] = val
	return check
}

func (check *HttpRequestCheck) Header(key, val string) *HttpRequestCheck {
	check.expectedHeaders[key] = val
	return check
}

func (check *HttpRequestCheck) Body(body string) *HttpRequestCheck {
	check.expectedBody = strings.Replace(body, "\t", "    ", -1)
	return check
}

func (check *HttpRequestCheck) BodyInspector(bodyInspector StringInspector) *HttpRequestCheck {
	check.bodyInspector = bodyInspector
	return check
}

func (check *HttpRequestCheck) Method(method string) *HttpRequestCheck {
	check.expectedMethod = method
	return check
}

func (check *HttpRequestCheck) Response() *HttpResponseConfig {
	check.response.requestCheck = check
	return check.response
}
