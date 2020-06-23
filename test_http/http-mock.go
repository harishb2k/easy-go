package test_http

import (
    "encoding/json"
    "github.com/jarcoal/httpmock"
    "net/http"
    "time"
)

type HttpMockSpec struct {
    Method      string
    Url         string
    Data        string
    ResponseObj interface{}
    Delay       int
    StatusCode  int
}

func (spec *HttpMockSpec) setupDefault() {
    if spec.Method == "" {
        spec.Method = "GET"
    }
    if spec.StatusCode == 0 {
        spec.StatusCode = 200
    }
}

func SetupMockHttpResponse(spec HttpMockSpec) {
    spec.setupDefault()

    if err := json.Unmarshal([]byte(spec.Data), spec.ResponseObj);
        err != nil {
        panic(err)
    }

    httpmock.RegisterResponder(spec.Method, spec.Url,
        func(req *http.Request) (*http.Response, error) {

            // Setup Delay if requested
            if spec.Delay > 0 {
                time.Sleep(time.Duration(spec.Delay) * time.Millisecond)
            }

            // Send final response with status code
            resp, err := httpmock.NewJsonResponse(spec.StatusCode, spec.ResponseObj)
            if err != nil {
                return httpmock.NewStringResponse(500, ""), nil
            }
            return resp, nil
        },
    )

}
