package emissary

import (
    "encoding/json"
    "errors"
    "strconv"
)

// A helper to get a command using service and api name
func (ctx *CommandContext) Get(service string, api string) (command Command, err error) {
    var commandList map[string]Command;
    var found bool

    if commandList, found = ctx.commandList[service]; !found {
        return nil, errors.New("Command not found service=" + service + " api=" + api)
    }

    var commandToReturn Command;
    if commandToReturn, found = commandList[api]; !found {
        return nil, errors.New("Command not found service=" + service + " api=" + api)
    }
    return commandToReturn, nil
}

// A helper to execute command using service and api name
func (ctx *CommandContext) Execute(service string, api string, request CommandRequest) (result interface{}, err error) {
    var command Command;
    if command, err = ctx.Get(service, api); err != nil {
        return nil, err
    }
    return command.Execute(request)
}

func DefaultJsonResultFunc(obj interface{}) (ResultFunc) {
    return func(bytes []byte) (interface{}, error) {
        if err := json.Unmarshal([]byte(bytes), obj); err != nil {
            return nil, err
        }
        return obj, nil
    }
}

func DefaultObjectToBodyFunc(obj interface{}) (BodyFunc) {
    return func() []byte {
        if out, err := json.Marshal(obj); err != nil {
            return out
        }
        return nil
    }
}

func DefaultPathParamFunc(qp ...interface{}) (PathParamFunc) {
    return func() map[string]interface{} {
        var ret = map[string]interface{}{}
        length := len(qp)
        for i := 0; i < length; i += 2 {
            ret[Stringify(qp[i])] = qp[i+1]
        }
        return ret
    }
}

func Stringify(input interface{}) (string) {
    switch v := input.(type) {
    case int:
        return strconv.Itoa(v)
    case string:
        return v;
    default:
        return input.(string);
    }
}

type NoOpLogger struct {
}

func (dl NoOpLogger) Log(l ...interface{}) {
}
