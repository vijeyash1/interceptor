package interceptor

import (
	"context"
	"fmt"

	"github.com/cerbos/cerbos/client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CerbosConfig struct {
	addr string
}

func NewCerbosConfig(addr string) *CerbosConfig {
	return &CerbosConfig{
		addr: addr,
	}
}

func (c *CerbosConfig) client() (client.Client, error) {
	cli, err := client.New(c.addr, client.WithPlaintext())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c *CerbosConfig) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Get the metadata from the incoming context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("couldn't parse incoming context metadata")
		}
		mdscope := md.Get("scope")
		if len(mdscope) == 0 {
			return nil, fmt.Errorf("scope not found in metadata")
		}
		mdkind := md.Get("kind")
		if len(mdkind) == 0 {
			return nil, fmt.Errorf("kind not found in metadata")
		}
		mdrole := md.Get("role")
		if len(mdrole) == 0 {
			return nil, fmt.Errorf("role not found in metadata")
		}
		mdaction := md.Get("action")
		if len(mdaction) == 0 {
			return nil, fmt.Errorf("actions not found in metadata")
		}
		cli, error := c.client()
		if error != nil {
			return nil, error
		}
		allowed, err := cli.IsAllowed(
			context.TODO(),
			client.NewPrincipal("id").WithRoles(mdrole[0]),
			client.NewResource(mdkind[0], "id"),
			mdaction[0],
		)
		if err != nil {
			return nil, err
		}
		h, err := handler(ctx, req)
		if allowed {
			return h, err
		} else {
			return nil, fmt.Errorf("access denied")
		}
	}
}
