package types

// HookType defines when a plugin hook executes in the request lifecycle.
// Each hook type represents a specific point where plugins can intercept
// and modify the request/response flow.
type HookType string

const (
	// Request Lifecycle Hooks
	PreRequest               HookType = "pre_request"                // Before request is sent (with template variables)
	PostVariableSubstitution HookType = "post_variable_substitution" // After environment variables are substituted
	PreSend                  HookType = "pre_send"                   // Just before sending the request
	PostSend                 HookType = "post_send"                  // After request is sent
	PostReceive              HookType = "post_receive"               // After response is received
	PostRequest              HookType = "post_request"               // After complete request/response cycle

	// Validation Hooks
	RequestValidation  HookType = "request_validation"  // Validate the request before sending
	ResponseValidation HookType = "response_validation" // Validate the response after receiving

	// Data Management Hooks
	PreSave  HookType = "pre_save"  // Before saving data (collections, environments, etc.)
	PostSave HookType = "post_save" // After saving data

	// UI Lifecycle Hooks
	PreUIUpdate  HookType = "pre_ui_update"  // Before UI is updated
	PostUIUpdate HookType = "post_ui_update" // After UI is updated
	OnUIInit     HookType = "on_ui_init"     // When UI is initialized
	OnUIClose    HookType = "on_ui_close"    // When UI is closed

	// Collection Management Hooks
	OnCollectionLoad HookType = "on_collection_load" // When a collection is loaded
	OnCollectionSave HookType = "on_collection_save" // When a collection is saved

	// Environment Management Hooks
	OnEnvironmentLoad HookType = "on_environment_load" // When an environment is loaded
	OnEnvironmentSave HookType = "on_environment_save" // When an environment is saved

	// Configuration Hooks
	OnConfigLoad HookType = "on_config_load" // When configuration is loaded
	OnConfigSave HookType = "on_config_save" // When configuration is saved

	// Error and Success Hooks
	OnError   HookType = "on_error"   // When an error occurs
	OnSuccess HookType = "on_success" // When operation succeeds

	// Advanced Request Hooks
	RequestRetry   HookType = "request_retry"   // When a request is retried
	RequestTimeout HookType = "request_timeout" // When a request times out

	// Response Processing Hooks
	ResponseTransform HookType = "response_transform" // Transform response data
	ResponseCache     HookType = "response_cache"     // Cache response data
)

// HookContext provides context data to plugin hooks at execution time.
// This struct contains all the information a plugin might need to process
// a request or response.
type HookContext struct {
	// Request contains the HTTP request data
	Request *RequestData

	// Response contains the HTTP response data (interface{} for flexibility)
	Response interface{}

	// Environment contains environment variables for variable substitution
	Environment map[string]string

	// Config contains plugin-specific configuration from the main config file
	Config map[string]interface{}
}

// PluginHook defines the function signature that all plugin hooks must implement.
// Hooks receive a HookContext and can modify data or perform side effects.
// They should return an error if something goes wrong, which may stop request processing.
type PluginHook func(ctx *HookContext) error
