package testing

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ServerConfig struct {
	t        *testing.T
	requests []*HttpRequestCheck
	server   *httptest.Server
}

func NewTestServer(t *testing.T) *ServerConfig {
	return &ServerConfig{
		t:        t,
		requests: make([]*HttpRequestCheck, 0),
	}
}

func (srv *ServerConfig) Start() *ServerConfig {
	srv.server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(srv.requests) > 0 {
			requestCheck := getNextRequest(srv)
			requestCheck.assertRequest(*r)
			buildResponse(w, requestCheck)
		} else {
			srv.t.Errorf("Missing configuration for Request %s", r)
			srv.t.Fail()
		}
	}))
	return srv
}

func buildResponse(w http.ResponseWriter, requestCheck *HttpRequestCheck) {
	response := requestCheck.Response()
	w.WriteHeader(response.GetStatusCode())
	if len(response.GetBody()) > 0 {
		io.WriteString(w, response.GetBody())
	}
	if len(response.GetHeaders()) > 0 {
		for key, val := range response.GetHeaders() {
			w.Header().Add(key, val)
		}
	}
}

func (srv ServerConfig) GetUrl() string {
	if srv.server == nil {
		srv.t.Error("Trying to get server url before the server was started")
		srv.t.Fail()
	}
	return srv.server.URL
}

func (srv *ServerConfig) AddRequest() *HttpRequestCheck {
	request := NewHttpRequestCheck(srv.t)
	srv.requests = append(srv.requests, request)
	return request
}

func getNextRequest(srv *ServerConfig) *HttpRequestCheck {
	requestCheck := srv.requests[0]
	srv.requests = srv.requests[1:]
	return requestCheck
}

func (srv *ServerConfig) Stop() {
	if srv.server != nil {
		srv.server.Close()
	}
}

func (srv *ServerConfig) Reset() *ServerConfig {
	srv.requests = make([]*HttpRequestCheck, 0)
	return srv
}
