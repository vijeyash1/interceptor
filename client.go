package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)
// role action scope resource 
// PolicyPayload is the payload that is sent to the server
// in the metadata of the request to check if the user has the required permissions
type CheckPolicy struct {
	// Kind   string `json:"kind"`
	Scope  string `json:"scope"`
	Role   string `json:"role"`
	Action []string `json:"actions"`
	Principle string `json:"principle"`
}

// Valid is for checking if all the fields are added to the CheckPolicy
type Valid struct {
	CheckPolicy *CheckPolicy `json:"check_policy"`
}

// AddKind adds the kind to the CheckPolicy
func (p *CheckPolicy) AddRole(role string) *CheckPolicy {
	p.Role = role
	return p
}

// AddResource adds the resource to the CheckPolicy
// func (p *CheckPolicy) AddKind(kind string) *CheckPolicy {
// 	p.Kind = kind
// 	return p
// }

// AddScope adds the scope to the CheckPolicy
func (p *CheckPolicy) AddScope(scope string) *CheckPolicy {
	p.Scope = scope
	return p
}

// AddActions adds the actions to the CheckPolicy
func (p *CheckPolicy) AddActions(action []string) *CheckPolicy {
	for _, a := range action {
		p.Action = append(p.Action, a)
	}
	return p
}

// AddPrinciple adds the principle to the CheckPolicy
func (p *CheckPolicy) AddPrinciple(principle string) *CheckPolicy {
	p.Principle = principle
	return p
}

func (p *CheckPolicy) IsValid() bool {
	if p.Scope == "" || p.Role == "" || len(p.Action) == 0 || p.Principle == ""{
		return false
	}
	return true
}

// func (p *CheckPolicy) GetKind() string {
// 	return p.Kind
// }

func (p *CheckPolicy) GetScope() string {
	return p.Scope
}

func (p *CheckPolicy) GetRole() string {
	return p.Role
}

func (p *CheckPolicy) GetAction() []string {
	return p.Action
}

func (p *CheckPolicy) GetPrinciple() string {
	return p.Principle
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
		// ctx = metadata.AppendToOutgoingContext(ctx, "kind", v.CheckPolicy.Kind)
		ctx = metadata.AppendToOutgoingContext(ctx, "role", v.CheckPolicy.Role)
		for _, a := range v.CheckPolicy.Action {
			ctx = metadata.AppendToOutgoingContext(ctx, "action", a)
		}
		ctx = metadata.AppendToOutgoingContext(ctx, "principle", v.CheckPolicy.Principle)
		// Invoke the original method call
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
