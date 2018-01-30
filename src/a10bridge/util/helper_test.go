package util

import "io"

type IoutilReadAllFunc func(r io.Reader) ([]byte, error)
type TestHelper struct{}

func (helper TestHelper) SetIoutilReadAllFunc(ioutilReadAllFunc IoutilReadAllFunc) IoutilReadAllFunc {
	old := ioutilReadAll
	ioutilReadAll = ioutilReadAllFunc
	return old
}
