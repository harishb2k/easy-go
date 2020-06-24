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
