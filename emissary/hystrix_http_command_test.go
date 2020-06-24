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
    api.RequestTimeout = 1

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
    assert.Nil(t, response)
    assert.Equal(t, ErrorCodeHttpServerTimeout, err.GetName())
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
    assert.NoError(t, err)
    assert.NotNil(t, response)
    result, ok := response.Result.(*TestServerObject)
    assert.True(t, ok)
    assert.NotNil(t, result)
    assert.Equal(t, "TestHystrixHttpCommand_ActualServer Response", result.ResponseStringValue)

    if err != nil {
        logger.Debug(err.FormattedDebugString())
    } else {
        logger.Debug(response.FormattedDebugString())
    }
}

func TestHystrixHttpCommand_ActualServer_HystrixCircuitOpen(t *testing.T) {
    setupTest()

    ctx, _ := NewContext(
        config.EmissaryConfiguration,
        logger,
    )

    var hystrixErrorCircuitOpen = 0
    for i := 1; i <= 100; i++ {
        response, err := ctx.Execute(
            "local", "simpleGetMethodWithError",
            &Request{
                PathParam:  map[string]interface{}{"id": 1},
                Body:       &TestServerObject{StringValue: "TestHystrixHttpCommand_ActualServer", ErrorCodeToReturn: 500},
                ResultFunc: DefaultJsonResultFunc(&TestServerObject{}),
            },
        )
        assert.Error(t, err)
        var _ = response
        if err != nil {
            if errors.Is(err.GetError(), ErrHystrixCircuitOpen) {
                hystrixErrorCircuitOpen++
            }
        }
    }
    assert.True(t, hystrixErrorCircuitOpen > 60)
}

func TestHystrixHttpCommand_ActualServer_Post_ExpectSuccess(t *testing.T) {
    setupTest()
    ctx, _ := NewContext(
        config.EmissaryConfiguration,
        logger,
    )

    response, err := ctx.Execute(
        "local", "simplePostMethod",
        &Request{
            Header:     map[string]interface{}{"A": "B", "C": 1, "D": false},
            PathParam:  map[string]interface{}{"id": 1},
            Body:       &TestServerObject{StringValue: "TestHystrixHttpCommand_ActualServer"},
            ResultFunc: DefaultJsonResultFunc(&TestServerObject{}),
        },
    )
    var _ = response
    assert.NoError(t, err)
    assert.NotNil(t, response)
    result, ok := response.Result.(*TestServerObject)
    assert.True(t, ok)
    assert.NotNil(t, result)
    assert.Equal(t, "TestHystrixHttpCommand_ActualServer Response", result.ResponseStringValue)

    assert.Equal(t, result.Headers["A"], "B")
    assert.Equal(t, result.Headers["C"], "1")
    assert.Equal(t, result.Headers["D"], "false")
}

func TestHystrixHttpCommand_ActualServer_Post_ExpectError(t *testing.T) {
    setupTest()
    ctx, _ := NewContext(
        config.EmissaryConfiguration,
        logger,
    )

    response, err := ctx.Execute(
        "local", "simplePostMethodWithError",
        &Request{
            Header:     map[string]interface{}{"A": "B", "C": 1, "D": false},
            PathParam:  map[string]interface{}{"id": 1},
            Body:       &TestServerObject{StringValue: "TestHystrixHttpCommand_ActualServer", ErrorCodeToReturn: 500},
            ResultFunc: DefaultJsonResultFunc(&TestServerObject{}),
        },
    )
    assert.Error(t, err)
    assert.Nil(t, response)
}
