package tools

import (
    "context"
    "errors"
    "fmt"
    . "github.com/harishb2k/gox-base"
    . "github.com/harishb2k/easy-go/test_http"
    "io/ioutil"
    "net/http"
    "time"
)

var DummyServer *http.Server

func readTestServerObjectFromRequest(req *http.Request) (out *TestServerObject, err error) {
    if in, err := ioutil.ReadAll(req.Body); err != nil {
        return nil, errors.New("error")
    } else {
        obj := TestServerObject{}
        if err := Objectify(in, &obj); err != nil {
            return nil, err
        } else {
            return &obj, nil
        }
    }
}

func simpleGetMethod(w http.ResponseWriter, req *http.Request) {

    if obj, err := readTestServerObjectFromRequest(req); err != nil {
        w.WriteHeader(500)
    } else {
        out := TestServerObject{}
        out.ResponseStringValue = obj.StringValue + " Response"
        out.ResponseIntValue = obj.IntValue + 1
        out.ResponseBoolValue = !obj.BoolValue

        if obj.Delay > 0 {
            time.Sleep(time.Duration(obj.Delay) * time.Millisecond)
        }

        if payloadAsString, err := StringifyWithError(out); err != nil {
            w.WriteHeader(500)
        } else {
            fmt.Fprintf(w, "%s", payloadAsString)
        }
    }
}

func simpleGetMethodWithError(w http.ResponseWriter, req *http.Request) {
    if obj, err := readTestServerObjectFromRequest(req); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    } else {
        http.Error(w, http.StatusText(http.StatusInternalServerError), obj.ErrorCodeToReturn)
    }
}

func simplePostMethod(w http.ResponseWriter, req *http.Request) {
    if obj, err := readTestServerObjectFromRequest(req); err == nil {
        out := TestServerObject{}
        out.ResponseStringValue = obj.StringValue + " Response"
        out.ResponseIntValue = obj.IntValue + 1
        out.ResponseBoolValue = !obj.BoolValue

        if obj.Delay > 0 {
            time.Sleep(time.Duration(obj.Delay) * time.Millisecond)
        }

        out.Headers = map[string]string{}
        for key, value := range req.Header {
            out.Headers[key] = value[0]
        }

        if payloadAsString, err := StringifyWithError(out); err != nil {
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        } else {
            fmt.Fprintf(w, "%s", payloadAsString)
        }
    }
}

func simplePostMethodWithError(w http.ResponseWriter, req *http.Request) {
    if obj, err := readTestServerObjectFromRequest(req); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    } else {
        http.Error(w, http.StatusText(http.StatusInternalServerError), obj.ErrorCodeToReturn)
    }
}

func StopServer() {
    if DummyServer != nil {
        if err := DummyServer.Shutdown(context.TODO()); err != nil {
            fmt.Println("Could not shutdown dummy server")
        }
    }
}

func RunServer(port string) {
    StopServer()
    go func() {
        http.HandleFunc("/v1/simpleGetMethodWithError", simpleGetMethodWithError)
        http.HandleFunc("/v1/simpleGetMethod", simpleGetMethod)
        http.HandleFunc("/v1/simplePostMethod", simplePostMethod)
        http.HandleFunc("/v1/simplePostMethodWithError", simplePostMethodWithError)
        fmt.Println("Running server at {} port", port)
        DummyServer := &http.Server{Addr: ":" + port}
        if err := DummyServer.ListenAndServe(); err != nil {
            panic("Failed to run test server at port: " + port)
        }
        fmt.Println("Running server at {} port over", port)
    }()
}
