package emissary

import "github.com/harishb2k/easy-go/easy"

// Client Config e.g. client name etc for logging
type ClientConfiguration struct {
    ClientName            string `yaml:"clientName"`
    EnableNewRelicLogging bool   `yaml:"enableNewRelicLogging"`
    NewRelicLoggingPrefix string `yaml:"newRelicLoggingPrefix"`
}

// This is a single API in a service
type Api struct {
    Name                     string
    Method                   string `yaml:"method"`
    Path                     string `yaml:"path"`
    AcceptableResponseCodes  []int  `yaml:"acceptableResponseCodes"`
    ConnectionRequestTimeout int    `yaml:"connectionRequestTimeout"`
    Concurrency              int    `yaml:"concurrency"`
    MaxRequestQueueSize      int    `yaml:"maxRequestQueueSize"`
    RequestTimeout           int    `yaml:"requestTimeout"`
    ElixirEnabled            bool   `yaml:"elixirEnabled"`
}

// This is a service which will have a list of APIs
type Service struct {
    Name    string
    Type    string         `yaml:"type"`
    Host    string         `yaml:"host"`
    Port    int            `yaml:"port"`
    ApiList map[string]Api `yaml:"apis"`
}

// This is the main emissary configuration
type Configuration struct {
    ClientConfiguration ClientConfiguration `yaml:"clientConfiguration"`
    ServiceList         map[string]Service  `yaml:"services"`
}

// This is a complete Emissary context - this context contains all service and api commands
type Context struct {
    commandList    map[string]map[string]Command
    DisableLogging bool
    logger         easy.Logger
}

// This will read YAML and make all Http Commands
func (ctx *Context) Setup(configuration *Configuration) (err error) {
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

func NewContext(configuration *Configuration, logger easy.Logger) (ctx *Context, err error) {
    ctx = &Context{
        commandList: map[string]map[string]Command{},
        logger:      logger,
    }
    ctx.Setup(configuration)
    return
}
