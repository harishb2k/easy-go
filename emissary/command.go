package emissary

import (
    "fmt"
    "github.com/harishb2k/easy-go/easy"
    . "github.com/harishb2k/easy-go/errors"
)

// This function returns body for HTTP request
type BodyFunc func() ([]byte)

// This function returns a user response object from Http response
type ResultFunc func([]byte) (interface{}, Error)

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
    Setup(logger easy.Logger) (err error)

    // Execute a request
    Execute(request *Request) (response *Response, err error)
}

func (r *Response) SerErrorIfNotNil(err *ErrorObj) {
    if err != nil && err.Err != nil {
        r.Error = err
    }
}

func (r *Response) HasError() bool {
    if r.Error != nil {
        if v, ok := r.Error.(*ErrorObj); ok {
            return v.Err != nil
        }
    }
    return r.Error != nil
}

func (r *Response) DoesNotHvaeResponseBody() bool {
    return r.ResponseBody == nil || len(r.ResponseBody) <= 0
}

func (r *Response) FormattedDebugString() string {
    if r.Result != nil {
        if resultString, err := easy.StringifyWithError(r.Result); err == nil {
            return fmt.Sprintf("StatusCode=%d \nError=%v \nResponse=%s ", r.StatusCode, r.Error, resultString)
        } else {
            return fmt.Sprintf("StatusCode=%d \nError=%v \nResponse=%v", r.StatusCode, r.Error, r.Result)
        }
    } else {
        return fmt.Sprintf("StatusCode=%d \nError=%v \nResponse=%v", r.StatusCode, r.Error, r.Result)
    }
}
