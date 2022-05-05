package consul

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	consul "github.com/hashicorp/consul/api"
)

// Registry defines a Consul-based service regisry.
type Registry struct {
	client *consul.Client
}

// NewRegistry creates a new Consul-based service registry instance.
func NewRegistry(addr string) (*Registry, error) {
	config := consul.DefaultConfig()
	config.Address = addr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{client: client}, nil
}

// Register creates a servie record in the registry.
func (r *Registry) Register(ctx context.Context, id string, serviceURL string) error {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return err
	}
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: u.Hostname(),
		ID:      id,
		Name:    id,
		Port:    port,
	})
}

// Deregister removes a servie record from the registry.
func (r *Registry) Deregister(ctx context.Context, id string) error {
	return r.client.Agent().ServiceDeregister(id)
}

// ServiceAddresses returns the list of service addresses associated with the given service id.
func (r *Registry) ServiceAddresses(ctx context.Context, id string) ([]string, error) {
	addrs, _, err := r.client.Health().Service(id, id, true, nil)
	if err != nil {
		return nil, err
	} else if len(addrs) == 0 {
		return nil, errors.New("not found")
	}
	var res []string
	for _, addr := range addrs {
		res = append(res, addr.Node.Address)
	}
	return res, nil
}
