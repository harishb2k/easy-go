package easy

import (
    "encoding/json"
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
