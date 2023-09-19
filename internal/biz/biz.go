package biz

import (
	"github.com/google/wire"
	jsoniter "github.com/json-iterator/go"
)

// json serializer
var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewRemoteConfigLogic, NewEnvManageLogic, NewRemoteConfigLogLogic, NewNetworkConfigLogic, NewAdminUserLogic, NewAuthLogUsecase)
