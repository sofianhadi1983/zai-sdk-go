//go:build tools
// +build tools

// Package tools tracks development tool dependencies.
// This file ensures that tool dependencies are included in go.mod
// without being imported in the actual code.
package tools

import (
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/stretchr/testify/assert"
	_ "go.uber.org/mock/mockgen"
)
