package easy

type Error interface {
    error
    GetName() string
    GetDescription() string
    GetError() error
    GetObject() interface{}
}

type ErrorObj struct {
    Name        string
    Description string
    Err         error
    Object      interface{}
}

func (e ErrorObj) Error() string {
    if e.Description != "" {
        return e.Description
    } else if e.Name != "" {
        return e.Name
    } else if e.Err != nil {
        return e.Err.Error()
    }
    return "Unknown"
}

func (e *ErrorObj) GetName() string {
    return e.Name
}

func (e *ErrorObj) GetDescription() string {
    return e.Description
}

func (e *ErrorObj) GetError() error {
    return e.Err
}

func (e *ErrorObj) GetObject() interface{} {
    return e.Object
}
