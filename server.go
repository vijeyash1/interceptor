package interceptor

import (
	"context"
	"fmt"

	effectv1 "github.com/cerbos/cerbos/api/genpb/cerbos/effect/v1"

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
		mdprinciple := md.Get("principle")
		if len(mdprinciple) == 0 {
			return nil, fmt.Errorf("principle not found in metadata")
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
		principal := client.NewPrincipal(mdprinciple[0], mdrole...)
		resource := client.NewResource(mdkind[0], "intelops")
		resource.WithScope(mdscope[0])
		batch := client.NewResourceBatch()
		batch.Add(resource, mdaction...)
		resp, err := cli.CheckResources(context.Background(), principal, batch)
		if err != nil {
			return nil, err
		}
		
		result := resp.GetResults()
		var allow bool
		for _, r := range mdaction {
			if result[0].Actions[r] == effectv1.Effect_EFFECT_ALLOW {
				allow = true
				
			} else {
				allow = false
			}
		}
		if allow {
			return  handler(ctx, req)
		} else {
			return nil, fmt.Errorf("not allowed")
		}

	}
}
