package emissary

import (
    "github.com/harishb2k/easy-go/basic"
    . "github.com/harishb2k/easy-go/test_http"
    "github.com/jarcoal/httpmock"
    "github.com/smartystreets/assertions"
    "github.com/stretchr/testify/assert"
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
                "id": 100,
                "title": "testme", 
                "completed": false
        }`

var config = localConfig{}

func setupTest() {
    data, err := ioutil.ReadFile("./testdata/app.yml")
    if err == nil {
        err = yaml.Unmarshal(data, &config)
        if err != nil {
            panic(err)
        }
    }
}

func TestHttpCommand(t *testing.T) {
    setupTest()

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Get service and api from config
    service := config.EmissaryConfiguration.ServiceList["serviceA"]
    api := service.ApiList["update"]
    service.Name = "serviceA"
    api.Name = "update"

    // Setup dummy http response
    SetupMockHttpResponse(
        HttpMockSpec{
            Url:         "http://jsonplaceholder.typicode.com:80/todos/1",
            Data:        dummyHttpResponseString,
            ResponseObj: &dummyHttpResponse{},
        },
    )

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
    assertions.ShouldEqual(100, result.ID)
    assertions.ShouldEqual("testme", result.Title)
}

func TestHttpCommandWithError400(t *testing.T) {
    setupTest()

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Get service and api from config
    service := config.EmissaryConfiguration.ServiceList["serviceB"]
    api := service.ApiList["updateWith400"]
    service.Name = "serviceB"
    api.Name = "updateWith400"

    // Setup dummy http response
    SetupMockHttpResponse(
        HttpMockSpec{
            Url:         "http://jsonplaceholder.typicode.com:80/todos/1",
            Data:        dummyHttpResponseString,
            ResponseObj: &dummyHttpResponse{},
            StatusCode:  400,
        },
    )

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
    assertions.ShouldEqual(100, result.ID)
    assertions.ShouldEqual(400, response.StatusCode)
}

func TestHttpCommandWithTimeout(t *testing.T) {
    setupTest()

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Get service and api from config
    service := config.EmissaryConfiguration.ServiceList["serviceA"]
    api := service.ApiList["update"]
    service.Name = "serviceA"
    api.Name = "update"
    api.RequestTimeout = 10

    // Setup dummy http response
    SetupMockHttpResponse(
        HttpMockSpec{
            Url:         "http://jsonplaceholder.typicode.com:80/todos/1",
            Data:        dummyHttpResponseString,
            ResponseObj: &dummyHttpResponse{},
            StatusCode:  200,
            Delay:       500,
        },
    )

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
    assert.Error(t, err, "we must get error from http command")
    assert.Nil(t, response, "We must have nil response - api timed out")
}

func TestHystrixHttpCommand(t *testing.T) {
    setupTest()

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Get service and api from config
    service := config.EmissaryConfiguration.ServiceList["serviceA"]
    api := service.ApiList["update"]
    service.Name = "serviceA"
    api.Name = "update"

    // Setup dummy http response
    SetupMockHttpResponse(
        HttpMockSpec{
            Url:         "http://jsonplaceholder.typicode.com:80/todos/1",
            Data:        dummyHttpResponseString,
            ResponseObj: &dummyHttpResponse{},
        },
    )

    // Make http command and set it up
    httpCommand := NewHystrixHttpCommand(
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
    assertions.ShouldEqual(100, result.ID)
    assertions.ShouldEqual("testme", result.Title)
}

func TestHystrixHttpCommandWithTimeout(t *testing.T) {
    setupTest()

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Get service and api from config
    service := config.EmissaryConfiguration.ServiceList["serviceA"]
    api := service.ApiList["update"]
    service.Name = "serviceA"
    api.Name = "update"
    api.RequestTimeout = 5

    // Setup dummy http response
    SetupMockHttpResponse(
        HttpMockSpec{
            Url:         "http://jsonplaceholder.typicode.com:80/todos/1",
            Data:        dummyHttpResponseString,
            ResponseObj: &dummyHttpResponse{},
            Delay:       500,
        },
    )

    // Make http command and set it up
    httpCommand := NewHystrixHttpCommand(
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
    assert.Error(t, err, "we must get error from hystrix http command")
    assert.Nil(t, response, "We must have nil response - api timed out")
}
