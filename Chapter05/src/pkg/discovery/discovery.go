package discovery

import "context"

// Registry defines a service registry.
type Registry interface {
	Register(ctx context.Context, id string, serviceURL string) error
	Deregister(ctx context.Context, id string) error
	ServiceAddresses(ctx context.Context, id string) ([]string, error)
}
