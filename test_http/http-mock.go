package test_http

import (
    "encoding/json"
    "github.com/jarcoal/httpmock"
    "net/http"
    "time"
)

func SetupGetHttpResponse(url string, data string, obj interface{}) {
    if err := json.Unmarshal([]byte(data), obj);
        err != nil {
        panic(err)
    }
    httpmock.RegisterResponder("GET", url,
        func(req *http.Request) (*http.Response, error) {
            resp, err := httpmock.NewJsonResponse(200, obj)
            if err != nil {
                return httpmock.NewStringResponse(500, ""), nil
            }
            return resp, nil
        },
    )
}

func SetupGetHttpResponseDelay(url string, data string, obj interface{}, delay int) {
    if err := json.Unmarshal([]byte(data), obj); err != nil {
        panic(err)
    }
    httpmock.RegisterResponder("GET", url,
        func(req *http.Request) (*http.Response, error) {
            time.Sleep(time.Duration(delay) * time.Millisecond)
            resp, err := httpmock.NewJsonResponse(200, obj)
            if err != nil {
                return httpmock.NewStringResponse(500, ""), nil
            }
            return resp, nil
        },
    )
}
