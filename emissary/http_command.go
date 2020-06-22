package emissary

import (
    "bytes"
    "context"
    "errors"
    "github.com/google/uuid"
    _ "github.com/google/uuid"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
)

type HttpCommand struct {
    Service Service
    Api     Api
    Logger  ILogger
}

// Build command name for hystrix
func (c *HttpCommand) commandName() (string) {
    return c.Service.Name + "_" + c.Api.Name
}

func (c *HttpCommand) getUrl() (string) {
    return c.Service.Type + "://" + c.Service.Host + ":" + strconv.Itoa(c.Service.Port) + c.Api.Path;
}

// Setup this command
func (c *HttpCommand) Setup(logger ILogger) (err error) {
    if logger != nil {
        c.Logger = logger
    } else {
        c.Logger = &DefaultLogger{};
    }
    return nil
}

// Make Http Request
func (c *HttpCommand) Execute(request CommandRequest) (result interface{}, err error) {
    requestId := uuid.New().String()

    // Setup http timeout to kill request if it takes longer
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Api.RequestTimeout)*time.Millisecond)
    defer cancel()

    var url = c.getUrl()
    c.Logger.Log(requestId, c.commandName(), "URL", url)

    // Update URL with all input params
    if request.PathParamFunc != nil {
        var pathParam = request.PathParamFunc()
        if pathParam != nil {
            for key, value := range pathParam {
                url = strings.Replace(url, "$${"+key+"}", Stringify(value), 1)
            }
        }
    }
    c.Logger.Log(requestId, c.commandName(), "Modified URL", url)

    var httpRequest *http.Request;
    switch c.Api.Method {
    case "GET":
        if httpRequest, err = http.NewRequest("GET", url, nil); err != nil {
            return nil, errors.New("Failed to create http request for " + c.commandName())
        }
        break

    case "POST":
        var body *bytes.Buffer;
        if request.BodyFunc != nil {
            bodyData := request.BodyFunc()
            if bodyData != nil {
                body = bytes.NewBuffer(bodyData)
            }
        }
        if httpRequest, err = http.NewRequest("POST", url, body); err != nil {
            return nil, errors.New("Failed to create http request for " + c.commandName() + " " + err.Error())
        }
        break
    }

    // Setup headers
    if request.HeaderParamFunc != nil {
        headers := request.HeaderParamFunc();
        if headers != nil {
            for key, value := range headers {
                httpRequest.Header.Set(key, value)
            }
        }
    }

    // Setup default content-type (if missing)
    if httpRequest.Header.Get("Content-Type") == "" && httpRequest.Header.Get("content-type") == "" {
        httpRequest.Header.Add("Content-Type", "application/json");
    }

    // Make http call
    var httpResponse *http.Response;
    if httpResponse, err = http.DefaultClient.Do(httpRequest.WithContext(ctx)); err != nil {
        return nil, errors.New("Http request failed with error " + c.commandName() + " " + err.Error())
    }
    defer httpResponse.Body.Close()

    // Read response body
    var body []byte;
    if httpResponse.Body != nil {
        if body, err = ioutil.ReadAll(httpResponse.Body); err != nil {
            return nil, errors.New("Failed to read http response" + c.commandName() + " " + err.Error())
        }
    }

    if request.ResultFunc != nil && body != nil && len(body) > 0 {
        result, err = request.ResultFunc(body)
        c.Logger.Log(requestId, c.commandName(), "ResultFunc Output", result, "ResultFunc Error", err)
        return result, err
    }

    return body, nil
}
