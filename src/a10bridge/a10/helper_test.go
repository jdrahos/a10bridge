package a10

import "a10bridge/a10/api"

type TestHelper struct{}

//BuildError allowing the a10_test package to build unexported implementation of api.A10Error interface residing in a10 package
func (helper TestHelper) BuildError(message string) api.A10Error {
	return buildError(message)
}
