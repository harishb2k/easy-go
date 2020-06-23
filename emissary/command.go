package emissary

import (
    "github.com/harishb2k/easy-go/basic"
    "github.com/harishb2k/easy-go/easy"
)

// This function returns body for HTTP request
type BodyFunc func() ([]byte)

// This function returns a user response object from Http response
type ResultFunc func([]byte) (interface{}, easy.Error)

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
    Error         easy.Error
    OriginalError easy.Error
}

type Command interface {
    // A method to setup a command when it is initialized by emissary framework
    Setup(logger basic.Logger) (err error)

    // Execute a request
    Execute(request *Request) (response *Response, err easy.Error)
}

func (r *Response) SerErrorIfNotNil(err *easy.ErrorObj) {
    if err != nil && err.Err != nil {
        r.Error = err
    }
}

func (r *Response) HasError() bool {
    if r.Error != nil {
        if v, ok := r.Error.(*easy.ErrorObj); ok {
            return v.Err != nil
        }
    }
    return r.Error != nil
}

func (r *Response) DoesNotHvaeResponseBody() bool {
    return r.ResponseBody == nil || len(r.ResponseBody) <= 0
}
