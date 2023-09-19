package middleware

import (
	"github.com/google/wire"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// ProviderSet is middleware providers.
var ProviderSet = wire.NewSet(NewValidateWithLogMiddleware, NewIpMiddleware, NewAuthMiddleware)
