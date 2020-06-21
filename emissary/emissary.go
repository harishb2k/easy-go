package emissary

import "fmt"

// Client Config e.g. client name etc for logging
type ClientConfiguration struct {
    ClientName            string `yaml:"clientName"`
    EnableNewRelicLogging bool   `yaml:"enableNewRelicLogging"`
    NewRelicLoggingPrefix string `yaml:"newRelicLoggingPrefix"`
}

// This is a single API in a service
type Api struct {
    Name                     string
    Method                   string `yaml:"method"`
    Path                     string `yaml:"path"`
    AcceptableResponseCodes  []int  `yaml:"acceptableResponseCodes"`
    ConnectionRequestTimeout int    `yaml:"connectionRequestTimeout"`
    Concurrency              int    `yaml:"concurrency"`
    MaxRequestQueueSize      int    `yaml:"maxRequestQueueSize"`
    RequestTimeout           int    `yaml:"requestTimeout"`
    ElixirEnabled            bool   `yaml:"elixirEnabled"`
}

// This is a service which will have a list of APIs
type Service struct {
    Name    string
    Type    string         `yaml:"type"`
    Host    string         `yaml:"host"`
    Port    int            `yaml:"port"`
    ApiList map[string]Api `yaml:"apis"`
}

// This is the main emissary configuration
type Configuration struct {
    ClientConfiguration ClientConfiguration `yaml:"clientConfiguration"`
    ServiceList         map[string]Service  `yaml:"services"`
}

type ILogger interface {
    Log(l ...interface{})
}

type DefaultLogger struct {
}

func (dl *DefaultLogger) Log(l ...interface{}) {
    fmt.Println(l...)
}

type NoOpLogger struct {
}

func (dl *NoOpLogger) Log(l ...interface{}) {
}
