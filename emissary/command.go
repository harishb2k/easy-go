package emissary

import (
    "errors"
    "github.com/jwangsadinata/go-multimap/slicemultimap"
)

// This function returns a multi-map to provide Query params for HTTP request
type QueryParamFunc func() (*slicemultimap.MultiMap)

// This function returns a multi-map to provide Path params for HTTP request
type PathParamFunc func() (map[string]interface{})

// This function returns a map to provide headers for HTTP request
type HeaderParamFunc func() (map[string]string)

// This function returns body for HTTP request
type BodyFunc func() ([]byte)

// This function returns a user response object from Http response
type ResultFunc func([]byte) (interface{}, error)

// A HTTP request
type CommandRequest struct {
    QueryParamFunc  QueryParamFunc;
    PathParamFunc   PathParamFunc;
    HeaderParamFunc HeaderParamFunc;
    BodyFunc        BodyFunc;
    ResultFunc      ResultFunc;
}

// This is a command interface. A valid implementation would be HTTP
type Command interface {
    // Setup a command at boot time
    Setup() (err error)

    // Execute command with given info
    Execute(request CommandRequest) (interface{}, error)
}

// This is a complete Emissary context - this context contains all service and api commands
type CommandContext struct {
    CommandList    map[string]map[string]Command
    DisableLogging bool
}

func (ctx *CommandContext) Setup(configuration *Configuration) (err error) {

    ctx.CommandList = map[string]map[string]Command{}

    for serviceName, service := range configuration.ServiceList {
        service.Name = serviceName

        var commandMap = map[string]Command{};
        for apiName, api := range service.ApiList {
            api.Name = apiName

            var logger ILogger
            if ctx.DisableLogging {
                logger = &NoOpLogger{}
            }

            // Command setup
            httpCommand := HttpCommand{
                Service: service,
                Api:     api,
                Logger:  logger,
            }
            httpCommand.Setup()

            commandMap[apiName] = &httpCommand
        }

        ctx.CommandList[serviceName] = commandMap
    }

    return nil
}

// A helper to get a command using service and api name
func (ctx *CommandContext) Get(service string, api string) (command Command, err error) {
    var commandList map[string]Command;
    var found bool

    if commandList, found = ctx.CommandList[service]; !found {
        return nil, errors.New("Command not found service=" + service + " api=" + api)
    }

    var commandToReturn Command;
    if commandToReturn, found = commandList[api]; !found {
        return nil, errors.New("Command not found service=" + service + " api=" + api)
    }
    return commandToReturn, nil
}
