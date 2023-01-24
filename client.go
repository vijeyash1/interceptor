package interceptor

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// PolicyPayload is the payload that is sent to the server
// in the metadata of the request to check if the user has the required permissions
type CheckPolicy struct {
	Kind string `json:"kind"`
	Scope    string `json:"scope"`
	Role     string `json:"role"`
	Action   string `json:"actions"`
}

// Valid is for checking if all the fields are added to the CheckPolicy
type Valid struct {
	CheckPolicy *CheckPolicy `json:"check_policy"`
}

// AddKind adds the kind to the CheckPolicy
func (p *CheckPolicy) AddRole(role string) *CheckPolicy {
	return &CheckPolicy{
		Role: role,
	}
}

// AddResource adds the resource to the CheckPolicy
func (p *CheckPolicy) AddKind(kind string) *CheckPolicy {
	return &CheckPolicy{
		Kind: kind,
	}
}

// AddScope adds the scope to the CheckPolicy
func (p *CheckPolicy) AddScope(scope string) *CheckPolicy {
	return &CheckPolicy{
		Scope: scope,
	}
}

// AddActions adds the actions to the CheckPolicy
func (p *CheckPolicy) AddActions(action string) *CheckPolicy {
	return &CheckPolicy{
		Action: action,
	}
}

func (p *CheckPolicy) IsValid() bool {
	if p.Kind == "" || p.Scope == "" || p.Role == "" || p.Action == "" {
		return false
	}
	return true
}

func (p *CheckPolicy) GetKind() string {
	return p.Kind
}

func (p *CheckPolicy) GetScope() string {
	return p.Scope
}

func (p *CheckPolicy) GetRole() string {
	return p.Role
}

func (p *CheckPolicy) GetAction() string {
	return p.Action
}

func (p *CheckPolicy) Valid() (*Valid, error) {
	if p.IsValid() {
		return &Valid{
			CheckPolicy: p,
		}, nil
	} else {
		return nil, errors.New("check if all the fields are added to the CheckPolicy")
	}

}

// UnaryClientInterceptor is a client interceptor that appends the scope, resource, role and actions to the metadata
// of the request
// use this function while dialing the server from client
/*
conn, err := grpc.Dial(
  "localhost:5565",
  grpc.WithInsecure(),
  grpc.WithUnaryInterceptor(interceptor.UnaryClientInterceptor()),
)
*/
func (v *Valid) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// append the scope, resource, role and actions to the metadata
		ctx = metadata.AppendToOutgoingContext(ctx, "scope", v.CheckPolicy.Scope)
		ctx = metadata.AppendToOutgoingContext(ctx, "kind", v.CheckPolicy.Kind)
		ctx = metadata.AppendToOutgoingContext(ctx, "role", v.CheckPolicy.Role)
		ctx = metadata.AppendToOutgoingContext(ctx, "action", v.CheckPolicy.Action)
		// Invoke the original method call
		err := invoker(ctx, method, req, reply, cc, opts...)
		log.Printf("client interceptor hit: appending scope,Kind,role,action: '%v %v %v %v ' to metadata", v.CheckPolicy.Scope, v.CheckPolicy.Kind, v.CheckPolicy.Role, v.CheckPolicy.Action)
		return err
	}
}


