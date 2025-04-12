package errors

// ErrorTrace represents debug information for an error.
type ErrorTrace struct {
	Message  string      `json:"message"`
	Code     ErrorCode   `json:"code"`
	Location string      `json:"location,omitempty"`
	Params   interface{} `json:"params,omitempty"`
	Cause    *ErrorTrace `json:"cause,omitempty"`
}
