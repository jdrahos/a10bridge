package testing

type HTTPResponseConfig struct {
	requestCheck *HTTPRequestCheck
	body         string
	statuscode   int
	headers      map[string]string
}

func (resp *HTTPResponseConfig) Body(respBody string, contentType string) *HTTPResponseConfig {
	resp.body = respBody
	resp.headers["Content-Type"] = contentType
	return resp
}

func (resp HTTPResponseConfig) GetBody() string {
	return resp.body
}

func (resp *HTTPResponseConfig) StatusCode(statusCode int) *HTTPResponseConfig {
	resp.statuscode = statusCode
	return resp
}

func (resp HTTPResponseConfig) GetStatusCode() int {
	return resp.statuscode
}

func (resp *HTTPResponseConfig) Header(key, val string) *HTTPResponseConfig {
	resp.headers[key] = val
	return resp
}

func (resp HTTPResponseConfig) GetHeaders() map[string]string {
	return resp.headers
}

func (resp *HTTPResponseConfig) Request() *HTTPRequestCheck {
	return resp.requestCheck
}
