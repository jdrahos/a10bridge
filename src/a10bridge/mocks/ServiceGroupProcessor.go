// Code generated by mockery v1.0.0
package mocks

import mock "github.com/stretchr/testify/mock"
import model "a10bridge/model"

// ServiceGroupProcessor is an autogenerated mock type for the ServiceGroupProcessor type
type ServiceGroupProcessor struct {
	mock.Mock
}

// ProcessServiceGroup provides a mock function with given fields: serviceGroup, failedNodeNames
func (_m *ServiceGroupProcessor) ProcessServiceGroup(serviceGroup *model.ServiceGroup, failedNodeNames []string) error {
	ret := _m.Called(serviceGroup, failedNodeNames)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.ServiceGroup, []string) error); ok {
		r0 = rf(serviceGroup, failedNodeNames)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
