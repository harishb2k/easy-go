package emissary

import (
    "github.com/harishb2k/easy-go/basic"
)

// This function returns body for HTTP request
type BodyFunc func() ([]byte)

// This function returns a user response object from Http response
type ResultFunc func([]byte) (interface{}, error)

// A HTTP request
type Request struct {
    PathParam  map[string]interface{}
    QueryParam []interface{}
    Header     map[string]interface{}
    Body       interface{}
    BodyFunc   BodyFunc
    ResultFunc ResultFunc
}

// A HTTP response
type Response struct {
    Result        interface{}
    ResponseBody  []byte
    StatusCode    int
    Status        string
    Error         error
    OriginalError error
}

type Command interface {
    // A method to setup a command when it is initialized by emissary framework
    Setup(logger basic.Logger) (err error)

    // Execute a request
    Execute(request *Request) (response *Response, err error)
}
