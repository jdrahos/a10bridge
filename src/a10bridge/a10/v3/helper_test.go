package v3

import (
	"a10bridge/a10/api"
	"strings"
)

type TestHelper struct{}

//BuildError allowing the v3_test package to build unexported implementation of api.A10Error interface residing in a10 package
func (helper TestHelper) BuildError(err error) api.A10Error {
	return buildA10Error(err)
}

func (helper TestHelper) SetErrorCode(err api.A10Error, code int) api.A10Error {
	v2err, _ := err.(a10Error)
	v2err.ErrorCode = code
	return v2err
}

func (helper TestHelper) GetSessionID(client api.Client) string {
	if a10Client, ok := client.(v3Client); ok {
		authHeader, _ := a10Client.commonHeaders["Authorization"]
		return strings.TrimPrefix(authHeader, "A10 ")
	}
	panic("client is not v3Client")
}
