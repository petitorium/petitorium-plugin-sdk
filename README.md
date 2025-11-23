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

func (p *MyPlugin) HookFuncs() map[types.HookType]types.PluginHook {
    return map[types.HookType]types.PluginHook{
        types.PreRequest:  p.handlePreRequest,
        types.PostReceive: p.handlePostReceive,
    }
}

func (p *MyPlugin) handlePreRequest(ctx *types.HookContext) error {
    // Your pre-request logic here
    return nil
}

func (p *MyPlugin) handlePostReceive(ctx *types.HookContext) error {
    // Your post-receive logic here
    return nil
}

// Export the plugin instance
var Plugin types.Plugin = &MyPlugin{}
```

4. Build your plugin:

```bash
go build -buildmode=plugin -o my-plugin.so .
```

### For Petitorium Core Developers

Import the SDK in your main application:

```go
import "github.com/petitorium/petitorium-plugin-sdk/types"
```

Use the interfaces to load and manage plugins:

```go
// Load a plugin
plugin, err := plugin.Open("path/to/plugin.so")
if err != nil {
    return err
}

// Lookup the plugin symbol
sym, err := plugin.Lookup("Plugin")
if err != nil {
    return err
}

// Cast to the interface
petitoriumPlugin, ok := sym.(types.Plugin)
if !ok {
    return fmt.Errorf("invalid plugin type")
}

// Use the plugin
hooks := petitoriumPlugin.Hooks()
hookFuncs := petitoriumPlugin.HookFuncs()
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
}
```

## Versioning

This SDK follows [Semantic Versioning](https://semver.org/). The version is independent of Petitorium's version to allow for independent evolution of the plugin API.

### Compatibility Promise

- **Major version changes** (1.x.x → 2.0.0): Breaking changes to interfaces
- **Minor version changes** (1.1.x → 1.2.0): New features, backward compatible
- **Patch version changes** (1.1.1 → 1.1.2): Bug fixes only

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

See the [examples directory](https://github.com/petitorium/petitorium-plugin-sdk/tree/main/examples) for complete plugin examples:

- **Request Logger** - Logs all requests and responses
- **Auth Injector** - Adds authentication headers
- **Response Transformer** - Modifies response data

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This SDK is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Support

- 📖 [Documentation](https://github.com/petitorium/petitorium-plugin-sdk/wiki)
- 🐛 [Issue Tracker](https://github.com/petitorium/petitorium-plugin-sdk/issues)
- 💬 [Discussions](https://github.com/petitorium/petitorium-plugin-sdk/discussions)

## Related Projects

- [Petitorium](https://github.com/petitorium/petitorium) - The main API client application
- [Petitorium Plugin Registry](https://github.com/petitorium/petitorium-plugins) - Community plugin registry

