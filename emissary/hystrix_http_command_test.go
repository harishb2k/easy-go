package emissary

import (
    "errors"
    "github.com/afex/hystrix-go/hystrix"
    . "github.com/harishb2k/easy-go/test_http"
    "github.com/harishb2k/easy-go/tools"
    "github.com/jarcoal/httpmock"
    "github.com/stretchr/testify/assert"
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

func init() {
    tools.RunServer("12345")
}

func TestHystrixHttpCommand_ExpectError_WithTimeout(t *testing.T) {
    setupTest()
    hystrix.Flush()

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
        logger,
    )

    // Make Http call
    response, err := httpCommand.Execute(
        &Request{
            PathParam:  map[string]interface{}{"id": 1},
            ResultFunc: DefaultJsonResultFunc(&dummyHttpResponse{}),
        },
    )
    assert.Error(t, err)
    assert.True(t, errors.Is(err.GetError(), ErrHystrixTimeout))
    assert.Nil(t, response)
    obj, ok := err.GetObject().(*Response)
    assert.True(t, ok)
    assert.Equal(t, 500, obj.StatusCode)
}

func TestHystrixHttpCommand_CircuitOpen(t *testing.T) {
    setupTest()
    hystrix.Flush()

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Get service and api from config
    service := config.EmissaryConfiguration.ServiceList["serviceB"]
    api := service.ApiList["updateWithCircuitOpenCase"]
    api.RequestTimeout = 200
    api.MaxRequestQueueSize = 2

    // Setup dummy http response
    SetupMockHttpResponse(
        HttpMockSpec{
            Url:         "http://jsonplaceholder.typicode.com:80/todos/1",
            Data:        dummyHttpResponseString,
            ResponseObj: &dummyHttpResponse{},
            Delay:       100,
        },
    )

    ctx, _ := NewContext(
        config.EmissaryConfiguration,
        logger,
    )

    var wg sync.WaitGroup
    var errHystrixRejectionCount int32 = 0
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go func() {
            time.Sleep(100 * time.Millisecond)
            response, err := ctx.Execute(
                "serviceB", "updateWithCircuitOpenCase",
                &Request{
                    PathParam:  map[string]interface{}{"id": 1},
                    ResultFunc: DefaultJsonResultFunc(&dummyHttpResponse{}),
                },
            )
            if err == nil {
                // fmt.Println(response.FormattedDebugString())
                var _ = response
            } else {
                // fmt.Println(err.FormattedDebugString())
                if errors.Is(err.GetError(), ErrHystrixRejection) {
                    atomic.AddInt32(&errHystrixRejectionCount, 1)
                }
            }
            wg.Done()
        }()
    }
    wg.Wait()
    assert.Equal(t, int32(8), errHystrixRejectionCount)
}

func TestHystrixHttpCommand_ActualServer(t *testing.T) {
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

    if err != nil {
        logger.Debug(err.FormattedDebugString())
    } else {
        logger.Debug(response.FormattedDebugString())
    }
}

func TestHystrixHttpCommand_ActualServer_HystrixCurcitOpen(t *testing.T) {
    setupTest()

    ctx, _ := NewContext(
        config.EmissaryConfiguration,
        logger,
    )

    response, err := ctx.Execute(
        "local", "simpleGetMethodWithError",
        &Request{
            PathParam:  map[string]interface{}{"id": 1},
            Body:       &TestServerObject{StringValue: "TestHystrixHttpCommand_ActualServer", ErrorCodeToReturn:500},
            ResultFunc: DefaultJsonResultFunc(&TestServerObject{}),
        },
    )

    if err != nil {
        logger.Debug(err.FormattedDebugString())
    } else {
        logger.Debug(response.FormattedDebugString())
    }
}
