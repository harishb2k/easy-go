package emissary

import (
    "github.com/afex/hystrix-go/hystrix"
    "github.com/harishb2k/easy-go/basic"
)

type HystrixHttpCommand struct {
    HttpCommand
}

func NewHystrixHttpCommand(Service Service, Api Api, logger basic.Logger) (*HystrixHttpCommand) {
    httpCommand := HystrixHttpCommand{
        HttpCommand: NewHttpCommand(
            Service,
            Api,
            logger,
        ),
    }
    return &httpCommand
}

// Setup a command at boot time
func (c *HystrixHttpCommand) Setup(logger basic.Logger) (err error) {
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
func (c *HystrixHttpCommand) Execute(request *Request) (response *Response, err error) {

    _error := make(chan error, 1)
    _output := make(chan *Response, 1)
    hystrixError := hystrix.Go(c.commandName(), func() (error) {
        if result, err := c.HttpCommand.Execute(request); err != nil {
            _error <- err
        } else {
            _output <- result
        }
        return err
    }, nil)

    select {
    case out := <-_output:
        close(_output)
        close(_error)
        return out, nil

    case err := <-_error:
        c.Error("HystrixHttpCommand: error to run command - ", "command=", c.commandName(), "error=", err)
        close(_output)
        close(_error)
        return nil, err

    case err := <-hystrixError:
        c.Error("HystrixHttpCommand: error to run command - ", "command=", c.commandName(), "error=", err)
        close(_output)
        close(_error)
        return nil, err
    }
}
