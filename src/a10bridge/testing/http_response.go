package testing

type HttpResponseConfig struct {
	requestCheck *HttpRequestCheck
	body         string
	statuscode   int
	headers      map[string]string
}

func (resp *HttpResponseConfig) Body(respBody string, contentType string) *HttpResponseConfig {
	resp.body = respBody
	resp.headers["Content-Type"] = contentType
	return resp
}

func (resp HttpResponseConfig) GetBody() string {
	return resp.body
}

func (resp *HttpResponseConfig) StatusCode(statusCode int) *HttpResponseConfig {
	resp.statuscode = statusCode
	return resp
}

func (resp HttpResponseConfig) GetStatusCode() int {
	return resp.statuscode
}

func (resp *HttpResponseConfig) Header(key, val string) *HttpResponseConfig {
	resp.headers[key] = val
	return resp
}

func (resp HttpResponseConfig) GetHeaders() map[string]string {
	return resp.headers
}

func (resp *HttpResponseConfig) Request() *HttpRequestCheck {
	return resp.requestCheck
}
