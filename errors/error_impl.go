package errors

import (
    "fmt"
)

type ErrorObj struct {
    Name        string
    Description string
    Err         error
    Object      interface{}
}

func (e *ErrorObj) Error() string {
    return fmt.Sprintf("Name=%s, Description=[%s] Err=[%v] Object=[%v]", e.Name, e.Description, e.Err, e.Object)
}

func (e *ErrorObj) FormattedDebugString() interface{} {
    return fmt.Sprintf("Name=%s \nDescription=%s \nErr=%v \nObject=%v", e.Name, e.Description, e.Err, e.Object)
}
