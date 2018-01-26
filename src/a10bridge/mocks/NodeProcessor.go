// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import model "a10bridge/model"

// NodeProcessor is an autogenerated mock type for the NodeProcessor type
type NodeProcessor struct {
	mock.Mock
}

// ProcessNode provides a mock function with given fields: node
func (_m *NodeProcessor) ProcessNode(node *model.Node) error {
	ret := _m.Called(node)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Node) error); ok {
		r0 = rf(node)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
