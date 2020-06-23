package emissary

import (
    "encoding/json"
    "github.com/harishb2k/easy-go/easy"
)

func DefaultJsonResultFunc(obj interface{}) (ResultFunc) {
    return func(bytes []byte) (interface{}, easy.Error) {
        if err := json.Unmarshal([]byte(bytes), obj); err != nil {
            return nil, &easy.ErrorObj{
                Err:         err,
                Name:        "failed_to_build_object",
                Description: "Failed to convert byte body to requested object type",
            }
        }
        return obj, nil
    }
}
