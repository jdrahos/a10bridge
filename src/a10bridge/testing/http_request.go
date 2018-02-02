package testing

import (
	"net/http"
	"regexp"
	"testing"

	tassert "github.com/stretchr/testify/require"
)

type HTTPRequestCheck struct {
	request http.Request
	t       *testing.T
	assert  *tassert.Assertions

	expectedPath    string
	expectedMethod  string
	expectedBody    string
	bodyInspector   StringInspector
	expectedQuery   map[string]string
	expectedHeaders map[string]string

	response *HTTPResponseConfig
}

func NewHTTPRequestCheck(t *testing.T) *HTTPRequestCheck {
	return &HTTPRequestCheck{
		t:               t,
		assert:          tassert.New(t),
		expectedMethod:  "",
		expectedBody:    "",
		expectedQuery:   make(map[string]string),
		expectedHeaders: make(map[string]string),
		response: &HTTPResponseConfig{
			body:       "",
			statuscode: 200,
			headers:    make(map[string]string),
		},
	}
}

func (check *HTTPRequestCheck) Path(path string) *HTTPRequestCheck {
	check.expectedPath = path
	return check
}

func (check *HTTPRequestCheck) Query(key, val string) *HTTPRequestCheck {
	check.expectedQuery[key] = val
	return check
}

func (check *HTTPRequestCheck) Header(key, val string) *HTTPRequestCheck {
	check.expectedHeaders[key] = val
	return check
}

func (check *HTTPRequestCheck) Body(body string) *HTTPRequestCheck {
	check.expectedBody = normalizeJson(body)
	return check
}

func (check *HTTPRequestCheck) BodyInspector(bodyInspector StringInspector) *HTTPRequestCheck {
	check.bodyInspector = bodyInspector
	return check
}

func (check *HTTPRequestCheck) Method(method string) *HTTPRequestCheck {
	check.expectedMethod = method
	return check
}

func (check *HTTPRequestCheck) Response() *HTTPResponseConfig {
	check.response.requestCheck = check
	return check.response
}

func normalizeJson(json string) string {
	rexp := regexp.MustCompile(`(\s{2,}|\n)`)
	return rexp.ReplaceAllString(json, " ")
}
