package client

import (
	"github.com/google/wire"
	jsoniter "github.com/json-iterator/go"
)

// json serializer
var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

var ProviderSet = wire.NewSet(NewAdminUserClient)
