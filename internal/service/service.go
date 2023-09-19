package service

import (
	"github.com/google/wire"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewRemoteConfigService, NewEnvManageService, NewRemoteConfigLogService, NewNetworkConfigService, NewAdminAuthService)
