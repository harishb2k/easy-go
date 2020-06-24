package emissary

import (
    "errors"
    "github.com/afex/hystrix-go/hystrix"
    . "github.com/harishb2k/gox-base"
    . "github.com/harishb2k/easy-go/errors"
    "time"
)

type HystrixHttpCommand struct {
    HttpCommand
}

func NewHystrixHttpCommand(Service *Service, Api *Api, logger Logger) (*HystrixHttpCommand) {
    httpCommand := HystrixHttpCommand{
        HttpCommand: NewHttpCommand(
            Service,
            Api,
            logger,
        ),
    }
    httpCommand.Setup(logger)
    return &httpCommand
}

func (a *Api) getRequestTimeoutWithBuffer() (timeout int) {
    timeout = a.RequestTimeout
    delta := timeout / 10
    if delta < 10 {
        timeout = timeout + 10
    } else {
        timeout = timeout + delta
    }
    return timeout
}

// Setup a command at boot time
func (c *HystrixHttpCommand) Setup(logger Logger) (err error) {
    hystrix.ConfigureCommand(
        c.commandName(),
        hystrix.CommandConfig{
            Timeout:               c.Api.getRequestTimeoutWithBuffer(),
            MaxConcurrentRequests: c.Api.MaxRequestQueueSize,
            ErrorPercentThreshold: 1,
        },
    )
    return
}

// Execute a request
func (c *HystrixHttpCommand) Execute(request *Request) (response *Response, err error) {
    errorChannel := make(chan Error, 1)
    responseOutputChannel := make(chan *Response, 1)
    hystrixError := hystrix.Go(c.commandName(), func() (error) {
        if result, err := c.HttpCommand.Execute(request); err != nil {
            errorChannel <- err
            return err
        } else {
            responseOutputChannel <- result
            return nil
        }
    }, nil)

    select {
    case out := <-responseOutputChannel:
        close(responseOutputChannel)
        close(errorChannel)
        return out, nil

    case err := <-errorChannel:
        c.Error("HystrixHttpCommand: error to run command - ", "command=", c.commandName(), "error=", err)
        close(responseOutputChannel)
        close(errorChannel)
        return nil, err

    case _ = <-time.After(time.Duration(c.Api.getRequestTimeoutWithBuffer()+10) * time.Millisecond):
        c.Error("HystrixHttpCommand: time based timeout - ", "command=", c.commandName(), "error=", err)
        return nil, &ErrorObj{
            Err:         ErrHystrixTimeout,
            Name:        hystrixErrorToInternalError(hystrix.ErrTimeout).Error(),
            Description: "Got hystrix error = " + hystrixErrorToInternalError(hystrix.ErrTimeout).Error(),
            Object:      &Response{StatusCode: 500, Status: "Unknown"},
        }

    case err := <-hystrixError:
        c.Error("HystrixHttpCommand: error to run command - ", "command=", c.commandName(), "error=", err)
        return nil, &ErrorObj{
            Err:         hystrixErrorToInternalError(err),
            Name:        hystrixErrorToInternalError(err).Error(),
            Description: "Got hystrix error = " + hystrixErrorToInternalError(err).Error(),
            Object:      &Response{StatusCode: 500, Status: "Unknown"},
        }
    }
}

func hystrixErrorToInternalError(err error) (error) {
    if errors.Is(err, hystrix.ErrMaxConcurrency) {
        return ErrHystrixRejection
    } else if errors.Is(err, hystrix.ErrCircuitOpen) {
        return ErrHystrixCircuitOpen
    } else if errors.Is(err, hystrix.ErrTimeout) {
        return ErrHystrixTimeout
    } else {
        return ErrHystrixUnknown
    }
}
