package emissary

import (
    "context"
    "github.com/google/uuid"
    "github.com/harishb2k/easy-go/easy"
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

    // Make correct http request
    var httpRequest *http.Request;
    switch c.Api.Method {
    case "GET":
        if httpRequest, err = http.NewRequest("GET", url, nil); err != nil {
            return nil, &easy.ErrorObj{
                Err:         err,
                Name:        "http_call_failed",
                Description: "Failed to create http request for " + c.commandName(),
                Object:      &Response{StatusCode: 500, Status: "Unknown"},
            }
        }
        break
    }

    // Make http call
    var httpResponse *http.Response;
    if httpResponse, err = http.DefaultClient.Do(httpRequest.WithContext(ctx)); err != nil {
        return nil, &easy.ErrorObj{
            Err:         err,
            Name:        "http_call_failed",
            Description: "Http request failed with error " + c.commandName() + " " + err.Error(),
            Object:      &Response{StatusCode: 500, Status: "Unknown"},
        }
    }
    defer httpResponse.Body.Close()

    // Fill response object from http response
    response = &Response{}
    c.populateResponse(requestId, request, response, httpResponse)
    c.Debug(requestId, "HttpCommand: response -", "statusCode=", response.StatusCode, "result=", response.Result, "error=", response.Error)

    // See if we accept a error
    c.handleAcceptedCodes(requestId, response)

    return response, response.Error
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

// Populate response
func (c *HttpCommand) populateResponse(reqId string, request *Request, response *Response, httpResponse *http.Response) {

    response.StatusCode = httpResponse.StatusCode
    response.Status = httpResponse.Status

    // Read body from Http response
    if httpResponse.Body != nil {

        var err error
        response.ResponseBody, err = ioutil.ReadAll(httpResponse.Body);

        // Setup correct error
        response.SerErrorIfNotNil(&easy.ErrorObj{
            Err:         err,
            Name:        "http_call_failed",
            Description: "Failed dto read http response body",
            Object:      &Response{StatusCode: 500, Status: "Unknown"},
        })

        if response.HasError() || response.DoesNotHvaeResponseBody() {
            return
        }
    }

    // Convert http response to requested Pojo
    if request.ResultFunc != nil {
        response.Result, response.Error = request.ResultFunc(response.ResponseBody)
    }
}

// Populate response
func (c *HttpCommand) handleAcceptedCodes(reqId string, response *Response) {
    response.OriginalError = response.Error

    var accepted = false
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
}
