// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.0
// - protoc             v4.23.4
// source: api/remoteconfig/remote_config.proto

package remoteconfig

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationRemoteConfigCancelPublishRemoteConfig = "/api.remoteconfig.RemoteConfig/CancelPublishRemoteConfig"
const OperationRemoteConfigCreateRemoteConfig = "/api.remoteconfig.RemoteConfig/CreateRemoteConfig"
const OperationRemoteConfigCreateRemoteConfigV1 = "/api.remoteconfig.RemoteConfig/CreateRemoteConfigV1"
const OperationRemoteConfigDeleteRemoteConfig = "/api.remoteconfig.RemoteConfig/DeleteRemoteConfig"
const OperationRemoteConfigGetRemoteConfig = "/api.remoteconfig.RemoteConfig/GetRemoteConfig"
const OperationRemoteConfigGetRemoteConfigList = "/api.remoteconfig.RemoteConfig/GetRemoteConfigList"
const OperationRemoteConfigGetRemoteConfigV1 = "/api.remoteconfig.RemoteConfig/GetRemoteConfigV1"
const OperationRemoteConfigPublishRemoteConfig = "/api.remoteconfig.RemoteConfig/PublishRemoteConfig"
const OperationRemoteConfigPublishRemoteConfigV1 = "/api.remoteconfig.RemoteConfig/PublishRemoteConfigV1"
const OperationRemoteConfigUpdateRemoteConfig = "/api.remoteconfig.RemoteConfig/UpdateRemoteConfig"
const OperationRemoteConfigUpdateRemoteConfigV1 = "/api.remoteconfig.RemoteConfig/UpdateRemoteConfigV1"

type RemoteConfigHTTPServer interface {
	// CancelPublishRemoteConfig 6. 取消发布远程配置：CancelPublishRemoteConfig
	CancelPublishRemoteConfig(context.Context, *CancelPublishRemoteConfigReq) (*CancelPublishRemoteConfigRsp, error)
	// CreateRemoteConfig 2. 创建远程配置：CreateRemoteConfig
	CreateRemoteConfig(context.Context, *CreateRemoteConfigReq) (*CreateRemoteConfigRsp, error)
	// CreateRemoteConfigV1 旧的创建远程配置：CreateRemoteConfig
	CreateRemoteConfigV1(context.Context, *CreateRemoteConfigV1Req) (*CreateRemoteConfigRsp, error)
	// DeleteRemoteConfig 3. 删除远程配置：DeleteRemoteConfig
	DeleteRemoteConfig(context.Context, *DeleteRemoteConfigReq) (*DeleteRemoteConfigRsp, error)
	// GetRemoteConfig 7. C端公网获取单个远程配置：GetRemoteConfig
	GetRemoteConfig(context.Context, *GetRemoteConfigReq) (*GetRemoteConfigRsp, error)
	// GetRemoteConfigList 1. 获取远程配置列表：GetRemoteConfigList
	GetRemoteConfigList(context.Context, *GetRemoteConfigListReq) (*GetRemoteConfigListRsp, error)
	// GetRemoteConfigV1 旧C端公网获取单个远程配置：GetRemoteConfigV1
	GetRemoteConfigV1(context.Context, *GetRemoteConfigV1Req) (*GetRemoteConfigRsp, error)
	// PublishRemoteConfig 5. 发布远程配置：PublicRemoteConfig
	PublishRemoteConfig(context.Context, *PublishRemoteConfigReq) (*PublishRemoteConfigRsp, error)
	// PublishRemoteConfigV1 旧的发布远程配置：PublicRemoteConfig
	PublishRemoteConfigV1(context.Context, *PublishRemoteConfigV1Req) (*PublishRemoteConfigRsp, error)
	// UpdateRemoteConfig 4. 修改远程配置：UpdateRemoteConfig
	UpdateRemoteConfig(context.Context, *UpdateRemoteConfigReq) (*UpdateRemoteConfigRsp, error)
	// UpdateRemoteConfigV1 旧的修改远程配置：UpdateRemoteConfig
	UpdateRemoteConfigV1(context.Context, *UpdateRemoteConfigV1Req) (*UpdateRemoteConfigRsp, error)
}

func RegisterRemoteConfigHTTPServer(s *http.Server, srv RemoteConfigHTTPServer) {
	r := s.Route("/")
	r.GET("/admin/remote_config/list/{common.appId}/{env}/{channel}/{platform}", _RemoteConfig_GetRemoteConfigList0_HTTP_Handler(srv))
	r.POST("/admin/remote_config/{common.appId}/{env}/{channel}/{platform}/{configName}", _RemoteConfig_CreateRemoteConfig0_HTTP_Handler(srv))
	r.POST("/admin/remote_config/create_remote_config", _RemoteConfig_CreateRemoteConfigV10_HTTP_Handler(srv))
	r.DELETE("/admin/remote_config/{common.appId}/{env}/{channel}/{platform}/{configName}", _RemoteConfig_DeleteRemoteConfig0_HTTP_Handler(srv))
	r.PUT("/admin/remote_config/{common.appId}/{env}/{channel}/{platform}/{configName}", _RemoteConfig_UpdateRemoteConfig0_HTTP_Handler(srv))
	r.POST("/admin/remote_config/update_remote_config", _RemoteConfig_UpdateRemoteConfigV10_HTTP_Handler(srv))
	r.PUT("/admin/remote_config/publish/{common.appId}/{env}/{channel}/{platform}/{configName}", _RemoteConfig_PublishRemoteConfig0_HTTP_Handler(srv))
	r.POST("/admin/remote_config/publish_remote_config", _RemoteConfig_PublishRemoteConfigV10_HTTP_Handler(srv))
	r.PUT("/admin/remote_config/cancel_publish/{common.appId}/{env}/{channel}/{platform}/{configName}", _RemoteConfig_CancelPublishRemoteConfig0_HTTP_Handler(srv))
	r.GET("/api/remote_config/{appId}/{env}/{channel}/{platform}/{configName}", _RemoteConfig_GetRemoteConfig0_HTTP_Handler(srv))
	r.GET("/api/remote_config/get_remote_config/{appid}/{env}/{channel}/{platform}/{configName}", _RemoteConfig_GetRemoteConfigV10_HTTP_Handler(srv))
}

func _RemoteConfig_GetRemoteConfigList0_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRemoteConfigListReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigGetRemoteConfigList)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRemoteConfigList(ctx, req.(*GetRemoteConfigListReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRemoteConfigListRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_CreateRemoteConfig0_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateRemoteConfigReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigCreateRemoteConfig)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateRemoteConfig(ctx, req.(*CreateRemoteConfigReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_CreateRemoteConfigV10_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateRemoteConfigV1Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigCreateRemoteConfigV1)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateRemoteConfigV1(ctx, req.(*CreateRemoteConfigV1Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_DeleteRemoteConfig0_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteRemoteConfigReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigDeleteRemoteConfig)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteRemoteConfig(ctx, req.(*DeleteRemoteConfigReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_UpdateRemoteConfig0_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRemoteConfigReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigUpdateRemoteConfig)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRemoteConfig(ctx, req.(*UpdateRemoteConfigReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_UpdateRemoteConfigV10_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateRemoteConfigV1Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigUpdateRemoteConfigV1)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateRemoteConfigV1(ctx, req.(*UpdateRemoteConfigV1Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_PublishRemoteConfig0_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in PublishRemoteConfigReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigPublishRemoteConfig)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.PublishRemoteConfig(ctx, req.(*PublishRemoteConfigReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*PublishRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_PublishRemoteConfigV10_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in PublishRemoteConfigV1Req
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigPublishRemoteConfigV1)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.PublishRemoteConfigV1(ctx, req.(*PublishRemoteConfigV1Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*PublishRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_CancelPublishRemoteConfig0_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CancelPublishRemoteConfigReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigCancelPublishRemoteConfig)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CancelPublishRemoteConfig(ctx, req.(*CancelPublishRemoteConfigReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CancelPublishRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_GetRemoteConfig0_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRemoteConfigReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigGetRemoteConfig)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRemoteConfig(ctx, req.(*GetRemoteConfigReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

func _RemoteConfig_GetRemoteConfigV10_HTTP_Handler(srv RemoteConfigHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetRemoteConfigV1Req
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRemoteConfigGetRemoteConfigV1)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetRemoteConfigV1(ctx, req.(*GetRemoteConfigV1Req))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetRemoteConfigRsp)
		return ctx.Result(200, reply)
	}
}

type RemoteConfigHTTPClient interface {
	CancelPublishRemoteConfig(ctx context.Context, req *CancelPublishRemoteConfigReq, opts ...http.CallOption) (rsp *CancelPublishRemoteConfigRsp, err error)
	CreateRemoteConfig(ctx context.Context, req *CreateRemoteConfigReq, opts ...http.CallOption) (rsp *CreateRemoteConfigRsp, err error)
	CreateRemoteConfigV1(ctx context.Context, req *CreateRemoteConfigV1Req, opts ...http.CallOption) (rsp *CreateRemoteConfigRsp, err error)
	DeleteRemoteConfig(ctx context.Context, req *DeleteRemoteConfigReq, opts ...http.CallOption) (rsp *DeleteRemoteConfigRsp, err error)
	GetRemoteConfig(ctx context.Context, req *GetRemoteConfigReq, opts ...http.CallOption) (rsp *GetRemoteConfigRsp, err error)
	GetRemoteConfigList(ctx context.Context, req *GetRemoteConfigListReq, opts ...http.CallOption) (rsp *GetRemoteConfigListRsp, err error)
	GetRemoteConfigV1(ctx context.Context, req *GetRemoteConfigV1Req, opts ...http.CallOption) (rsp *GetRemoteConfigRsp, err error)
	PublishRemoteConfig(ctx context.Context, req *PublishRemoteConfigReq, opts ...http.CallOption) (rsp *PublishRemoteConfigRsp, err error)
	PublishRemoteConfigV1(ctx context.Context, req *PublishRemoteConfigV1Req, opts ...http.CallOption) (rsp *PublishRemoteConfigRsp, err error)
	UpdateRemoteConfig(ctx context.Context, req *UpdateRemoteConfigReq, opts ...http.CallOption) (rsp *UpdateRemoteConfigRsp, err error)
	UpdateRemoteConfigV1(ctx context.Context, req *UpdateRemoteConfigV1Req, opts ...http.CallOption) (rsp *UpdateRemoteConfigRsp, err error)
}

type RemoteConfigHTTPClientImpl struct {
	cc *http.Client
}

func NewRemoteConfigHTTPClient(client *http.Client) RemoteConfigHTTPClient {
	return &RemoteConfigHTTPClientImpl{client}
}

func (c *RemoteConfigHTTPClientImpl) CancelPublishRemoteConfig(ctx context.Context, in *CancelPublishRemoteConfigReq, opts ...http.CallOption) (*CancelPublishRemoteConfigRsp, error) {
	var out CancelPublishRemoteConfigRsp
	pattern := "/admin/remote_config/cancel_publish/{common.appId}/{env}/{channel}/{platform}/{configName}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRemoteConfigCancelPublishRemoteConfig))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) CreateRemoteConfig(ctx context.Context, in *CreateRemoteConfigReq, opts ...http.CallOption) (*CreateRemoteConfigRsp, error) {
	var out CreateRemoteConfigRsp
	pattern := "/admin/remote_config/{common.appId}/{env}/{channel}/{platform}/{configName}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRemoteConfigCreateRemoteConfig))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) CreateRemoteConfigV1(ctx context.Context, in *CreateRemoteConfigV1Req, opts ...http.CallOption) (*CreateRemoteConfigRsp, error) {
	var out CreateRemoteConfigRsp
	pattern := "/admin/remote_config/create_remote_config"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRemoteConfigCreateRemoteConfigV1))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) DeleteRemoteConfig(ctx context.Context, in *DeleteRemoteConfigReq, opts ...http.CallOption) (*DeleteRemoteConfigRsp, error) {
	var out DeleteRemoteConfigRsp
	pattern := "/admin/remote_config/{common.appId}/{env}/{channel}/{platform}/{configName}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRemoteConfigDeleteRemoteConfig))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) GetRemoteConfig(ctx context.Context, in *GetRemoteConfigReq, opts ...http.CallOption) (*GetRemoteConfigRsp, error) {
	var out GetRemoteConfigRsp
	pattern := "/api/remote_config/{appId}/{env}/{channel}/{platform}/{configName}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRemoteConfigGetRemoteConfig))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) GetRemoteConfigList(ctx context.Context, in *GetRemoteConfigListReq, opts ...http.CallOption) (*GetRemoteConfigListRsp, error) {
	var out GetRemoteConfigListRsp
	pattern := "/admin/remote_config/list/{common.appId}/{env}/{channel}/{platform}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRemoteConfigGetRemoteConfigList))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) GetRemoteConfigV1(ctx context.Context, in *GetRemoteConfigV1Req, opts ...http.CallOption) (*GetRemoteConfigRsp, error) {
	var out GetRemoteConfigRsp
	pattern := "/api/remote_config/get_remote_config/{appid}/{env}/{channel}/{platform}/{configName}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRemoteConfigGetRemoteConfigV1))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) PublishRemoteConfig(ctx context.Context, in *PublishRemoteConfigReq, opts ...http.CallOption) (*PublishRemoteConfigRsp, error) {
	var out PublishRemoteConfigRsp
	pattern := "/admin/remote_config/publish/{common.appId}/{env}/{channel}/{platform}/{configName}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRemoteConfigPublishRemoteConfig))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) PublishRemoteConfigV1(ctx context.Context, in *PublishRemoteConfigV1Req, opts ...http.CallOption) (*PublishRemoteConfigRsp, error) {
	var out PublishRemoteConfigRsp
	pattern := "/admin/remote_config/publish_remote_config"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRemoteConfigPublishRemoteConfigV1))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) UpdateRemoteConfig(ctx context.Context, in *UpdateRemoteConfigReq, opts ...http.CallOption) (*UpdateRemoteConfigRsp, error) {
	var out UpdateRemoteConfigRsp
	pattern := "/admin/remote_config/{common.appId}/{env}/{channel}/{platform}/{configName}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRemoteConfigUpdateRemoteConfig))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *RemoteConfigHTTPClientImpl) UpdateRemoteConfigV1(ctx context.Context, in *UpdateRemoteConfigV1Req, opts ...http.CallOption) (*UpdateRemoteConfigRsp, error) {
	var out UpdateRemoteConfigRsp
	pattern := "/admin/remote_config/update_remote_config"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationRemoteConfigUpdateRemoteConfigV1))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
