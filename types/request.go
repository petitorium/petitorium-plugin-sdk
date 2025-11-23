package types

// RequestData represents an HTTP request being processed by Petitorium.
// This structure contains all the information about an outgoing request.
type RequestData struct {
	// Method is the HTTP method (GET, POST, PUT, DELETE, etc.)
	Method string

	// URL is the request URL (may contain template variables like {{protocol}}{{domain}})
	URL string

	// Headers contains the request headers
	Headers map[string]string

	// Body contains the request body (if any)
	Body string

	// Collection is the name of the collection this request belongs to
	Collection string

	// RequestName is the name/identifier of this specific request
	RequestName string
}

// ResponseData represents an HTTP response received from the server.
// This provides a consistent interface for plugins to access response information.
type ResponseData struct {
	// StatusCode is the HTTP status code (200, 404, 500, etc.)
	StatusCode int

	// Status is the HTTP status text ("OK", "Not Found", "Internal Server Error", etc.)
	Status string

	// Headers contains the response headers
	Headers map[string]string

	// Body contains the response body
	Body string

	// Duration is how long the request took (in milliseconds)
	Duration int64
}
