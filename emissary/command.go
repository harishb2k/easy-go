package emissary

import (
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
    Setup(logger ILogger) (err error)

    // Execute command with given info
    Execute(request CommandRequest) (interface{}, error)
}

// This is a complete Emissary context - this context contains all service and api commands
type CommandContext struct {
    commandList    map[string]map[string]Command
    DisableLogging bool
    logger         ILogger
}

func NewCommandContext(configuration *Configuration, logger ILogger) (ctx CommandContext, err error) {
    ctx = CommandContext{
        commandList: map[string]map[string]Command{},
        logger:      logger,
    }
    ctx.Setup(configuration)
    return ctx, nil
}

// This will read YAML and make all Http Commands
func (ctx *CommandContext) Setup(configuration *Configuration) (err error) {
    for serviceName, service := range configuration.ServiceList {
        service.Name = serviceName
        ctx.commandList[serviceName] = map[string]Command{};
        for apiName, api := range service.ApiList {
            api.Name = apiName
            ctx.commandList[serviceName][apiName] = NewHystrixHttpCommand(service, api, ctx.logger)
        }
    }
    return nil
}
