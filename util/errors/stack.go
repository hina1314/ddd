package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// StackTrace represents a filtered stack trace.
type StackTrace struct {
	Frames []Frame
}

// Frame represents a single stack frame.
type Frame struct {
	File     string
	Line     int
	Function string
}

// CaptureStack captures a stack trace, skipping the specified number of frames.
func CaptureStack(skip int) *StackTrace {
	const maxDepth = 32
	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(skip+1, pcs) // Skip captureStack itself
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(pcs[:n])
	var stack []Frame
	for {
		frame, more := frames.Next()
		if shouldSkipFrame(frame) {
			continue
		}

		file := filepath.Base(frame.File) // Use base name for brevity
		function := simplifyFunctionName(frame.Function)
		stack = append(stack, Frame{
			File:     file,
			Line:     frame.Line,
			Function: function,
		})

		if !more {
			break
		}
	}

	if len(stack) == 0 {
		return nil
	}
	return &StackTrace{Frames: stack}
}

// shouldSkipFrame determines if a frame should be excluded from the stack trace.
func shouldSkipFrame(frame runtime.Frame) bool {
	file := frame.File

	// 跳过标准库
	if strings.HasPrefix(file, runtime.GOROOT()) {
		return true
	}

	// 跳过 Go modules 缓存（一般是第三方库）
	if strings.Contains(file, "pkg/mod") {
		return true
	}

	return false
}

// simplifyFunctionName extracts a concise function name.
func simplifyFunctionName(name string) string {
	parts := strings.Split(name, "/")
	if len(parts) > 0 {
		name = parts[len(parts)-1]
	}
	if idx := strings.LastIndex(name, "."); idx != -1 {
		return name[idx+1:]
	}
	return name
}

// String formats the stack trace as a single-line location or multi-line stack.
func (s *StackTrace) String(full bool) string {
	if s == nil || len(s.Frames) == 0 {
		return ""
	}

	if !full {
		f := s.Frames[0]
		return fmt.Sprintf("%s:%d (%s)", f.File, f.Line, f.Function)
	}

	var builder strings.Builder
	for i, f := range s.Frames {
		if i > 0 {
			builder.WriteString(" -> ")
		}
		fmt.Fprintf(&builder, "%s:%d (%s)", f.File, f.Line, f.Function)
	}
	return builder.String()
}
