# dCloud Topology Builder Go Client

A `Go` client for [Cisco dCloud](https://dcloud.cisco.com) Topology Builder

This library contains an SDK for integrating `Go` applications with Cisco dCloud Topology Builder

## Building

Build with the following command:

`go build ./tbclient`

## Usage

To get a pointer to the client use `tbclient.NewClient`

```go
apiUrl := "https://tbv3-production.ciscodcloud.com/api/"
authToken := "xyz" // Use real Auth Token from Cisco SSO

client := tbclient.NewClient(&apiUrl, &authToken)
```
The client has a debug mode which if turned on will output debug information on its http calls

```go
client.Debug = true
```
Applications should declare and use their own unique agent value in the client

```go
client.UserAgent = "my-application/v1.0.0"
```

The client exposes functions for interacting with the various Topology Builder Resources

For example to retrieve all the available Topologies:
```go
topologies, err := client.GetAllTopologies()
```

## Testing

There is a single test class `tbclient_test.go` which uses wiremock mappings declared within the `test_stubs` directory.

The test suite starts wiremock in a [Docker](https://www.docker.com/) container and Docker must be running locally in order to run the tests.

The wiremock mappings are based on Topology Builder API Contract Stubs which are internal to Cisco.  A Cisco employee with access to the contracts will be able to update the test stubs with any new functionality as it is released.

Use the following command to execute the test suite:

`go test ./tbclient`

Copyright (c) 2023 Cisco Systems, Inc. and its affiliates