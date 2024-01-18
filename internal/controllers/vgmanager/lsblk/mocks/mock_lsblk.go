// Code generated by mockery v2.40.1. DO NOT EDIT.

package lsblk

import (
	lsblk "github.com/openshift/lvm-operator/internal/controllers/vgmanager/lsblk"
	mock "github.com/stretchr/testify/mock"
)

// MockLSBLK is an autogenerated mock type for the LSBLK type
type MockLSBLK struct {
	mock.Mock
}

type MockLSBLK_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLSBLK) EXPECT() *MockLSBLK_Expecter {
	return &MockLSBLK_Expecter{mock: &_m.Mock}
}

// BlockDeviceInfos provides a mock function with given fields: bs
func (_m *MockLSBLK) BlockDeviceInfos(bs []lsblk.BlockDevice) (lsblk.BlockDeviceInfos, error) {
	ret := _m.Called(bs)

	if len(ret) == 0 {
		panic("no return value specified for BlockDeviceInfos")
	}

	var r0 lsblk.BlockDeviceInfos
	var r1 error
	if rf, ok := ret.Get(0).(func([]lsblk.BlockDevice) (lsblk.BlockDeviceInfos, error)); ok {
		return rf(bs)
	}
	if rf, ok := ret.Get(0).(func([]lsblk.BlockDevice) lsblk.BlockDeviceInfos); ok {
		r0 = rf(bs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(lsblk.BlockDeviceInfos)
		}
	}

	if rf, ok := ret.Get(1).(func([]lsblk.BlockDevice) error); ok {
		r1 = rf(bs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLSBLK_BlockDeviceInfos_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlockDeviceInfos'
type MockLSBLK_BlockDeviceInfos_Call struct {
	*mock.Call
}

// BlockDeviceInfos is a helper method to define mock.On call
//   - bs []lsblk.BlockDevice
func (_e *MockLSBLK_Expecter) BlockDeviceInfos(bs interface{}) *MockLSBLK_BlockDeviceInfos_Call {
	return &MockLSBLK_BlockDeviceInfos_Call{Call: _e.mock.On("BlockDeviceInfos", bs)}
}

func (_c *MockLSBLK_BlockDeviceInfos_Call) Run(run func(bs []lsblk.BlockDevice)) *MockLSBLK_BlockDeviceInfos_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]lsblk.BlockDevice))
	})
	return _c
}

func (_c *MockLSBLK_BlockDeviceInfos_Call) Return(_a0 lsblk.BlockDeviceInfos, _a1 error) *MockLSBLK_BlockDeviceInfos_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLSBLK_BlockDeviceInfos_Call) RunAndReturn(run func([]lsblk.BlockDevice) (lsblk.BlockDeviceInfos, error)) *MockLSBLK_BlockDeviceInfos_Call {
	_c.Call.Return(run)
	return _c
}

// IsUsableLoopDev provides a mock function with given fields: b
func (_m *MockLSBLK) IsUsableLoopDev(b lsblk.BlockDevice) (bool, error) {
	ret := _m.Called(b)

	if len(ret) == 0 {
		panic("no return value specified for IsUsableLoopDev")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(lsblk.BlockDevice) (bool, error)); ok {
		return rf(b)
	}
	if rf, ok := ret.Get(0).(func(lsblk.BlockDevice) bool); ok {
		r0 = rf(b)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(lsblk.BlockDevice) error); ok {
		r1 = rf(b)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLSBLK_IsUsableLoopDev_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsUsableLoopDev'
type MockLSBLK_IsUsableLoopDev_Call struct {
	*mock.Call
}

// IsUsableLoopDev is a helper method to define mock.On call
//   - b lsblk.BlockDevice
func (_e *MockLSBLK_Expecter) IsUsableLoopDev(b interface{}) *MockLSBLK_IsUsableLoopDev_Call {
	return &MockLSBLK_IsUsableLoopDev_Call{Call: _e.mock.On("IsUsableLoopDev", b)}
}

func (_c *MockLSBLK_IsUsableLoopDev_Call) Run(run func(b lsblk.BlockDevice)) *MockLSBLK_IsUsableLoopDev_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(lsblk.BlockDevice))
	})
	return _c
}

func (_c *MockLSBLK_IsUsableLoopDev_Call) Return(_a0 bool, _a1 error) *MockLSBLK_IsUsableLoopDev_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLSBLK_IsUsableLoopDev_Call) RunAndReturn(run func(lsblk.BlockDevice) (bool, error)) *MockLSBLK_IsUsableLoopDev_Call {
	_c.Call.Return(run)
	return _c
}

// ListBlockDevices provides a mock function with given fields:
func (_m *MockLSBLK) ListBlockDevices() ([]lsblk.BlockDevice, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListBlockDevices")
	}

	var r0 []lsblk.BlockDevice
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]lsblk.BlockDevice, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []lsblk.BlockDevice); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]lsblk.BlockDevice)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockLSBLK_ListBlockDevices_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListBlockDevices'
type MockLSBLK_ListBlockDevices_Call struct {
	*mock.Call
}

// ListBlockDevices is a helper method to define mock.On call
func (_e *MockLSBLK_Expecter) ListBlockDevices() *MockLSBLK_ListBlockDevices_Call {
	return &MockLSBLK_ListBlockDevices_Call{Call: _e.mock.On("ListBlockDevices")}
}

func (_c *MockLSBLK_ListBlockDevices_Call) Run(run func()) *MockLSBLK_ListBlockDevices_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockLSBLK_ListBlockDevices_Call) Return(_a0 []lsblk.BlockDevice, _a1 error) *MockLSBLK_ListBlockDevices_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockLSBLK_ListBlockDevices_Call) RunAndReturn(run func() ([]lsblk.BlockDevice, error)) *MockLSBLK_ListBlockDevices_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockLSBLK creates a new instance of MockLSBLK. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockLSBLK(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockLSBLK {
	mock := &MockLSBLK{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
