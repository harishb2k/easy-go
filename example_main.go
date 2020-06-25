package main

import (
    "fmt"
    "github.com/harishb2k/easy-go/examples"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

func main() {

    fmt.Println("\n\n\n------------------ Start: ScyllaMain ------------------")
    examples.ScyllaMain()

    fmt.Println("\n\n\n------------------ Start: EmissaryMain ----------------")
    examples.EmissaryMain()

    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":2112", nil)
}
