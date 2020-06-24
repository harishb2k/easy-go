package examples

import (
    "fmt"
    "github.com/harishb2k/gox-base"
    . "github.com/harishb2k/easy-go/emissary"
    . "github.com/harishb2k/easy-go/test_http"
    "github.com/harishb2k/easy-go/tools"
    "gopkg.in/yaml.v3"
    "io/ioutil"
)

var config = localEmissaryConfiguration{}
// var logger = DefaultLogger{}
var logger = base.NoOpLogger{}

type localEmissaryConfiguration struct {
    EmissaryConfiguration *Configuration `yaml:"emissaryConfiguration"`
}

func setupTest() {
    data, err := ioutil.ReadFile("./examples/exampledata/app.yml")
    if err == nil {
        err = yaml.Unmarshal(data, &config)
        if err != nil {
            panic(err)
        }
    } else {
        panic(err)
    }

    tools.RunServer("12345")
}

func EmissaryMain() {
    setupTest()

    ctx, _ := NewContext(
        config.EmissaryConfiguration,
        logger,
    )

    response, err := ctx.Execute(
        "local", "simpleGetMethod",
        &Request{
            PathParam:  map[string]interface{}{"id": 1},
            Body:       &TestServerObject{StringValue: "TestHystrixHttpCommand_ActualServer"},
            ResultFunc: DefaultJsonResultFunc(&TestServerObject{}),
        },
    )
    if err == nil {
        fmt.Println(response.FormattedDebugString())
    } else {
        fmt.Println(err)
    }

}
