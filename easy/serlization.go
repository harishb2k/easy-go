package easy

import (
    "encoding/json"
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

func StringifyWithError(input interface{}) (string, error) {
    switch v := input.(type) {
    case int:
        return strconv.Itoa(v), nil
    case bool:
        if v == true {
            return "true", nil
        } else {
            return "false", nil
        }
    case string:
        return v, nil

    default:
        out, err := json.Marshal(v)
        if err != nil {
            return "", err
        }
        return string(out), nil
    }
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
