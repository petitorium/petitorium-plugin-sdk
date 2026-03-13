package shared

import (
	"context"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/petitorium/petitorium-plugin-sdk/proto"
	"github.com/petitorium/petitorium-plugin-sdk/types"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PETITORIUM_PLUGIN",
	MagicCookieValue: "petitorium",
}

// RPCServer is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	Impl types.Plugin
}

func (s *RPCServer) Name(args interface{}, resp *string) error {
	*resp = s.Impl.Name()
	return nil
}

func (s *RPCServer) Version(args interface{}, resp *string) error {
	*resp = s.Impl.Version()
	return nil
}

func (s *RPCServer) Description(args interface{}, resp *string) error {
	*resp = s.Impl.Description()
	return nil
}

func (s *RPCServer) Hooks(args interface{}, resp *[]types.HookType) error {
	*resp = s.Impl.Hooks()
	return nil
}

type ExecuteHookArgs struct {
	HookType types.HookType
	Context  *types.HookContext
}

type ExecuteHookResponse struct {
	Context *types.HookContext
	Error   string
}

func (s *RPCServer) ExecuteHook(args ExecuteHookArgs, resp *ExecuteHookResponse) error {
	ctx, err := s.Impl.ExecuteHook(args.HookType, args.Context)
	if err != nil {
		resp.Error = err.Error()
	}
	resp.Context = ctx
	return nil
}

// RPCClient is an implementation of Plugin that talks over RPC.
type RPCClient struct {
	client *rpc.Client
}

func (c *RPCClient) Name() string {
	var resp string
	err := c.client.Call("Plugin.Name", new(interface{}), &resp)
	if err != nil {
		return ""
	}
	return resp
}

func (c *RPCClient) Version() string {
	var resp string
	err := c.client.Call("Plugin.Version", new(interface{}), &resp)
	if err != nil {
		return ""
	}
	return resp
}

func (c *RPCClient) Description() string {
	var resp string
	err := c.client.Call("Plugin.Description", new(interface{}), &resp)
	if err != nil {
		return ""
	}
	return resp
}

func (c *RPCClient) Hooks() []types.HookType {
	var resp []types.HookType
	err := c.client.Call("Plugin.Hooks", new(interface{}), &resp)
	if err != nil {
		return nil
	}
	return resp
}

func (c *RPCClient) ExecuteHook(hookType types.HookType, ctx *types.HookContext) (*types.HookContext, error) {
	var resp ExecuteHookResponse
	err := c.client.Call("Plugin.ExecuteHook", ExecuteHookArgs{
		HookType: hookType,
		Context:  ctx,
	}, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return resp.Context, fmt.Errorf("%s", resp.Error)
	}
	return resp.Context, nil
}

// PetitoriumPlugin is the implementation of plugin.Plugin so we can serve/consume this
type PetitoriumPlugin struct {
	Impl types.Plugin
}

func (p *PetitoriumPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (PetitoriumPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

func (p *PetitoriumPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterPluginServer(s, &GRPCServer{
		Impl: p.Impl,
	})
	return nil
}

func (p *PetitoriumPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: proto.NewPluginClient(c),
	}, nil
}
