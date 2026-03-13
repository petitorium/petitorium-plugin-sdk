package shared

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/petitorium/petitorium-plugin-sdk/proto"
	"github.com/petitorium/petitorium-plugin-sdk/types"
)

// GRPCServer is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	proto.UnimplementedPluginServer
	Impl types.Plugin
}

func (m *GRPCServer) GetInfo(ctx context.Context, req *proto.Empty) (*proto.PluginInfo, error) {
	hooks := m.Impl.Hooks()
	hookStrings := make([]string, len(hooks))
	for i, h := range hooks {
		hookStrings[i] = string(h)
	}

	return &proto.PluginInfo{
		Name:        m.Impl.Name(),
		Version:     m.Impl.Version(),
		Description: m.Impl.Description(),
		Hooks:       hookStrings,
	}, nil
}

func (m *GRPCServer) ExecuteHook(ctx context.Context, req *proto.HookRequest) (*proto.HookResponse, error) {
	// Convert proto.HookContext to types.HookContext
	tCtx := protoToTypesContext(req.Context)

	resCtx, err := m.Impl.ExecuteHook(types.HookType(req.HookType), tCtx)

	resp := &proto.HookResponse{}
	if err != nil {
		resp.Error = err.Error()
	}
	if resCtx != nil {
		resp.Context = typesToProtoContext(resCtx)
	}

	return resp, nil
}

// GRPCClient is an implementation of types.Plugin that talks over RPC.
type GRPCClient struct {
	client proto.PluginClient

	// Cache these so we don't need to do a network call for static info
	name        string
	version     string
	description string
	hooks       []types.HookType
	infoFetched bool
}

func (m *GRPCClient) fetchInfo() {
	if m.infoFetched {
		return
	}
	resp, err := m.client.GetInfo(context.Background(), &proto.Empty{})
	if err == nil {
		m.name = resp.Name
		m.version = resp.Version
		m.description = resp.Description
		m.hooks = make([]types.HookType, len(resp.Hooks))
		for i, h := range resp.Hooks {
			m.hooks[i] = types.HookType(h)
		}
	}
	m.infoFetched = true
}

func (m *GRPCClient) Name() string {
	m.fetchInfo()
	return m.name
}

func (m *GRPCClient) Version() string {
	m.fetchInfo()
	return m.version
}

func (m *GRPCClient) Description() string {
	m.fetchInfo()
	return m.description
}

func (m *GRPCClient) Hooks() []types.HookType {
	m.fetchInfo()
	return m.hooks
}

func (m *GRPCClient) ExecuteHook(hookType types.HookType, ctx *types.HookContext) (*types.HookContext, error) {
	req := &proto.HookRequest{
		HookType: string(hookType),
		Context:  typesToProtoContext(ctx),
	}

	resp, err := m.client.ExecuteHook(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return protoToTypesContext(resp.Context), errors.New(resp.Error)
	}

	return protoToTypesContext(resp.Context), nil
}

func protoToTypesContext(p *proto.HookContext) *types.HookContext {
	if p == nil {
		return nil
	}

	ctx := &types.HookContext{
		Environment: p.Environment,
	}

	if len(p.ConfigJson) > 0 {
		var cfg map[string]interface{}
		json.Unmarshal(p.ConfigJson, &cfg)
		ctx.Config = cfg
	}

	if p.Request != nil {
		ctx.Request = &types.RequestData{
			Method:      p.Request.Method,
			URL:         p.Request.Url,
			Headers:     p.Request.Headers,
			Body:        p.Request.Body,
			Collection:  p.Request.Collection,
			RequestName: p.Request.RequestName,
		}
	}

	if p.Response != nil {
		headers := make(map[string][]string)
		for k, v := range p.Response.Headers {
			headers[k] = v.Values
		}

		ctx.Response = &types.ResponseData{
			StatusCode: int(p.Response.StatusCode),
			Status:     p.Response.Status,
			Headers:    headers,
			Body:       p.Response.Body,
			Duration:   p.Response.Duration,
		}
	}

	return ctx
}

func typesToProtoContext(t *types.HookContext) *proto.HookContext {
	if t == nil {
		return nil
	}

	ctx := &proto.HookContext{
		Environment: t.Environment,
	}

	if t.Config != nil {
		b, _ := json.Marshal(t.Config)
		ctx.ConfigJson = b
	}

	if t.Request != nil {
		ctx.Request = &proto.RequestData{
			Method:      t.Request.Method,
			Url:         t.Request.URL,
			Headers:     t.Request.Headers,
			Body:        t.Request.Body,
			Collection:  t.Request.Collection,
			RequestName: t.Request.RequestName,
		}
	}

	if t.Response != nil {
		headers := make(map[string]*proto.HeaderList)
		for k, v := range t.Response.Headers {
			headers[k] = &proto.HeaderList{Values: v}
		}

		ctx.Response = &proto.ResponseData{
			StatusCode: int32(t.Response.StatusCode),
			Status:     t.Response.Status,
			Headers:    headers,
			Body:       t.Response.Body,
			Duration:   t.Response.Duration,
		}
	}

	return ctx
}
