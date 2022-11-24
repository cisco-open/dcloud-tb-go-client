package tbclient

// Errors
type vndError struct {
	Logref  string `json:"logref"`
	Message string `json:"message"`
}

// Collections
type collection[T any] struct {
	Embedded T `json:"_embedded"`
}

type embeddedData[R any] interface {
	getData() []R
}

// Topology
type Topology struct {
	Uid         string `json:"uid,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Datacenter  string `json:"datacenter"`
	Notes       string `json:"notes"`
	Status      string `json:"status,omitempty"`
}

type topologyCollection struct {
	Data []Topology `json:"topologies"`
}

func (t topologyCollection) getData() []Topology {
	return t.Data
}

// InventoryNetwork
type InventoryNetwork struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Subnet string `json:"subnet"`
}

type inventoryNetworkCollection struct {
	Data []InventoryNetwork `json:"inventoryNetworks"`
}

func (t inventoryNetworkCollection) getData() []InventoryNetwork {
	return t.Data
}
