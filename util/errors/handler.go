package errors

// ErrorHandler processes errors for debug tracing.
type ErrorHandler struct {
	debugMode bool
	fullStack bool // Toggle full stack trace
}

// NewErrorHandler creates a new ErrorHandler.
func NewErrorHandler(debugMode, fullStack bool) *ErrorHandler {
	return &ErrorHandler{
		debugMode: debugMode,
		fullStack: fullStack,
	}
}

// SetDebugMode sets the debug mode.
func (h *ErrorHandler) SetDebugMode(debug bool) {
	h.debugMode = debug
}

// GetErrorTrace creates an ErrorTrace from an error, using the stored stack trace.
func (h *ErrorHandler) GetErrorTrace(err error) *ErrorTrace {
	if !h.debugMode || err == nil {
		return nil
	}

	var trace *ErrorTrace
	if de, ok := err.(*DomainError); ok {
		location := ""
		if de.Stack != nil {
			location = de.Stack.String(h.fullStack)
		}

		trace = &ErrorTrace{
			Message:  de.Message,
			Code:     de.Code,
			Location: location,
			Params:   de.Params,
		}

		if de.Cause != nil {
			trace.Cause = h.GetErrorTrace(de.Cause)
		}
	} else {
		trace = &ErrorTrace{
			Message: err.Error(),
		}
	}

	return trace
}
