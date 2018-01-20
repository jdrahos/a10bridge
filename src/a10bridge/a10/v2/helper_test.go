package v2

import "a10bridge/a10/api"

type TestHelper struct{}

//BuildError allowing the v2_test package to build unexported implementation of api.A10Error interface residing in a10 package
func (helper TestHelper) BuildError(err error) api.A10Error {
	return buildA10Error(err)
}

func (helper TestHelper) SetErrorCode(err api.A10Error, code int) api.A10Error {
	v2err, _ := err.(a10Error)
	v2err.ErrorCode = code
	return v2err
}

func (helper TestHelper) GetSessionID(client api.Client) string {
	if a10Client, ok := client.(v2Client); ok {
		return a10Client.baseRequest.SessionID
	}
	panic("client is not v2Client")
}
