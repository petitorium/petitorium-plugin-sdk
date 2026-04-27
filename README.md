# Petitorium Plugin SDK

The official Plugin Software Development Kit (SDK) for Petitorium. This module provides the shared interfaces and types that enable plugin development for the Petitorium API client.

## Overview

This SDK serves as the contract between Petitorium and its plugins, ensuring compatibility and providing a stable API for plugin developers.

## Installation

```bash
go get github.com/petitorium/petitorium-plugin-sdk
```

## Quick Start

### For Plugin Developers

1. Create a new Go module for your plugin:

```bash
go mod init github.com/yourusername/petitorium-plugin-my-plugin
```

2. Add the SDK dependency:

```bash
go get github.com/petitorium/petitorium-plugin-sdk
```

3. Implement the Plugin interface:

```go
package main

import (
    "github.com/hashicorp/go-plugin"
    "github.com/petitorium/petitorium-plugin-sdk/shared"
    "github.com/petitorium/petitorium-plugin-sdk/types"
)

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}

func (p *MyPlugin) Description() string {
    return "My custom Petitorium plugin"
}

func (p *MyPlugin) Hooks() []types.HookType {
    return []types.HookType{types.PreRequest, types.PostReceive}
}

func (p *MyPlugin) ExecuteHook(hookType types.HookType, ctx *types.HookContext) (*types.HookContext, error) {
    switch hookType {
    case types.PreRequest:
        // Handle pre-request logic here
    case types.PostReceive:
        // Handle post-receive logic here
    }
    return ctx, nil
}

func main() {
    plugin.Serve(&plugin.ServeConfig{
        HandshakeConfig: shared.Handshake,
        Plugins: map[string]plugin.Plugin{
            "my-plugin": &shared.PetitoriumPlugin{Impl: &MyPlugin{}},
        },
        GRPCServer: plugin.DefaultGRPCServer,
    })
}
```

4. Build your plugin as a standalone executable:

```bash
go build -o my-plugin .
# Or cross-compile for other platforms
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o my-plugin .
```

### For Petitorium Core Developers

Import the SDK and `go-plugin` in your main application:

```go
import (
    "os/exec"
    "github.com/hashicorp/go-plugin"
    "github.com/petitorium/petitorium-plugin-sdk/shared"
    "github.com/petitorium/petitorium-plugin-sdk/types"
)
```

Use `go-plugin` to load and manage plugins over gRPC:

```go
// Create an exec.Cmd to launch the plugin process
client := plugin.NewClient(&plugin.ClientConfig{
    HandshakeConfig: shared.Handshake,
    Plugins: map[string]plugin.Plugin{
        "my-plugin": &shared.PetitoriumPlugin{},
    },
    Cmd:              exec.Command("path/to/my-plugin"),
    AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
})
defer client.Kill()

// Connect via RPC
rpcClient, err := client.Client()
if err != nil {
    return err
}

// Request the plugin instance
raw, err := rpcClient.Dispense("my-plugin")
if err != nil {
    return err
}

// Cast to the interface
petitoriumPlugin, ok := raw.(types.Plugin)
if !ok {
    return fmt.Errorf("invalid plugin type")
}

// Use the plugin
hooks := petitoriumPlugin.Hooks()
ctx, err := petitoriumPlugin.ExecuteHook(types.PreRequest, &types.HookContext{})
```

## Hook Types

The SDK provides a comprehensive set of hook types that plugins can implement:

### Request Lifecycle

- `PreRequest` - Before request is sent (with template variables)
- `PostVariableSubstitution` - After environment variables are substituted
- `PreSend` - Just before sending the request
- `PostSend` - After request is sent
- `PostReceive` - After response is received
- `PostRequest` - After complete request/response cycle

### Validation

- `RequestValidation` - Validate the request before sending
- `ResponseValidation` - Validate the response after receiving

### Data Management

- `PreSave` / `PostSave` - Before/after saving data

### UI Lifecycle

- `PreUIUpdate` / `PostUIUpdate` - Before/after UI updates
- `OnUIInit` / `OnUIClose` - UI initialization/cleanup

### And more... see the full list in `types/hooks.go`

## Data Structures

### RequestData

```go
type RequestData struct {
    Method      string            // HTTP method
    URL         string            // Request URL (may contain templates)
    Headers     map[string]string // Request headers
    Body        string            // Request body
    Collection  string            // Collection name
    RequestName string            // Request identifier
}
```

### ResponseData

```go
type ResponseData struct {
    StatusCode int               // HTTP status code
    Status     string            // HTTP status text
    Headers    map[string]string // Response headers
    Body       string            // Response body
    Duration   int64             // Request duration (ms)
}
```

### HookContext

```go
type HookContext struct {
    Request     *RequestData      // HTTP request data
    Response    interface{}       // HTTP response data
    Environment map[string]string // Environment variables
    Config      map[string]interface{} // Plugin config
    Workspace   string            // Active workspace name
}
```

## Versioning

This SDK follows [Semantic Versioning](https://semver.org/). The version is independent of Petitorium's version to allow for independent evolution of the plugin API.

## Best Practices

### For Plugin Developers

1. **Version your plugins** - Use semantic versioning
2. **Handle errors gracefully** - Return meaningful errors from hooks
3. **Don't modify request data unless necessary** - Be a good citizen
4. **Document your hooks** - Let users know what your plugin does
5. **Test thoroughly** - Test with different request types and edge cases
6. **Follow naming conventions** - Use descriptive plugin names

### For Petitorium Core

1. **Maintain backward compatibility** - Don't break existing plugins
2. **Provide good error messages** - Help plugin developers debug issues
3. **Document hook execution order** - Let developers know when hooks run
4. **Test with example plugins** - Ensure the SDK works correctly

## Examples

- [**Request Logger**](https://github.com/petitorium/petitorium-plugin-request-logger) - A comprehensive request and response logging plugin for [Petitorium](https://github.com/petitorium/petitorium) that provides detailed logging of HTTP requests and responses with support for both raw template variables and expanded environment variables.
- [**Auth Injector**](https://github.com/petitorium/petitorium-plugin-auth-injector) - An authentication injection plugin for [Petitorium](https://github.com/petitorium/petitorium) that automatically injects authentication headers into HTTP requests and captures authentication tokens from responses.
- [**Auth Retriever**](https://github.com/petitorium/petitorium-plugin-auth-retriever) - An authentication retrieval plugin for [Petitorium](https://github.com/petitorium/petitorium) that automatically captures authentication tokens from responses and stores them in environment variables for flexible usage.


## License

This SDK is licensed under the MIT License.

## Support

- 🐛 [Issue Tracker](https://github.com/petitorium/petitorium-plugin-sdk/issues)

## Related Projects

- [Petitorium](https://github.com/petitorium/petitorium) - The main API client application

