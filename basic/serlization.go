package basic

import (
    "encoding/json"
    "strconv"
)

func Stringify(input interface{}) (string) {
    result, _ := StringifyWithError(input)
    return result
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
