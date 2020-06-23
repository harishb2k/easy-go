package emissary

import "encoding/json"

func DefaultJsonResultFunc(obj interface{}) (ResultFunc) {
    return func(bytes []byte) (interface{}, error) {
        if err := json.Unmarshal([]byte(bytes), obj); err != nil {
            return nil, err
        }
        return obj, nil
    }
}
