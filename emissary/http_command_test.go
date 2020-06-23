package emissary

import (
    "github.com/harishb2k/easy-go/basic"
    . "github.com/harishb2k/easy-go/test_http"
    "github.com/jarcoal/httpmock"
    "github.com/smartystreets/assertions"
    "gopkg.in/yaml.v3"
    "io/ioutil"
    "testing"
)

type localConfig struct {
    EmissaryConfiguration Configuration `yaml:"emissaryConfiguration"`
}

type dummyHttpResponse struct {
    UserID    int    `json:"userId"`
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

var dummyHttpResponseString = `{
                "userId": 1,  
                "id": 1,
                "title": "testme", 
                "completed": false
        }`

var config = localConfig{}

func init() {
    data, err := ioutil.ReadFile("./testdata/app.yml")
    if err == nil {
        err = yaml.Unmarshal(data, &config)
        if err != nil {
            panic(err)
        }
    }
}

func TestHttpCommand(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Get service and api from config
    service := config.EmissaryConfiguration.ServiceList["serviceA"]
    api := service.ApiList["update"]
    service.Name = "serviceA"
    api.Name = "update"

    // Setup dummy http response
    dummyHttpResponseObj := dummyHttpResponse{}
    SetupGetHttpResponse("http://jsonplaceholder.typicode.com:80/todos/1", dummyHttpResponseString, &dummyHttpResponseObj)

    // Make http command and set it up
    httpCommand := NewHttpCommand(
        service,
        api,
        basic.DefaultLogger{},
    )

    // Make Http call
    response, err := httpCommand.Execute(
        &Request{
            PathParam:  map[string]interface{}{"id": 1},
            ResultFunc: DefaultJsonResultFunc(&dummyHttpResponse{}),
        },
    )
    assertions.ShouldBeNil(err)
    assertions.ShouldNotBeNil(response)

    // verify result
    result, ok := response.Result.(*dummyHttpResponse)
    assertions.ShouldBeTrue(ok)
    assertions.ShouldNotBeNil(result)
}
