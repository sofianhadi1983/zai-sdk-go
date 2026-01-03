package zai

import "github.com/z-ai/zai-sdk-go/internal/constants"

// Version returns the current version of the Z.ai Go SDK.
func Version() string {
	return constants.SDKVersion
}
