package errors

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

type testingInterface interface {
    SomeMethod() (err error)
    SomeMethodWithErrorObj() (err error)
    SomeMethodWithErrorObjFull() (err error)
}

type testingInterfaceImpl struct {
    Data string
}

func (t *testingInterfaceImpl) SomeMethod() (err error) {
    return New("some_error")
}

func (t *testingInterfaceImpl) SomeMethodWithErrorWrap() (err error) {
    return Wrap(New("some_error"), "SomeMethodWithErrorObj")
}

func (t *testingInterfaceImpl) SomeMethodWithErrorObj() (err error) {
    return Wrap(
        &ErrorObj{
            Name:        "error_name",
            Description: "error_description",
            Err:         New("error_err"),
            Object:      "error_object",
        },
        "SomeMethodWithErrorObj",
    )
}

func TestUsage(t *testing.T) {
    var ti = &testingInterfaceImpl{Data: "Data"}

    err := ti.SomeMethod()
    assert.NotNil(t, err)

    errorObj, ok := Cause(err).(*ErrorObj)
    assert.False(t, ok)
    assert.Nil(t, errorObj)
}

func TestUsageWrap(t *testing.T) {
    var ti = &testingInterfaceImpl{Data: "Data"}

    err := ti.SomeMethodWithErrorWrap()
    assert.NotNil(t, err)

    errorObj, ok := Cause(err).(*ErrorObj)
    assert.False(t, ok)
    assert.Nil(t, errorObj)
}

func TestUsageErrorObj(t *testing.T) {
    var ti = &testingInterfaceImpl{Data: "Data"}

    err := ti.SomeMethodWithErrorObj()
    assert.NotNil(t, err)

    errorObj, ok := Cause(err).(*ErrorObj)
    assert.True(t, ok)
    assert.NotNil(t, errorObj)
    assert.Equal(t, "error_name", errorObj.Name)
    assert.Equal(t, "error_description", errorObj.Description)
    assert.Equal(t, "error_err", errorObj.Err.Error())
    assert.Equal(t, "error_object", errorObj.Object)

    if e, ok := AsErrorObj(err); ok {
        assert.Equal(t, "error_name", e.Name)
        assert.Equal(t, "error_description", e.Description)
        assert.Equal(t, "error_err", e.Err.Error())
        assert.Equal(t, "error_object", e.Object)
    } else {
        assert.Fail(t, "Expected to get Error Object from err")
    }
}
