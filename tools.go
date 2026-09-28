//go:build tools
// +build tools

package main

// Track Wire in go.mod without requiring IDE support for the tool directive.
// Run it with go run github.com/google/wire/cmd/wire ./internal/di.
import _ "github.com/google/wire/cmd/wire"
