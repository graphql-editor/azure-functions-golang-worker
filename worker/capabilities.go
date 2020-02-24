package worker

// Capability represents host capabilities
type Capability string

func (c Capability) String() string { return string(c) }

const (
	// RPCHttpTriggerMetadataRemoved capability
	RPCHttpTriggerMetadataRemoved Capability = "RpcHttpTriggerMetadataRemoved"
	// RPCHttpBodyOnly capability
	RPCHttpBodyOnly Capability = "RpcHttpBodyOnly"
	// RawHTTPBodyBytes capability
	RawHTTPBodyBytes Capability = "RawHttpBodyBytes"
)

// Capabilities is a map of capabilites
type Capabilities map[Capability]string

// ToRPC marshals map of capabilities to it's RPC representation
func (c Capabilities) ToRPC() map[string]string {
	m := make(map[string]string, len(c))
	for k, v := range c {
		m[k.String()] = v
	}
	return m
}
