package testing

import (
	"io/ioutil"
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
	expectedQuery   map[string]string
	expectedHeaders map[string]string
}

func NewHttpRequestCheck(t *testing.T) HttpRequestCheck {
	return HttpRequestCheck{
		t:               t,
		assert:          tassert.New(t),
		expectedMethod:  "",
		expectedBody:    "",
		expectedQuery:   make(map[string]string),
		expectedHeaders: make(map[string]string),
	}
}

func (check HttpRequestCheck) Path(path string) HttpRequestCheck {
	check.expectedPath = path
	return check
}

func (check HttpRequestCheck) Query(key, val string) HttpRequestCheck {
	check.expectedQuery[key] = val
	return check
}

func (check HttpRequestCheck) Header(key, val string) HttpRequestCheck {
	check.expectedHeaders[key] = val
	return check
}

func (check HttpRequestCheck) Body(body string) HttpRequestCheck {
	check.expectedBody = strings.Replace(body, "\t", "    ", -1)
	return check
}

func (check HttpRequestCheck) Method(method string) HttpRequestCheck {
	check.expectedMethod = method
	return check
}

func (check HttpRequestCheck) Assert(req http.Request) {
	check.request = req
	if len(check.expectedPath) != 0 {
		check.assertPath()
	}
	if len(check.expectedQuery) != 0 {
		check.assertQuery()
	}
	if len(check.expectedHeaders) != 0 {
		check.assertHeaders()
	}
	if len(check.expectedMethod) != 0 {
		check.assertMethod()
	}
	if len(check.expectedBody) != 0 {
		check.assertBody()
	}
}

func (check HttpRequestCheck) assertPath() {
	if check.request.URL.Path != check.expectedPath {
		check.t.Errorf("Unexpected call to '%s'. Expected call to '%s'", check.request.URL.Path, check.expectedPath)
	}
}

func (check HttpRequestCheck) assertMethod() {
	if check.request.Method != check.expectedMethod {
		check.t.Errorf("Unexpected method '%s'. Expected '%s' method", check.request.URL.Path, check.expectedPath)
	}
}

func (check HttpRequestCheck) assertQuery() {
	for qvar, qval := range check.expectedQuery {
		actual := check.request.URL.Query().Get(qvar)
		check.assert.Equalf(qval, actual, "'%s' query paramter should have value '%s' but was '%s'", qvar, qval, actual)
	}
}

func (check HttpRequestCheck) assertHeaders() {
	for qvar, qval := range check.expectedHeaders {
		actual := check.request.Header.Get(qvar)
		check.assert.Equalf(qval, actual, "'%s' header should have value '%s' but was '%s'", qvar, qval, actual)
	}
}

func (check HttpRequestCheck) assertBody() {
	binary, err := ioutil.ReadAll(check.request.Body)
	if err != nil {
		check.t.Errorf("Failed to read request body to check %s", err)
		return
	}
	httpBody := string(binary)
	check.assert.EqualValues(check.expectedBody, httpBody, "Unexpected body '%s'. Expected '%s'", httpBody, check.expectedBody)
}
