// Package types provides the core interfaces and types for Petitorium plugins.
// This package is designed to be imported by both the Petitorium application
// and individual plugins to ensure interface compatibility.
package types

// Plugin represents a loaded plugin that can be used with Petitorium.
// All plugins must implement this interface to be loaded by the plugin manager.
type Plugin interface {
	// Name returns the unique name of the plugin
	Name() string

	// Version returns the semantic version of the plugin (e.g., "1.0.0")
	Version() string

	// Description returns a brief description of what the plugin does
	Description() string

	// Hooks returns the list of hook types this plugin implements
	Hooks() []HookType

	// HookFuncs returns a map of hook functions keyed by hook type
	HookFuncs() map[HookType]PluginHook
}
