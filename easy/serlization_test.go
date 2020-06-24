package easy

import (
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "testing"
)

type testType struct {
    A string
    B int
    C bool
}

func str(in interface{}) (string) {
    out, _ := json.Marshal(in)
    return string(out)
}

func TestStringify(t *testing.T) {

    tests := []struct {
        TestName string
        Input    interface{}
        Output   string
        err      error
    }{
        {TestName: "1", Input: 1, Output: "1", err: nil},
        {TestName: "2", Input: true, Output: "true", err: nil},
        {TestName: "3", Input: false, Output: "false", err: nil},
        {TestName: "4", Input: false, Output: "false", err: nil},
        {TestName: "5", Input: testType{A: "1", B: 1, C: true}, Output: str(testType{A: "1", B: 1, C: true}), err: nil},
        {TestName: "6", Input: int8(123), Output: "123", err: nil},
        {TestName: "7", Input: int16(12345), Output: "12345", err: nil},
        {TestName: "8", Input: int32(1234567890), Output: "1234567890", err: nil},
        {TestName: "9", Input: int64(1234567890), Output: "1234567890", err: nil},
        {TestName: "10", Input: float32(1.34), Output: "1.34", err: nil},
        {TestName: "11", Input: float64(1.34), Output: "1.34", err: nil},
    }

    for _, tt := range tests {
        t.Run(tt.TestName, func(t *testing.T) {
            out, _ := StringifyWithError(tt.Input)
            if out != tt.Output {
                t.Errorf("got %q, want %q", out, tt.Output)
            }
        })
    }
}

func TestStringifyOne(t *testing.T) {
    out, err := StringifyWithError((int32(1)))
    assert.NoError(t, err)
    assert.Equal(t, "1", out)
}
