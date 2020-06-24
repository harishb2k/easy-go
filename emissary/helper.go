package emissary

import (
    "encoding/json"
    . "github.com/harishb2k/gox-errors"
)

// A helper to get a command using service and api name
func (ctx *Context) Get(service string, api string) (command Command, err error) {
    var commandList map[string]Command;
    var found bool

    if commandList, found = ctx.commandList[service]; !found {
        return nil, &ErrorObj{
            Name:        "service_not_found",
            Description: "service not found service=" + service + " api=" + api,
        }
    }

    var commandToReturn Command;
    if commandToReturn, found = commandList[api]; !found {
        return nil, &ErrorObj{
            Name:        "api_not_found",
            Description: "api not found service=" + service + " api=" + api,
        }
    }
    return commandToReturn, nil
}

// A helper to execute command using service and api name
func (ctx *Context) Execute(service string, api string, request *Request) (response *Response, err error) {
    var command Command;
    if command, err = ctx.Get(service, api); err != nil {
        return nil, err
    }
    return command.Execute(request)
}

func DefaultJsonResultFunc(obj interface{}) (ResultFunc) {
    return func(bytes []byte) (interface{}, Error) {
        if err := json.Unmarshal([]byte(bytes), obj); err != nil {
            return nil, &ErrorObj{
                Err:         err,
                Name:        "failed_to_build_object",
                Description: "Failed to convert byte body to requested object type",
            }
        }
        return obj, nil
    }
}

func DefaultBodyFunc(obj interface{}) (BodyFunc) {
    return func() ([]byte, error) {
        if bytes, err := json.Marshal(obj); err != nil {
            return nil, &ErrorObj{
                Err:         err,
                Name:        "failed_to_build_object",
                Description: "Failed to convert byte body to requested object type",
            }
        } else {
            return bytes, nil
        }
    }
}
