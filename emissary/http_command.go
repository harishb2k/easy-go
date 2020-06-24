package emissary

import (
    "bytes"
    "context"
    "errors"
    "github.com/google/uuid"
    "github.com/harishb2k/easy-go/easy"
    "io"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type HttpCommand struct {
    Service *Service
    Api     *Api
    easy.Logger
}

func NewHttpCommand(service *Service, api *Api, logger easy.Logger) (HttpCommand) {
    command := HttpCommand{
        Service: service,
        Api:     api,
    }
    command.Setup(logger)
    return command
}

// Setup this command
func (c *HttpCommand) Setup(logger easy.Logger) (err error) {
    if logger != nil {
        c.Logger = logger
    } else {
        c.Logger = &easy.DefaultLogger{};
    }
    return nil
}

// Execute a request
func (c *HttpCommand) Execute(request *Request) (response *Response, e easy.Error) {
    var err error
    requestId := uuid.New().String()

    // Setup http timeout to kill request if it takes longer
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Api.RequestTimeout)*time.Millisecond)
    defer cancel()

    // Build a URL with path param
    var url = c.getUrl(requestId, request)

    var payload []byte
    var body io.Reader

    if request.Body != nil {
        if payload, err = easy.BytesWithError(request.Body); err != nil {
            return nil, &easy.ErrorObj{
                Err:         err,
                Name:        "http_call_failed",
                Description: "Failed to convert object to byte body for " + c.commandName(),
                Object:      &Response{StatusCode: 500, Status: "Unknown"},
            }
        } else if payload != nil {
            body = bytes.NewReader(payload)
        }
    }

    // Make correct http request
    var httpRequest *http.Request;
    switch c.Api.Method {
    case "GET":
        if httpRequest, err = http.NewRequest("GET", url, body); err != nil {
            return nil, &easy.ErrorObj{
                Err:         err,
                Name:        "http_call_failed",
                Description: "Failed to create http request for " + c.commandName(),
                Object:      &Response{StatusCode: 500, Status: "Unknown"},
            }
        }
        break

    case "POST":
        if httpRequest, err = http.NewRequest("POST", url, body); err != nil {
            return nil, &easy.ErrorObj{
                Err:         err,
                Name:        "http_call_failed",
                Description: "Failed to create http request for " + c.commandName(),
                Object:      &Response{StatusCode: 500, Status: "Unknown"},
            }
        }
        break

    case "PUT":
        if httpRequest, err = http.NewRequest("PUT", url, body); err != nil {
            return nil, &easy.ErrorObj{
                Err:         err,
                Name:        "http_call_failed",
                Description: "Failed to create http request for " + c.commandName(),
                Object:      &Response{StatusCode: 500, Status: "Unknown"},
            }
        }
        break

    case "DELETE":
        if httpRequest, err = http.NewRequest("DELETE", url, body); err != nil {
            return nil, &easy.ErrorObj{
                Err:         err,
                Name:        "http_call_failed",
                Description: "Failed to create http request for " + c.commandName(),
                Object:      &Response{StatusCode: 500, Status: "Unknown"},
            }
        }
        break

    default:
        return nil, &easy.ErrorObj{
            Err:         errors.New("method not supported"),
            Name:        "http_call_failed",
            Description: "Http request does not support method type " + c.commandName(),
            Object:      &Response{StatusCode: 500, Status: "Unknown"},
        }
    }

    // Setup Headers
    c.populateHeaders(requestId, request, httpRequest)

    // Make http call
    var httpResponse *http.Response;
    if httpResponse, err = http.DefaultClient.Do(httpRequest.WithContext(ctx)); err != nil {
        return nil, &easy.ErrorObj{
            Err:         err,
            Name:        ErrorCodeHttpServerTimeout,
            Description: "Http request failed with error " + c.commandName() + " " + err.Error(),
            Object:      &Response{StatusCode: 500, Status: "Unknown"},
        }
    }
    defer httpResponse.Body.Close()

    // Fill response object from http response
    response = &Response{}
    c.populateResponse(requestId, request, response, httpResponse)
    c.Debug(requestId, "HttpCommand: response -", "statusCode=", response.StatusCode, "result=", response.Result, "error=", response.Error)

    if response.Error == nil {
        return response, response.Error
    } else {
        return nil, response.Error
    }
}

// Build command name for hystrix or debug name
func (c *HttpCommand) commandName() (string) {
    return c.Service.Name + "_" + c.Api.Name
}

// Make a URL from command
func (c *HttpCommand) getUrl(reqId string, request *Request) (string) {

    var url = c.Service.Type + "://" + c.Service.Host + ":" + strconv.Itoa(c.Service.Port) + c.Api.Path
    c.Debug(reqId, "HttpCommand: url=", url)

    if request.PathParam != nil {
        for key, value := range request.PathParam {
            url = strings.Replace(url, "$${"+key+"}", easy.Stringify(value), 1)
        }
        c.Debug(reqId, "HttpCommand: after path param replacement url=", url)
    }

    return url
}

func (c *HttpCommand) populateHeaders(reqId string, request *Request, httpRequest *http.Request) {

    // Setup headers
    if request.Header != nil {
        for key, value := range request.Header {
            httpRequest.Header.Set(key, easy.Stringify(value))
        }
    }

    // Setup default content-type (if missing)
    if httpRequest.Header.Get("Content-Type") == "" && httpRequest.Header.Get("content-type") == "" {
        httpRequest.Header.Add("Content-Type", "application/json");
    }
}

// Populate response
func (c *HttpCommand) populateResponse(reqId string, request *Request, response *Response, httpResponse *http.Response) {
    var err error
    var acceptedErr error

    response.StatusCode = httpResponse.StatusCode
    response.Status = httpResponse.Status

    if response.StatusCode >= 200 && response.StatusCode < 300 {
        // No-OP
    } else if c.Api.AcceptableResponseCodes != nil && len(c.Api.AcceptableResponseCodes) > 0 {
        var accepted = false
        for _, c := range c.Api.AcceptableResponseCodes {
            if c == response.StatusCode {
                accepted = true
                break
            }
        }
        if !accepted {
            acceptedErr = errors.New("request failed with statusCode=" + strconv.Itoa(response.StatusCode) + " status=" + response.Status)
        }
    } else {
        acceptedErr = errors.New("request failed with statusCode=" + strconv.Itoa(response.StatusCode) + " status=" + response.Status)
    }

    // Read body from Http response
    if httpResponse.Body != nil {
        if body, err := ioutil.ReadAll(httpResponse.Body); err == nil {
            response.ResponseBody = body
        }
    }

    if err == nil && acceptedErr == nil {
        if request.ResultFunc != nil {
            response.Result, response.Error = request.ResultFunc(response.ResponseBody)
        }
    } else if acceptedErr != nil {
        response.SerErrorIfNotNil(&easy.ErrorObj{
            Err:         acceptedErr,
            Name:        ErrorCodeHttpServerApiError,
            Description: "Api call failed with error: code=" + easy.Stringify(response.StatusCode),
            Object:      &Response{StatusCode: response.StatusCode, Status: easy.Stringify(response.StatusCode)},
        })
    } else if err != nil {
        response.SerErrorIfNotNil(&easy.ErrorObj{
            Err:         err,
            Name:        ErrorCodeHttpFailedToReadBody,
            Description: "Failed to read body from http response",
            Object:      &Response{StatusCode: response.StatusCode, Status: easy.Stringify(response.StatusCode)},
        })
    }
}

// Populate response
func (c *HttpCommand) handleAcceptedCodes(reqId string, response *Response) (accepted bool) {

    if c.Api.AcceptableResponseCodes != nil && len(c.Api.AcceptableResponseCodes) > 0 {
        for c := range c.Api.AcceptableResponseCodes {
            if c == response.StatusCode {
                accepted = true
            }
        }
    } else {
        if response.StatusCode >= 200 && response.StatusCode < 300 {
            accepted = true
        }
    }

    if accepted {
        response.Error = nil
    }

    return
}
