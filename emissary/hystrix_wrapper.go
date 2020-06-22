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
        HttpCommand: HttpCommand{
            Service: Service,
            Api:     Api,
        },
    }
    httpCommand.Setup(logger)
    return &httpCommand
}

// Setup a command at boot time
func (c *HystrixHttpCommand) Setup(logger basic.Logger) (err error) {

    // Setup hystrix command
    hystrix.ConfigureCommand(
        c.commandName(),
        hystrix.CommandConfig{
            Timeout:               c.Api.RequestTimeout + (c.Api.RequestTimeout / 10),
            MaxConcurrentRequests: c.Api.MaxRequestQueueSize,
            ErrorPercentThreshold: 25,
        },
    )

    c.HttpCommand.Setup(logger)

    return nil
}

// Execute command with given info
func (c *HystrixHttpCommand) Execute(request CommandRequest) (_r interface{}, err error) {

    _output := make(chan interface{}, 1)
    hystrixError := hystrix.Go(c.commandName(), func() (error) {
        if result, err := c.HttpCommand.Execute(request); err != nil {
            close(_output)
        } else {
            _output <- result
        }
        return err
    }, nil)

    select {
    case out := <-_output:
        return out, nil

    case err := <-hystrixError:
        c.Error("Error to run", c.commandName(), err)
        return nil, err
    }
}
