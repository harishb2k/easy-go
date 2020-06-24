package easy

import (
    "encoding/json"
    "fmt"
    "strconv"
)

func Stringify(input interface{}) (string) {
    result, _ := StringifyWithError(input)
    return result
}

func Objectify(data []byte, obj interface{}) (err error) {
    return json.Unmarshal(data, &obj)
}

func ObjectifyString(data string, obj interface{}) (err error) {
    return Objectify([]byte(data), obj)
}

func StringifyWithError(input interface{}) (out string, err error) {
    switch v := input.(type) {
    case int:
        out = strconv.Itoa(v)

    case int8, int16, int32, int64:
        out = fmt.Sprintf("%d", v)

    case bool:
        if v == true {
            out = "true"
        } else {
            out = "false"
        }
    case string:
        out = v

    default:
        if _out, err := json.Marshal(v); err != nil {
            out = ""
        } else {
            out = string(_out)
        }
    }
    return
}

func BytesWithError(input interface{}) ([]byte, error) {
    switch v := input.(type) {
    default:
        out, err := json.Marshal(v)
        if err != nil {
            return nil, err
        }
        return out, err
    }
}
