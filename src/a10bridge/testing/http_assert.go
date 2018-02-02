package testing

import (
	"io/ioutil"
	"net/http"
)

type StringInspector func(toBeInspected string) (bool, string)

func (check HTTPRequestCheck) assertRequest(req http.Request) {
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
	if check.bodyInspector != nil {
		check.inspectBody()
	}
}

func (check HTTPRequestCheck) assertPath() {
	if check.request.URL.Path != check.expectedPath {
		check.t.Errorf("Unexpected call to '%s'. Expected call to '%s'", check.request.URL.Path, check.expectedPath)
	}
}

func (check HTTPRequestCheck) assertMethod() {
	if check.request.Method != check.expectedMethod {
		check.t.Errorf("Unexpected method '%s'. Expected '%s' method", check.request.Method, check.expectedMethod)
	}
}

func (check HTTPRequestCheck) assertQuery() {
	for qvar, qval := range check.expectedQuery {
		actual := check.request.URL.Query().Get(qvar)
		check.assert.Equalf(qval, actual, "'%s' query paramter should have value '%s' but was '%s'", qvar, qval, actual)
	}
}

func (check HTTPRequestCheck) assertHeaders() {
	for qvar, qval := range check.expectedHeaders {
		actual := check.request.Header.Get(qvar)
		check.assert.Equalf(qval, actual, "'%s' header should have value '%s' but was '%s'", qvar, qval, actual)
	}
}

func (check HTTPRequestCheck) assertBody() {
	binary, err := ioutil.ReadAll(check.request.Body)
	if err != nil {
		check.t.Errorf("Failed to read request body to check %s", err)
		return
	}
	httpBody := normalizeJson(string(binary))
	check.assert.EqualValues(check.expectedBody, httpBody, "Unexpected body '%s'. Expected '%s'", httpBody, check.expectedBody)
}

func (check HTTPRequestCheck) inspectBody() {
	binary, err := ioutil.ReadAll(check.request.Body)
	if err != nil {
		check.t.Errorf("Failed to read request body to check %s", err)
		return
	}
	httpBody := string(binary)
	result, msg := check.bodyInspector(httpBody)
	if !result {
		check.t.Errorf("Body inspection has failed: %s", msg)
	}

}

func (server ServerConfig) AssertNoPendingRequests() {
	if len(server.requests) > 0 {
		server.t.Errorf("There are still %d pending requests", len(server.requests))
	}
}
