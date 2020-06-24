package errors

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestErrorObj(t *testing.T) {
    var err error

    err = &ErrorObj{
        Name:        "error_name",
        Description: "error_description",
        Err:         New("error_err"),
        Object:      "error_object",
    }

    // Test 1 - Make sure wrap works
    errWrapped := Wrap(err, "error_1")
    assert.True(t, Is(errWrapped, err))

    // Test 2 - make sure wrapped error is able to get back the original error
    errorObj, ok := Cause(errWrapped).(*ErrorObj)
    assert.True(t, ok)
    assert.NotNil(t, errorObj)
    assert.Equal(t, "error_name", errorObj.Name)
    assert.Equal(t, "error_description", errorObj.Description)
    assert.Equal(t, "error_err", errorObj.Err.Error())
    assert.Equal(t, "error_object", errorObj.Object)

    // Test 3 - make sure wrapped error is able to get back the original error
    errorObj, ok = Cause(err).(*ErrorObj)
    assert.True(t, ok)
    assert.NotNil(t, errorObj)
    assert.Equal(t, "error_name", errorObj.Name)
    assert.Equal(t, "error_description", errorObj.Description)
    assert.Equal(t, "error_err", errorObj.Err.Error())
    assert.Equal(t, "error_object", errorObj.Object)
}
