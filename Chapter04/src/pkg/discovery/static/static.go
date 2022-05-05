package static

import (
	"context"
	"errors"
)

// Registry defines a static service regisry.
type Registry struct {
	serviceAddrs map[string][]string
}

// NewRegistry creates a new static service registry instance.
func NewRegistry(serviceAddrs map[string][]string) *Registry {
	return &Registry{serviceAddrs: serviceAddrs}
}

// Register creates a servie record in the registry.
func (r *Registry) Register(ctx context.Context, id string, url string) error {
	for i, urls := range r.serviceAddrs {
		if i != id {
			continue
		}
		for _, u := range urls {
			if u == url {
				return nil // record already exists.
			}
		}
	}
	return errors.New("adding new records is not supported")
}

// Deregister removes a servie record from the registry.
func (r *Registry) Deregister(ctx context.Context, id string) error {
	return errors.New("removing records is not supported")
}

// ServiceAddresses returns the list of service addresses associated with the given service id.
func (r *Registry) ServiceAddresses(ctx context.Context, id string) ([]string, error) {
	return r.serviceAddrs[id], nil
}
