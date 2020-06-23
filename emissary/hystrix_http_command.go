package emissary

import (
    "errors"
    "github.com/afex/hystrix-go/hystrix"
    "github.com/harishb2k/easy-go/easy"
)

type HystrixHttpCommand struct {
    HttpCommand
}

func NewHystrixHttpCommand(Service *Service, Api *Api, logger easy.Logger) (*HystrixHttpCommand) {
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

// Setup a command at boot time
func (c *HystrixHttpCommand) Setup(logger easy.Logger) (err error) {
    hystrix.ConfigureCommand(
        c.commandName(),
        hystrix.CommandConfig{
            Timeout:               c.Api.RequestTimeout + (c.Api.RequestTimeout / 10),
            MaxConcurrentRequests: c.Api.MaxRequestQueueSize,
            ErrorPercentThreshold: 25,
        },
    )
    return
}

// Execute a request
func (c *HystrixHttpCommand) Execute(request *Request) (response *Response, err easy.Error) {
    errorChannel := make(chan easy.Error, 1)
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

    case err := <-hystrixError:
        c.Error("HystrixHttpCommand: error to run command - ", "command=", c.commandName(), "error=", err)
        return nil, &easy.ErrorObj{
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
